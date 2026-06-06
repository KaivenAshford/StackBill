package repository

import (
	"time"

	"github.com/kingqaquuu/stackbill/internal/model"
	"gorm.io/gorm"
)

type AssetRepository struct {
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) *AssetRepository {
	return &AssetRepository{db: db}
}

func (r *AssetRepository) List(userID uint, page, pageSize int, assetType, status, keyword string, expiringDays int) ([]model.Asset, int64, error) {
	q := r.db.Where("user_id = ?", userID)
	if assetType != "" {
		q = q.Where("asset_type = ?", assetType)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if keyword != "" {
		q = q.Where("name ILIKE ?", "%"+keyword+"%")
	}
	if expiringDays > 0 {
		deadline := time.Now().AddDate(0, 0, expiringDays)
		q = q.Where("expire_date <= ? AND expire_date IS NOT NULL", deadline)
	}
	var total int64
	if err := q.Model(&model.Asset{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var assets []model.Asset
	offset := (page - 1) * pageSize
	err := q.Order("expire_date ASC, id DESC").Offset(offset).Limit(pageSize).Find(&assets).Error
	return assets, total, err
}

func (r *AssetRepository) FindByID(id uint) (*model.Asset, error) {
	var asset model.Asset
	if err := r.db.First(&asset, id).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *AssetRepository) Create(asset *model.Asset) error {
	return r.db.Create(asset).Error
}

func (r *AssetRepository) Update(asset *model.Asset) error {
	return r.db.Save(asset).Error
}

func (r *AssetRepository) Delete(id uint) error {
	return r.db.Delete(&model.Asset{}, id).Error
}

func (r *AssetRepository) GetByUserID(userID uint) ([]model.Asset, error) {
	var assets []model.Asset
	err := r.db.Where("user_id = ?", userID).Find(&assets).Error
	return assets, err
}

func (r *AssetRepository) CountByUserID(userID uint) int64 {
	var count int64
	r.db.Model(&model.Asset{}).Where("user_id = ?", userID).Count(&count)
	return count
}

func (r *AssetRepository) GetRecentByUserID(userID uint, limit int) ([]model.Asset, error) {
	var assets []model.Asset
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Find(&assets).Error
	return assets, err
}
