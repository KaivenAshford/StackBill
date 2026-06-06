package service

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/kingqaquuu/stackbill/internal/model"
	"github.com/kingqaquuu/stackbill/internal/repository"
)

type ImportService struct {
	subscriptionRepo *repository.SubscriptionRepository
	assetRepo        *repository.AssetRepository
}

func NewImportService(
	subscriptionRepo *repository.SubscriptionRepository,
	assetRepo *repository.AssetRepository,
) *ImportService {
	return &ImportService{
		subscriptionRepo: subscriptionRepo,
		assetRepo:        assetRepo,
	}
}

type ImportResult struct {
	Created int `json:"created"`
	Skipped int `json:"skipped"`
}

func (s *ImportService) ImportSubscriptionsCSV(userID uint, reader io.Reader) (*ImportResult, error) {
	r := csv.NewReader(reader)
	// Skip header row
	if _, err := r.Read(); err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}

	result := &ImportResult{}
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read row: %w", err)
		}
		if len(row) < 6 {
			result.Skipped++
			continue
		}

		amount, _ := strconv.ParseFloat(row[3], 64)
		interval, _ := strconv.Atoi(row[6])
		if interval <= 0 {
			interval = 1
		}
		categoryID, _ := strconv.ParseUint(row[2], 10, 64)

		sub := &model.Subscription{
			UserID:          userID,
			Name:            row[0],
			Description:     row[1],
			CategoryID:      uint(categoryID),
			Amount:          amount,
			Currency:        row[4],
			BillingCycle:    row[5],
			BillingInterval: interval,
			PaymentMethod:   col(row, 9),
			AutoRenew:       col(row, 10) != "false",
			Status:          colDefault(row, 11, "active"),
			WebsiteURL:      col(row, 12),
			Remark:          col(row, 13),
		}

		if sd := col(row, 7); sd != "" {
			t, _ := time.Parse("2006-01-02", sd)
			sub.StartDate = &t
		}
		if np := col(row, 8); np != "" {
			t, _ := time.Parse("2006-01-02", np)
			sub.NextPaymentDate = &t
		}

		if sub.Name == "" || sub.BillingCycle == "" {
			result.Skipped++
			continue
		}

		if err := s.subscriptionRepo.Create(sub); err != nil {
			result.Skipped++
			continue
		}
		result.Created++
	}

	return result, nil
}

func (s *ImportService) ImportAssetsCSV(userID uint, reader io.Reader) (*ImportResult, error) {
	r := csv.NewReader(reader)
	if _, err := r.Read(); err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}

	result := &ImportResult{}
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read row: %w", err)
		}
		if len(row) < 2 {
			result.Skipped++
			continue
		}

		costAmount, _ := strconv.ParseFloat(col(row, 6), 64)
		subID, _ := strconv.ParseUint(col(row, 12), 10, 64)

		asset := &model.Asset{
			UserID:         userID,
			Name:           row[0],
			AssetType:      row[1],
			Provider:       col(row, 2),
			Identifier:     col(row, 3),
			URL:            col(row, 4),
			CostAmount:     costAmount,
			CostCurrency:   colDefault(row, 7, "USD"),
			BillingCycle:   col(row, 8),
			Status:         colDefault(row, 9, "active"),
			Description:    col(row, 10),
			Remark:         col(row, 11),
			SubscriptionID: uint(subID),
		}

		if ed := col(row, 5); ed != "" {
			t, _ := time.Parse("2006-01-02", ed)
			asset.ExpireDate = &t
		}

		if asset.Name == "" || asset.AssetType == "" {
			result.Skipped++
			continue
		}

		if err := s.assetRepo.Create(asset); err != nil {
			result.Skipped++
			continue
		}
		result.Created++
	}

	return result, nil
}

func col(row []string, idx int) string {
	if idx < len(row) {
		return row[idx]
	}
	return ""
}

func colDefault(row []string, idx int, def string) string {
	v := col(row, idx)
	if v == "" {
		return def
	}
	return v
}
