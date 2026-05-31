package model

import "time"

type ReminderDismissed struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	UserID     uint      `gorm:"uniqueIndex:idx_reminder_dismissed;not null" json:"user_id"`
	TargetType string    `gorm:"uniqueIndex:idx_reminder_dismissed;size:30;not null" json:"target_type"`
	TargetID   uint      `gorm:"uniqueIndex:idx_reminder_dismissed;not null" json:"target_id"`
	CreatedAt  time.Time `json:"created_at"`
}
