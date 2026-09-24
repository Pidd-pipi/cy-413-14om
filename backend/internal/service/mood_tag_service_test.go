package service

import (
	"errors"
	"github.com/blueship581/mindgarden/backend/internal/constants"
	"github.com/blueship581/mindgarden/backend/internal/dto"
	"github.com/blueship581/mindgarden/backend/internal/model"
	"github.com/blueship581/mindgarden/backend/internal/repository"
	"github.com/blueship581/mindgarden/backend/internal/util"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"
)

type fakeMoodTagRepo struct {
	tags map[uint]*model.MoodTag
	next uint
}

func newFakeMoodTagRepo() *fakeMoodTagRepo { return &fakeMoodTagRepo{tags: map[uint]*model.MoodTag{}} }
func (f *fakeMoodTagRepo) Create(v *model.MoodTag) error {
	f.next++
	v.ID = f.next
	f.tags[v.ID] = v
	return nil
}
func (f *fakeMoodTagRepo) Update(v *model.MoodTag) error { f.tags[v.ID] = v; return nil }
func (f *fakeMoodTagRepo) ByID(id, uid uint) (*model.MoodTag, error) {
	if v, ok := f.tags[id]; ok && v.UserID == uid {
		return v, nil
	}
	return nil, repository.ErrNotFound
}
func (f *fakeMoodTagRepo) ByName(uid uint, name string) (*model.MoodTag, error) {
	for _, v := range f.tags {
		if v.UserID == uid && v.Name == name {
			return v, nil
		}
	}
	return nil, repository.ErrNotFound
}
func (f *fakeMoodTagRepo) ListActive(uid uint) (out []model.MoodTag, e error) {
	for _, v := range f.tags {
		if v.UserID == uid && v.IsActive {
			out = append(out, *v)
		}
	}
	return
}
func (f *fakeMoodTagRepo) CountActive(uid uint) (n int64, e error) {
	for _, v := range f.tags {
		if v.UserID == uid && v.IsActive {
			n++
		}
	}
	return
}

type fakeMoodRepo struct{ moods []*model.Mood }

func (f *fakeMoodRepo) Create(v *model.Mood) error {
	v.ID = uint(len(f.moods) + 1)
	f.moods = append(f.moods, v)
	return nil
}
func (f *fakeMoodRepo) List(uint, *time.Time) ([]model.Mood, error) { return nil, nil }
func (f *fakeMoodRepo) ByID(id, uid uint) (*model.Mood, error) {
	for _, v := range f.moods {
		if v.ID == id && v.UserID == uid {
			return v, nil
		}
	}
	return nil, repository.ErrNotFound
}
func (f *fakeMoodRepo) Update(v *model.Mood) error { return nil }
func (f *fakeMoodRepo) Delete(v *model.Mood) error { return nil }

func testLogger() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }

func TestMoodTagServiceCreateValidation(t *testing.T) {
	cases := []struct {
		name    string
		tagName string
		wantErr bool
	}{
		{"blank", "", true},
		{"whitespace only", "   ", true},
		{"nine chars", strings.Repeat("好", 9), true},
		{"eight chars", strings.Repeat("好", 8), false},
		{"normal", "兴奋", false},
		{"trimmed", "  期待  ", false},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMoodTagService(newFakeMoodTagRepo(), testLogger())
			v, e := s.Create(1, tt.tagName)
			if (e != nil) != tt.wantErr {
				t.Fatalf("err=%v", e)
			}
			if tt.wantErr {
				var app *util.AppError
				if !errors.As(e, &app) || app.Code != constants.CodeValidation {
					t.Fatalf("expected validation app error, got %v", e)
				}
			} else if v == nil || v.Name != strings.TrimSpace(tt.tagName) || !v.IsActive {
				t.Fatalf("unexpected tag %+v", v)
			}
		})
	}
}

