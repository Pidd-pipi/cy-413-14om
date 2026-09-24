package service

import (
	"errors"
	"log/slog"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/blueship581/mindgarden/backend/internal/constants"
	"github.com/blueship581/mindgarden/backend/internal/model"
	"github.com/blueship581/mindgarden/backend/internal/repository"
	"github.com/blueship581/mindgarden/backend/internal/util"
)

type MoodTagService struct {
	repo   repository.MoodTagRepository
	logger *slog.Logger
}

func NewMoodTagService(r repository.MoodTagRepository, l *slog.Logger) *MoodTagService {
	return &MoodTagService{r, l}
}

// List 返回当前账号当前可选择的自定义标签（不含已停用）。
func (s *MoodTagService) List(uid uint) ([]model.MoodTagDef, error) {
	vs, e := s.repo.ListActive(uid)
	if e != nil {
		return nil, errors.New("MoodTag[user_id] list failed: " + e.Error())
	}
	return vs, nil
}

// Create 校验空白、八字上限、重名与十二个上限；重名若为已停用标签则原地恢复。
func (s *MoodTagService) Create(uid uint, name string) (*model.MoodTagDef, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, util.NewAppError(constants.CodeValidation, constants.MessageMoodTagNameBlank, nil)
	}
	if utf8.RuneCountInString(name) > constants.CustomMoodTagMaxRunes {
		return nil, util.NewAppError(constants.CodeValidation, constants.MessageMoodTagNameTooLong, nil)
	}
	if isBuiltInMoodTag(name) {
		return nil, util.NewAppError(constants.CodeValidation, constants.MessageMoodTagNameBuiltIn, nil)
	}
	existing, e := s.repo.ByName(uid, name)
	if e != nil && !errors.Is(e, repository.ErrNotFound) {
		return nil, errors.New("MoodTag[name] create failed: " + e.Error())
	}
	if existing != nil {
		if existing.IsActive {
			return nil, util.NewAppError(constants.CodeConflict, constants.MessageMoodTagNameExists, nil)
		}
		existing.IsActive = true
		existing.UpdatedAt = time.Now()
		if e = s.repo.Save(existing); e != nil {
			return nil, errors.New("MoodTag[name] reactivate failed: " + e.Error())
		}
		s.logger.Info(constants.LogMoodTagReactivated, "user_id", uid, "mood_tag_id", existing.ID)
		return existing, nil
	}
	count, e := s.repo.CountActive(uid)
	if e != nil {
		return nil, errors.New("MoodTag[user_id] count failed: " + e.Error())
	}
	if count >= constants.CustomMoodTagMaxCount {
		return nil, util.NewAppError(constants.CodeValidation, constants.MessageMoodTagLimit, nil)
	}
	v := &model.MoodTagDef{UserID: uid, Name: name, IsActive: true}
	if e = s.repo.Create(v); e != nil {
		if errors.Is(e, repository.ErrConflict) {
			return nil, util.NewAppError(constants.CodeConflict, constants.MessageMoodTagNameExists, nil)
		}
		return nil, errors.New("MoodTag[name] create failed: " + e.Error())
	}
	s.logger.Info(constants.LogMoodTagCreated, "user_id", uid, "mood_tag_id", v.ID)
	return v, nil
}

// Deactivate 停用标签：停用后不再出现在选择区，历史记录文本不受影响。
func (s *MoodTagService) Deactivate(uid, id uint) error {
	v, e := s.repo.ByID(id, uid)
	if e != nil {
		if errors.Is(e, repository.ErrNotFound) {
			return util.NewAppError(constants.CodeNotFound, constants.MessageMoodTagNotOwner, nil)
		}
		return errors.New("MoodTag[id] deactivate failed: " + e.Error())
	}
	if v.IsActive {
		v.IsActive = false
		v.UpdatedAt = time.Now()
		if e = s.repo.Save(v); e != nil {
			return errors.New("MoodTag[id] deactivate failed: " + e.Error())
		}
	}
	s.logger.Info(constants.LogMoodTagDeactivated, "user_id", uid, "mood_tag_id", id)
	return nil
}
