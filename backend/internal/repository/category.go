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

func (r *CategoryRepository) BatchCreate(categories []model.Category) error {
	return r.db.Create(&categories).Error
}