func TestMoodTagServiceSameNameKeepsOne(t *testing.T) {
	repo := newFakeMoodTagRepo()
	s := NewMoodTagService(repo, testLogger())
	a, e := s.Create(1, "兴奋")
	if e != nil {
		t.Fatalf("create: %v", e)
	}
	b, e := s.Create(1, "兴奋")
	if e != nil {
		t.Fatalf("duplicate create must not fail: %v", e)
	}
	if a.ID != b.ID {
		t.Fatalf("same name must return the same tag, got %d and %d", a.ID, b.ID)
	}
	if n, _ := repo.CountActive(1); n != 1 {
		t.Fatalf("same name must keep only one tag, got %d", n)
	}
	if e = s.Deactivate(1, a.ID); e != nil {
		t.Fatalf("deactivate: %v", e)
	}
	if vs, _ := s.ListActive(1); len(vs) != 0 {
		t.Fatalf("deactivated tag must leave the selectable list, got %+v", vs)
	}
	c, e := s.Create(1, "兴奋")
	if e != nil || c.ID != a.ID || !c.IsActive {
		t.Fatalf("recreating a deactivated name must reactivate the same tag: %+v %v", c, e)
	}
	if n, _ := repo.CountActive(1); n != 1 {
		t.Fatalf("reactivation must keep only one tag, got %d", n)
	}
}

func TestMoodTagServiceLimit(t *testing.T) {
	repo := newFakeMoodTagRepo()
	s := NewMoodTagService(repo, testLogger())
	for i := 0; i < constants.MoodTagMaxUserCount; i++ {
		if _, e := s.Create(1, strings.Repeat("签", 1)+strings.Repeat("好", i%7)+string(rune('a'+i))); e != nil {
			t.Fatalf("create %d: %v", i, e)
		}
	}
	if _, e := s.Create(1, "超出的标签"); e == nil {
		t.Fatal("13th active tag must be rejected")
	} else if !strings.Contains(e.Error(), constants.MessageMoodTagLimit) {
		t.Fatalf("limit error must explain itself, got %v", e)
	}
	vs, _ := s.ListActive(1)
	if e := s.Deactivate(1, vs[0].ID); e != nil {
		t.Fatalf("deactivate: %v", e)
	}
	if _, e := s.Create(1, "补位的标签"); e != nil {
		t.Fatalf("a freed slot must allow a new tag: %v", e)
	}
	other := NewMoodTagService(repo, testLogger())
	if _, e := other.Create(2, "自己的标签"); e != nil {
		t.Fatalf("limit is per account: %v", e)
	}
}

func TestMoodServiceOnlySelectableTags(t *testing.T) {
	tagRepo := newFakeMoodTagRepo()
	tags := NewMoodTagService(tagRepo, testLogger())
	own, e := tags.Create(1, "兴奋")
	if e != nil {
		t.Fatalf("create tag: %v", e)
	}
	if _, e = tags.Create(2, "期待"); e != nil {
		t.Fatalf("create other account tag: %v", e)
	}
	moods := NewMoodService(&fakeMoodRepo{}, tagRepo, testLogger())
	req := func(names ...string) dto.MoodRequest {
		return dto.MoodRequest{MoodLevel: 6, MoodTags: names, RecordDate: "2026-09-24"}
	}
	cases := []struct {
		name    string
		uid     uint
		tags    []string
		wantErr bool
	}{
		{"builtin only", 1, []string{"happy", "calm"}, false},
		{"own custom tag", 1, []string{"happy", "兴奋"}, false},
		{"other account tag", 1, []string{"期待"}, true},
		{"unknown tag", 1, []string{"不存在"}, true},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			_, e := moods.Create(tt.uid, req(tt.tags...))
			if (e != nil) != tt.wantErr {
				t.Fatalf("err=%v", e)
			}
		})
	}
	if e = tags.Deactivate(1, own.ID); e != nil {
		t.Fatalf("deactivate: %v", e)
	}
	if _, e = moods.Create(1, req("兴奋")); e == nil {
		t.Fatal("deactivated tag is no longer selectable and must be rejected")
	}
	if _, e = moods.Create(1, req("happy")); e != nil {
		t.Fatalf("builtin tags stay selectable: %v", e)
	}
}
