package service

import (
	"errors"
	"time"

	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/model"
	"github.com/kingqaquuu/stackbill/internal/repository"
	"github.com/kingqaquuu/stackbill/pkg/response"
	"gorm.io/gorm"
)

type AssetService struct {
	repo *repository.AssetRepository
}

func NewAssetService(repo *repository.AssetRepository) *AssetService {
	return &AssetService{repo: repo}
}

func (s *AssetService) List(userID uint, query *dto.AssetListQuery) (*response.PageResult, error) {
	assets, total, err := s.repo.List(userID, query.Page, query.PageSize, query.AssetType, query.Status, query.Keyword, query.ExpiringDays)
	if err != nil {
		return nil, err
	}
	items := make([]dto.AssetResponse, len(assets))
	for i, asset := range assets {
		items[i] = s.toResponse(&asset)
	}
	return &response.PageResult{Items: items, Total: total, Page: query.Page, PageSize: query.PageSize}, nil
}

func (s *AssetService) GetByID(userID uint, id uint) (*dto.AssetResponse, error) {
	asset, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewServiceError(404, ErrCodeNotFound, "asset not found")
		}
		return nil, err
	}
	if asset.UserID != userID {
		return nil, NewServiceError(403, ErrCodeForbidden, "forbidden")
	}
	resp := s.toResponse(asset)
	return &resp, nil
}

func (s *AssetService) Create(userID uint, req *dto.CreateAssetRequest) (*dto.AssetResponse, error) {
	asset := &model.Asset{
		UserID:         userID,
		Name:           req.Name,
		AssetType:      req.AssetType,
		Provider:       req.Provider,
		Identifier:     req.Identifier,
		URL:            req.URL,
		CostAmount:     req.CostAmount,
		CostCurrency:   req.CostCurrency,
		BillingCycle:   req.BillingCycle,
		SubscriptionID: req.SubscriptionID,
		Description:    req.Description,
		Remark:         req.Remark,
		Status:         "active",
	}
	if req.Status != "" {
		asset.Status = req.Status
	}
	if req.ExpireDate != nil {
		t, _ := time.Parse("2006-01-02", *req.ExpireDate)
		asset.ExpireDate = &t
	}
	if err := s.repo.Create(asset); err != nil {
		return nil, err
	}
	resp := s.toResponse(asset)
	return &resp, nil
}

func (s *AssetService) Update(userID uint, id uint, req *dto.UpdateAssetRequest) (*dto.AssetResponse, error) {
	asset, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewServiceError(404, ErrCodeNotFound, "asset not found")
		}
		return nil, err
	}
	if asset.UserID != userID {
		return nil, NewServiceError(403, ErrCodeForbidden, "forbidden")
	}

	if req.Name != "" {
		asset.Name = req.Name
	}
	if req.AssetType != "" {
		asset.AssetType = req.AssetType
	}
	asset.Provider = req.Provider
	asset.Identifier = req.Identifier
	asset.URL = req.URL
	if req.ExpireDate != nil {
		t, _ := time.Parse("2006-01-02", *req.ExpireDate)
		asset.ExpireDate = &t
	}
	asset.CostAmount = req.CostAmount
	asset.CostCurrency = req.CostCurrency
	asset.BillingCycle = req.BillingCycle
	if req.Status != "" {
		asset.Status = req.Status
	}
	asset.SubscriptionID = req.SubscriptionID
	asset.Description = req.Description
	asset.Remark = req.Remark

	if err := s.repo.Update(asset); err != nil {
		return nil, err
	}
	resp := s.toResponse(asset)
	return &resp, nil
}

func (s *AssetService) Delete(userID uint, id uint) error {
	asset, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewServiceError(404, ErrCodeNotFound, "asset not found")
		}
		return err
	}
	if asset.UserID != userID {
		return NewServiceError(403, ErrCodeForbidden, "forbidden")
	}
	return s.repo.Delete(id)
}

func (s *AssetService) toResponse(asset *model.Asset) dto.AssetResponse {
	var expireDate *string
	if asset.ExpireDate != nil {
		ed := asset.ExpireDate.Format("2006-01-02")
		expireDate = &ed
	}
	return dto.AssetResponse{
		ID:             asset.ID,
		Name:           asset.Name,
		AssetType:      asset.AssetType,
		Provider:       asset.Provider,
		Identifier:     asset.Identifier,
		URL:            asset.URL,
		ExpireDate:     expireDate,
		CostAmount:     asset.CostAmount,
		CostCurrency:   asset.CostCurrency,
		BillingCycle:   asset.BillingCycle,
		Status:         asset.Status,
		SubscriptionID: asset.SubscriptionID,
		Description:    asset.Description,
		Remark:         asset.Remark,
		CreatedAt:      asset.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      asset.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
