package model

import "time"

type MoodTagDef struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null;uniqueIndex:idx_mood_tag_owner_name,priority:2" json:"user_id"`
	Name      string    `gorm:"type:varchar(32);not null;uniqueIndex:idx_mood_tag_owner_name,priority:3" json:"name"`
	IsActive  bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
