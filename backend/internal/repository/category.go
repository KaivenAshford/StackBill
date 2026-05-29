package repository

import (
	"github.com/kingqaquuu/stackbill/internal/model"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) List(userID uint, categoryType string) ([]model.Category, error) {
	var categories []model.Category
	q := r.db.Where("user_id = ?", userID)
	if categoryType != "" {
		q = q.Where("type = ?", categoryType)
	}
	err := q.Order("sort_order ASC, id ASC").Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) FindByID(id uint) (*model.Category, error) {
	var cat model.Category
	if err := r.db.First(&cat, id).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *CategoryRepository) FindByName(userID uint, name string) (*model.Category, error) {
	var cat model.Category
	if err := r.db.Where("user_id = ? AND name = ?", userID, name).First(&cat).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *CategoryRepository) Create(cat *model.Category) error {
	return r.db.Create(cat).Error
}

func (r *CategoryRepository) Update(cat *model.Category) error {
	return r.db.Save(cat).Error
}

func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Delete(&model.Category{}, id).Error
}

func (r *CategoryRepository) BatchCreate(categories []model.Category) error {
	return r.db.Create(&categories).Error
}
