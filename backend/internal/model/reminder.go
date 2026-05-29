package model

import "time"

type Reminder struct {
	Model
	UserID     uint       `gorm:"index;not null" json:"user_id"`
	TargetType string     `gorm:"size:30;not null" json:"target_type"`
	TargetID   uint       `gorm:"not null" json:"target_id"`
	RemindType string     `gorm:"size:30;not null" json:"remind_type"`
	RemindDate *time.Time `json:"remind_date"`
	Title      string     `gorm:"size:200;not null" json:"title"`
	Content    string     `gorm:"size:500" json:"content"`
	IsRead     bool       `gorm:"default:false" json:"is_read"`
}
