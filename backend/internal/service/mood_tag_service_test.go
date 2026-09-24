package service

import (
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/blueship581/mindgarden/backend/internal/constants"
	"github.com/blueship581/mindgarden/backend/internal/dto"
	"github.com/blueship581/mindgarden/backend/internal/model"
	"github.com/blueship581/mindgarden/backend/internal/repository"
	"github.com/blueship581/mindgarden/backend/internal/util"
)

type fakeMoodTagRepo struct {
	tags map[uint]*model.MoodTagDef
	next uint
}

func newFakeMoodTagRepo() *fakeMoodTagRepo {
	return &fakeMoodTagRepo{tags: map[uint]*model.MoodTagDef{}}
}
func (f *fakeMoodTagRepo) Create(v *model.MoodTagDef) error {
	f.next++
	v.ID = f.next
	cp := *v
	f.tags[v.ID] = &cp
	return nil
}
func (f *fakeMoodTagRepo) ListActive(uid uint) (out []model.MoodTagDef, e error) {
	for _, t := range f.tags {
		if t.UserID == uid && t.IsActive {
			out = append(out, *t)
		}
	}
	return
}
func (f *fakeMoodTagRepo) CountActive(uid uint) (n int64, e error) {
	for _, t := range f.tags {
		if t.UserID == uid && t.IsActive {
			n++
		}
	}
	return
}
func (f *fakeMoodTagRepo) ByName(uid uint, name string) (*model.MoodTagDef, error) {
	for _, t := range f.tags {
		if t.UserID == uid && t.Name == name {
			cp := *t
			return &cp, nil
		}
	}
	return nil, repository.ErrNotFound
}
func (f *fakeMoodTagRepo) ByID(id, uid uint) (*model.MoodTagDef, error) {
	t, ok := f.tags[id]
	if !ok || t.UserID != uid {
		return nil, repository.ErrNotFound
	}
	cp := *t
	return &cp, nil
}
func (f *fakeMoodTagRepo) Save(v *model.MoodTagDef) error {
	cp := *v
	f.tags[v.ID] = &cp
	return nil
}

type fakeMoodRepoForTags struct{ tagRepo *fakeMoodTagRepo }

func (f *fakeMoodRepoForTags) Create(*model.Mood) error                    { return nil }
func (f *fakeMoodRepoForTags) List(uint, *time.Time) ([]model.Mood, error) { return nil, nil }
func (f *fakeMoodRepoForTags) ByID(uint, uint) (*model.Mood, error)        { return nil, repository.ErrNotFound }
func (f *fakeMoodRepoForTags) Update(*model.Mood) error                    { return nil }
func (f *fakeMoodRepoForTags) Delete(*model.Mood) error                    { return nil }

func TestMoodTagServiceCreateValidation(t *testing.T) {
	s := NewMoodTagService(newFakeMoodTagRepo(), slog.New(slog.NewTextHandler(io.Discard, nil)))
	cases := []struct {
		name    string
		input   string
		wantErr bool
		code    int
	}{
		{"blank", "   ", true, constants.CodeValidation},
		{"too long", strings.Repeat("兴", 9), true, constants.CodeValidation},
		{"built in name", "happy", true, constants.CodeValidation},
		{"ok", "兴奋", false, constants.CodeOK},
		{"trimmed ok", "  期待  ", false, constants.CodeOK},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			v, e := s.Create(1, tt.input)
			if (e != nil) != tt.wantErr {
				t.Fatalf("err=%v", e)
			}
			if tt.wantErr {
				var app *util.AppError
				if !errors.As(e, &app) || app.Code != tt.code {
					t.Fatalf("want code %d, got %v", tt.code, e)
				}
				return
			}
			if v.Name != strings.TrimSpace(tt.input) || !v.IsActive {
				t.Fatalf("unexpected tag %+v", v)
			}
		})
	}
}

func TestMoodTagServiceDuplicateSameAccount(t *testing.T) {
	s := NewMoodTagService(newFakeMoodTagRepo(), slog.New(slog.NewTextHandler(io.Discard, nil)))
	if _, e := s.Create(1, "兴奋"); e != nil {
		t.Fatal(e)
	}
	_, e := s.Create(1, "兴奋")
	var app *util.AppError
	if !errors.As(e, &app) || app.Code != constants.CodeConflict {
		t.Fatalf("same account same active name must conflict, got %v", e)
	}
	// 另一账号同名允许存在。
	if _, e := s.Create(2, "兴奋"); e != nil {
		t.Fatalf("other account should own the same name: %v", e)
	}
}

