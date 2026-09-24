package repository

import (
	"errors"
	"github.com/blueship581/mindgarden/backend/internal/model"
	"gorm.io/gorm"
)

type MoodTagRepository interface {
	Create(*model.MoodTag) error
	Update(*model.MoodTag) error
	ByID(uint, uint) (*model.MoodTag, error)
	ByName(uint, string) (*model.MoodTag, error)
	ListActive(uint) ([]model.MoodTag, error)
	CountActive(uint) (int64, error)
}
type moodTagRepository struct{ db *gorm.DB }

func NewMoodTagRepository(db *gorm.DB) MoodTagRepository { return &moodTagRepository{db} }
func (r *moodTagRepository) Create(v *model.MoodTag) error {
	return r.db.Create(v).Error
}
func (r *moodTagRepository) Update(v *model.MoodTag) error { return r.db.Save(v).Error }
func (r *moodTagRepository) ByID(id, uid uint) (*model.MoodTag, error) {
	var v model.MoodTag
	e := r.db.Where("id = ? AND user_id = ?", id, uid).First(&v).Error
	if errors.Is(e, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &v, e
}
func (r *moodTagRepository) ByName(uid uint, name string) (*model.MoodTag, error) {
	var v model.MoodTag
	e := r.db.Where("user_id = ? AND name = ?", uid, name).First(&v).Error
	if errors.Is(e, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &v, e
}
func (r *moodTagRepository) ListActive(uid uint) (out []model.MoodTag, e error) {
	e = r.db.Where("user_id = ? AND is_active = ?", uid, true).Order("id asc").Find(&out).Error
	return
}
func (r *moodTagRepository) CountActive(uid uint) (n int64, e error) {
	e = r.db.Model(&model.MoodTag{}).Where("user_id = ? AND is_active = ?", uid, true).Count(&n).Error
	return
}
