package repository

import (
	"fmt"
	"time"

	"github.com/kingqaquuu/stackbill/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReminderRepository struct {
	db *gorm.DB
}

func NewReminderRepository(db *gorm.DB) *ReminderRepository {
	return &ReminderRepository{db: db}
}

func (r *ReminderRepository) GetReadKeys(userID uint) (map[string]bool, error) {
	var reads []model.ReminderRead
	if err := r.db.Where("user_id = ?", userID).Find(&reads).Error; err != nil {
		return nil, err
	}
	keys := make(map[string]bool, len(reads))
	for _, rd := range reads {
		key := fmt.Sprintf("%s-%d", rd.TargetType, rd.TargetID)
		keys[key] = true
	}
	return keys, nil
}

func (r *ReminderRepository) MarkRead(userID uint, targetType string, targetID uint) error {
	read := model.ReminderRead{
		UserID:     userID,
		TargetType: targetType,
		TargetID:   targetID,
	}
	return r.db.Where("user_id = ? AND target_type = ? AND target_id = ?",
		userID, targetType, targetID).FirstOrCreate(&read).Error
}

func (r *ReminderRepository) MarkAllRead(userID uint, items []model.ReminderRead) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&items).Error
}

func (r *ReminderRepository) GetDismissedKeys(userID uint) (map[string]bool, error) {
	var dismissed []model.ReminderDismissed
	if err := r.db.Where("user_id = ?", userID).Find(&dismissed).Error; err != nil {
		return nil, err
	}
	keys := make(map[string]bool, len(dismissed))
	for _, d := range dismissed {
		key := fmt.Sprintf("%s-%d", d.TargetType, d.TargetID)
		keys[key] = true
	}
	return keys, nil
}

func (r *ReminderRepository) Dismiss(userID uint, targetType string, targetID uint) error {
	d := model.ReminderDismissed{
		UserID:     userID,
		TargetType: targetType,
		TargetID:   targetID,
	}
	return r.db.Where("user_id = ? AND target_type = ? AND target_id = ?",
		userID, targetType, targetID).FirstOrCreate(&d).Error
}

func (r *ReminderRepository) GetSubscriptionsRenewingSoon(userID uint, withinDays int) ([]model.Subscription, error) {
	deadline := time.Now().AddDate(0, 0, withinDays)
	var subs []model.Subscription
	err := r.db.Where("user_id = ? AND status = ? AND next_payment_date IS NOT NULL AND next_payment_date <= ?",
		userID, "active", deadline).Find(&subs).Error
	return subs, err
}

func (r *ReminderRepository) GetAssetsExpiringSoon(userID uint, withinDays int) ([]model.Asset, error) {
	deadline := time.Now().AddDate(0, 0, withinDays)
	var assets []model.Asset
	err := r.db.Where("user_id = ? AND expire_date IS NOT NULL AND expire_date <= ?",
		userID, deadline).Find(&assets).Error
	return assets, err
}

func (r *ReminderRepository) GetWarningAssets(userID uint) ([]model.Asset, error) {
	var assets []model.Asset
	err := r.db.Where("user_id = ? AND status IN ?", userID, []string{"warning", "inactive"}).Find(&assets).Error
	return assets, err
}
