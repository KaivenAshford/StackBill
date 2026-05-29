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

func (r *SubscriptionRepository) CountByUserID(userID uint) int64 {
	var count int64
	r.db.Model(&model.Subscription{}).Where("user_id = ?", userID).Count(&count)
	return count
}

func (r *SubscriptionRepository) GetRecentByUserID(userID uint, limit int) ([]model.Subscription, error) {
	var subs []model.Subscription
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Find(&subs).Error
	return subs, err
}

type CategoryExpenseRow struct {
	CategoryID   uint
	CategoryName string
	Color        string
	Amount       float64
}

func (r *SubscriptionRepository) GetCategoryExpense(userID uint) ([]CategoryExpenseRow, error) {
	var results []CategoryExpenseRow
	err := r.db.Raw(`
		SELECT c.id as category_id, c.name as category_name, c.color,
			COALESCE(SUM(
				CASE s.billing_cycle
					WHEN 'weekly' THEN s.amount * 4.33 / GREATEST(s.billing_interval, 1)
					WHEN 'monthly' THEN s.amount / GREATEST(s.billing_interval, 1)
					WHEN 'quarterly' THEN s.amount / 3.0 / GREATEST(s.billing_interval, 1)
					WHEN 'yearly' THEN s.amount / 12.0 / GREATEST(s.billing_interval, 1)
					ELSE 0
				END
			), 0) as amount
		FROM subscriptions s
		JOIN categories c ON s.category_id = c.id
		WHERE s.user_id = ? AND s.status = 'active' AND s.deleted_at IS NULL AND c.deleted_at IS NULL
		GROUP BY c.id, c.name, c.color
		ORDER BY amount DESC
	`, userID).Scan(&results).Error
	return results, err
}