func TestMoodTagServiceLimitTwelve(t *testing.T) {
	s := NewMoodTagService(newFakeMoodTagRepo(), slog.New(slog.NewTextHandler(io.Discard, nil)))
	for i := 0; i < 12; i++ {
		if _, e := s.Create(1, "标签"+string(rune('A'+i))); e != nil {
			t.Fatalf("create %d failed: %v", i, e)
		}
	}
	_, e := s.Create(1, "超出")
	var app *util.AppError
	if !errors.As(e, &app) || app.Code != constants.CodeValidation {
		t.Fatalf("13th tag must be rejected, got %v", e)
	}
}

func TestMoodTagServiceDeactivateAndReactivate(t *testing.T) {
	s := NewMoodTagService(newFakeMoodTagRepo(), slog.New(slog.NewTextHandler(io.Discard, nil)))
	v, e := s.Create(1, "兴奋")
	if e != nil {
		t.Fatal(e)
	}
	if e = s.Deactivate(1, v.ID); e != nil {
		t.Fatal(e)
	}
	if active, _ := s.List(1); len(active) != 0 {
		t.Fatalf("deactivated tag must leave selection list")
	}
	// 停用后再用同名新建：恢复原标签，不算新增名额占用。
	v2, e := s.Create(1, "兴奋")
	if e != nil {
		t.Fatalf("recreate after deactivate should reactivate: %v", e)
	}
	if v2.ID != v.ID || !v2.IsActive {
		t.Fatalf("expected same tag reactivated, got %+v", v2)
	}
	// 别人账号的标签不能被停用。
	if e = s.Deactivate(2, v.ID); e == nil {
		t.Fatal("deactivating another account tag must fail")
	}
}

func TestMoodServiceRejectsTagsNotSelectable(t *testing.T) {
	tagRepo := newFakeMoodTagRepo()
	tagSvc := NewMoodTagService(tagRepo, slog.New(slog.NewTextHandler(io.Discard, nil)))
	moodSvc := NewMoodService(&fakeMoodRepoForTags{tagRepo}, tagRepo, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if _, e := tagSvc.Create(1, "兴奋"); e != nil {
		t.Fatal(e)
	}
	req := dto.MoodRequest{MoodLevel: 7, MoodTags: []string{"兴奋"}, RecordDate: "2026-09-24"}
	if _, e := moodSvc.Create(1, req); e != nil {
		t.Fatalf("own active custom tag should be accepted: %v", e)
	}
	// 停用后不再可用于保存。
	if _, e := tagSvc.Create(1, "待停"); e != nil {
		t.Fatal(e)
	}
	tags, _ := tagSvc.List(1)
	var id uint
	for _, x := range tags {
		if x.Name == "待停" {
			id = x.ID
		}
	}
	if e := tagSvc.Deactivate(1, id); e != nil {
		t.Fatal(e)
	}
	if _, e := moodSvc.Create(1, dto.MoodRequest{MoodLevel: 7, MoodTags: []string{"待停"}, RecordDate: "2026-09-24"}); e == nil {
		t.Fatal("deactivated tag must be rejected when saving mood")
	}
	// 别人账号的标签不能被拿来使用。
	if _, e := tagSvc.Create(2, "期待"); e != nil {
		t.Fatal(e)
	}
	if _, e := moodSvc.Create(1, dto.MoodRequest{MoodLevel: 7, MoodTags: []string{"期待"}, RecordDate: "2026-09-24"}); e == nil {
		t.Fatal("other account tag must not be usable")
	}
	// 未知标签同样拒绝。
	_, e := moodSvc.Create(1, dto.MoodRequest{MoodLevel: 7, MoodTags: []string{"瞎编的"}, RecordDate: "2026-09-24"})
	var app *util.AppError
	if !errors.As(e, &app) || app.Code != constants.CodeValidation {
		t.Fatalf("unknown tag must fail validation, got %v", e)
	}
}
