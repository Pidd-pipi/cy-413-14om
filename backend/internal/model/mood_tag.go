package model

import "time"

type MoodTag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex:idx_mood_tags_user_name;not null" json:"user_id"`
	Name      string    `gorm:"type:text;uniqueIndex:idx_mood_tags_user_name;not null" json:"name"`
	IsActive  bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}
