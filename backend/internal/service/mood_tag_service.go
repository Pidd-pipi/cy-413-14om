package service

import (
	"errors"
	"fmt"
	"github.com/blueship581/mindgarden/backend/internal/constants"
	"github.com/blueship581/mindgarden/backend/internal/model"
	"github.com/blueship581/mindgarden/backend/internal/repository"
	"github.com/blueship581/mindgarden/backend/internal/util"
	"log/slog"
	"strings"
	"unicode/utf8"
)

type MoodTagService struct {
	repo   repository.MoodTagRepository
	logger *slog.Logger
}

func NewMoodTagService(r repository.MoodTagRepository, l *slog.Logger) *MoodTagService {
	return &MoodTagService{r, l}
}
func (s *MoodTagService) Create(uid uint, name string) (*model.MoodTag, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, util.NewAppError(constants.CodeValidation, "MoodTag[name] create failed: "+constants.MessageMoodTagBlank, nil)
	}
	if utf8.RuneCountInString(name) > constants.MoodTagMaxNameLength {
		return nil, util.NewAppError(constants.CodeValidation, "MoodTag[name] create failed: "+constants.MessageMoodTagTooLong, nil)
	}
	existing, e := s.repo.ByName(uid, name)
	if e != nil && !errors.Is(e, repository.ErrNotFound) {
		return nil, fmt.Errorf("MoodTag[user_id] fetch failed: %w", e)
	}
	if existing != nil {
		if !existing.IsActive {
			existing.IsActive = true
			if e = s.repo.Update(existing); e != nil {
				return nil, fmt.Errorf("MoodTag[id=%d] reactivate failed: %w", existing.ID, e)
			}
			s.logger.Info(constants.LogMoodTagCreated, "user_id", uid, "name", name, "reactivated", true)
		}
		return existing, nil
	}
	n, e := s.repo.CountActive(uid)
	if e != nil {
		return nil, fmt.Errorf("MoodTag[user_id] count failed: %w", e)
	}
	if n >= constants.MoodTagMaxUserCount {
		return nil, util.NewAppError(constants.CodeValidation, "MoodTag[user_id] create failed: "+constants.MessageMoodTagLimit, nil)
	}
	v := &model.MoodTag{UserID: uid, Name: name, IsActive: true}
	if e = s.repo.Create(v); e != nil {
		return nil, fmt.Errorf("MoodTag[user_id] create failed: %w", e)
	}
	s.logger.Info(constants.LogMoodTagCreated, "user_id", uid, "name", name)
	return v, nil
}
func (s *MoodTagService) ListActive(uid uint) ([]model.MoodTag, error) {
	vs, e := s.repo.ListActive(uid)
	if e != nil {
		return nil, fmt.Errorf("MoodTag[user_id] list failed: %w", e)
	}
	s.logger.Info(constants.LogMoodTagListed, "user_id", uid)
	return vs, nil
}
func (s *MoodTagService) Deactivate(uid, id uint) error {
	v, e := s.repo.ByID(id, uid)
	if e != nil {
		return fmt.Errorf("MoodTag[id=%d] fetch failed: %w", id, e)
	}
	v.IsActive = false
	if e = s.repo.Update(v); e != nil {
		return util.WrapEntity("MoodTag", "is_active", id, constants.CodeInternal, e)
	}
	s.logger.Info(constants.LogMoodTagDeactivated, "mood_tag_id", id)
	return nil
}
