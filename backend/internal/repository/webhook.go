package repository

import (
	"github.com/kingqaquuu/stackbill/internal/model"
	"gorm.io/gorm"
)

type WebhookRepository struct {
	db *gorm.DB
}

func NewWebhookRepository(db *gorm.DB) *WebhookRepository {
	return &WebhookRepository{db: db}
}

func (r *WebhookRepository) List(userID uint) ([]model.Webhook, error) {
	var webhooks []model.Webhook
	err := r.db.Where("user_id = ?", userID).Find(&webhooks).Error
	return webhooks, err
}

func (r *WebhookRepository) FindByID(id uint) (*model.Webhook, error) {
	var webhook model.Webhook
	if err := r.db.First(&webhook, id).Error; err != nil {
		return nil, err
	}
	return &webhook, nil
}

func (r *WebhookRepository) Create(webhook *model.Webhook) error {
	return r.db.Create(webhook).Error
}

func (r *WebhookRepository) Update(webhook *model.Webhook) error {
	return r.db.Save(webhook).Error
}

func (r *WebhookRepository) Delete(id uint) error {
	return r.db.Delete(&model.Webhook{}, id).Error
}

func (r *WebhookRepository) GetActiveByUserID(userID uint) ([]model.Webhook, error) {
	var webhooks []model.Webhook
	err := r.db.Where("user_id = ? AND active = ?", userID, true).Find(&webhooks).Error
	return webhooks, err
}
