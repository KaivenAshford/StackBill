package model

import "time"

type NotificationSetting struct {
	ID              uint      `gorm:"primarykey" json:"id"`
	UserID          uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	EmailEnabled    bool      `gorm:"default:false" json:"email_enabled"`
	RemindDaysBefore int      `gorm:"default:3" json:"remind_days_before"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
