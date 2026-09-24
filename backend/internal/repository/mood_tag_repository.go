package repository

import (
	"errors"
	"github.com/blueship581/mindgarden/backend/internal/model"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type MoodTagRepository interface {
	Create(*model.MoodTagDef) error
	ListActive(uint) ([]model.MoodTagDef, error)
	CountActive(uint) (int64, error)
	ByName(uid uint, name string) (*model.MoodTagDef, error)
	ByID(id, uid uint) (*model.MoodTagDef, error)
	Save(*model.MoodTagDef) error
}

type moodTagRepository struct{ db *gorm.DB }

func NewMoodTagRepository(db *gorm.DB) MoodTagRepository { return &moodTagRepository{db} }

func (r *moodTagRepository) Create(v *model.MoodTagDef) error {
	err := r.db.Create(v).Error
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrConflict
	}
	return err
}
func (r *moodTagRepository) ListActive(uid uint) (out []model.MoodTagDef, e error) {
	e = r.db.Where("user_id = ? AND is_active = ?", uid, true).Order("id asc").Find(&out).Error
	return
}
func (r *moodTagRepository) CountActive(uid uint) (n int64, e error) {
	e = r.db.Model(&model.MoodTagDef{}).Where("user_id = ? AND is_active = ?", uid, true).Count(&n).Error
	return
}
func (r *moodTagRepository) ByName(uid uint, name string) (*model.MoodTagDef, error) {
	var v model.MoodTagDef
	e := r.db.Where("user_id = ? AND name = ?", uid, name).First(&v).Error
	if errors.Is(e, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &v, e
}
func (r *moodTagRepository) ByID(id, uid uint) (*model.MoodTagDef, error) {
	var v model.MoodTagDef
	e := r.db.Where("id = ? AND user_id = ?", id, uid).First(&v).Error
	if errors.Is(e, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &v, e
}
func (r *moodTagRepository) Save(v *model.MoodTagDef) error { return r.db.Save(v).Error }
