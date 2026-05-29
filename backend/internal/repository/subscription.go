package repository

import (
	"time"

	"github.com/kingqaquuu/stackbill/internal/model"
	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) List(userID uint, page, pageSize int, categoryID *uint, status string, upcomingRenewal bool) ([]model.Subscription, int64, error) {
	q := r.db.Where("user_id = ?", userID)
	if categoryID != nil {
		q = q.Where("category_id = ?", *categoryID)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if upcomingRenewal {
		sevenDays := time.Now().Add(7 * 24 * time.Hour)
		q = q.Where("next_payment_date <= ? AND next_payment_date IS NOT NULL", sevenDays)
	}
	var total int64
	if err := q.Model(&model.Subscription{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var subs []model.Subscription
	offset := (page - 1) * pageSize
	err := q.Order("next_payment_date ASC, id DESC").Offset(offset).Limit(pageSize).Find(&subs).Error
	return subs, total, err
}

func (r *SubscriptionRepository) FindByID(id uint) (*model.Subscription, error) {
	var sub model.Subscription
	if err := r.db.First(&sub, id).Error; err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *SubscriptionRepository) Create(sub *model.Subscription) error {
	return r.db.Create(sub).Error
}

func (r *SubscriptionRepository) Update(sub *model.Subscription) error {
	return r.db.Save(sub).Error
}

func (r *SubscriptionRepository) Delete(id uint) error {
	return r.db.Delete(&model.Subscription{}, id).Error
}

func (r *SubscriptionRepository) GetActiveByUserID(userID uint) ([]model.Subscription, error) {
	var subs []model.Subscription
	err := r.db.Where("user_id = ? AND status = ?", userID, "active").Find(&subs).Error
	return subs, err
}
