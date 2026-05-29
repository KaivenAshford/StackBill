package service

import (
	"errors"

	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/model"
	"github.com/kingqaquuu/stackbill/internal/repository"
	"gorm.io/gorm"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) List(userID uint, query *dto.CategoryListQuery) ([]dto.CategoryResponse, error) {
	categories, err := s.repo.List(userID, query.Type)
	if err != nil {
		return nil, err
	}
	result := make([]dto.CategoryResponse, len(categories))
	for i, cat := range categories {
		result[i] = s.toResponse(&cat)
	}
	return result, nil
}

func (s *CategoryService) GetByID(userID uint, id uint) (*dto.CategoryResponse, error) {
	cat, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewServiceError(404, 40400, "category not found")
		}
		return nil, err
	}
	if cat.UserID != userID {
		return nil, NewServiceError(403, 40301, "forbidden")
	}
	resp := s.toResponse(cat)
	return &resp, nil
}

func (s *CategoryService) Create(userID uint, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	if _, err := s.repo.FindByName(userID, req.Name); err == nil {
		return nil, NewServiceError(409, 40901, "category name already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	cat := &model.Category{
		UserID:    userID,
		Name:      req.Name,
		Type:      req.Type,
		Color:     req.Color,
		Icon:      req.Icon,
		SortOrder: req.SortOrder,
	}
	if err := s.repo.Create(cat); err != nil {
		return nil, err
	}
	resp := s.toResponse(cat)
	return &resp, nil
}

func (s *CategoryService) Update(userID uint, id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	cat, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewServiceError(404, 40400, "category not found")
		}
		return nil, err
	}
	if cat.UserID != userID {
		return nil, NewServiceError(403, 40301, "forbidden")
	}

	cat.Name = req.Name
	cat.Type = req.Type
	cat.Color = req.Color
	cat.Icon = req.Icon
	cat.SortOrder = req.SortOrder

	if err := s.repo.Update(cat); err != nil {
		return nil, err
	}
	resp := s.toResponse(cat)
	return &resp, nil
}

func (s *CategoryService) Delete(userID uint, id uint) error {
	cat, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewServiceError(404, 40400, "category not found")
		}
		return err
	}
	if cat.UserID != userID {
		return NewServiceError(403, 40301, "forbidden")
	}
	return s.repo.Delete(id)
}

func (s *CategoryService) toResponse(cat *model.Category) dto.CategoryResponse {
	return dto.CategoryResponse{
		ID:        cat.ID,
		Name:      cat.Name,
		Type:      cat.Type,
		Color:     cat.Color,
		Icon:      cat.Icon,
		SortOrder: cat.SortOrder,
	}
}
