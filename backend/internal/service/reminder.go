package service

import (
	"fmt"
	"time"

	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/model"
	"github.com/kingqaquuu/stackbill/internal/repository"
	"github.com/kingqaquuu/stackbill/pkg/response"
)

const (
	reminderOffsetRenewal = 0
	reminderOffsetExpiry  = 1_000_000
	reminderOffsetWarning = 2_000_000
)

type ReminderService struct {
	repo *repository.ReminderRepository
}

func NewReminderService(repo *repository.ReminderRepository) *ReminderService {
	return &ReminderService{repo: repo}
}

func (s *ReminderService) List(userID uint, query *dto.ReminderListQuery) (*response.PageResult, error) {
	var reminders []dto.ReminderResponse

	subs, _ := s.repo.GetSubscriptionsRenewingSoon(userID, 7)
	for _, sub := range subs {
		amount := sub.Amount
		reminders = append(reminders, dto.ReminderResponse{
			ID:         sub.ID + reminderOffsetRenewal,
			TargetType: "subscription",
			TargetID:   sub.ID,
			RemindType: "subscription_renewal",
			RemindDate: sub.NextPaymentDate.Format("2006-01-02"),
			Title:      sub.Name,
			Content:    fmt.Sprintf("Renewal on %s, amount %.2f %s", sub.NextPaymentDate.Format("2006-01-02"), sub.Amount, sub.Currency),
			Amount:     &amount,
			Currency:   sub.Currency,
		})
	}

	assets, _ := s.repo.GetAssetsExpiringSoon(userID, 30)
	for _, asset := range assets {
		reminders = append(reminders, dto.ReminderResponse{
			ID:         asset.ID + reminderOffsetExpiry,
			TargetType: "asset",
			TargetID:   asset.ID,
			RemindType: "asset_expiration",
			RemindDate: asset.ExpireDate.Format("2006-01-02"),
			Title:      asset.Name,
			Content:    fmt.Sprintf("Expiring on %s", asset.ExpireDate.Format("2006-01-02")),
			ExpireDate: asset.ExpireDate.Format("2006-01-02"),
		})
	}

	warnings, _ := s.repo.GetWarningAssets(userID)
	for _, asset := range warnings {
		reminders = append(reminders, dto.ReminderResponse{
			ID:          asset.ID + reminderOffsetWarning,
			TargetType:  "asset",
			TargetID:    asset.ID,
			RemindType:  "service_warning",
			RemindDate:  time.Now().Format("2006-01-02"),
			Title:       asset.Name,
			Content:     fmt.Sprintf("Status warning: %s", asset.Status),
			AssetStatus: asset.Status,
		})
	}

	if query.Type != "" {
		filtered := make([]dto.ReminderResponse, 0)
		for _, r := range reminders {
			if r.RemindType == query.Type {
				filtered = append(filtered, r)
			}
		}
		reminders = filtered
	}

	if query.IsRead != nil {
		filtered := make([]dto.ReminderResponse, 0)
		for _, r := range reminders {
			if r.IsRead == *query.IsRead {
				filtered = append(filtered, r)
			}
		}
		reminders = filtered
	}

	dismissedKeys, _ := s.repo.GetDismissedKeys(userID)
	if len(dismissedKeys) > 0 {
		filtered := make([]dto.ReminderResponse, 0, len(reminders))
		for _, r := range reminders {
			key := fmt.Sprintf("%s-%d", r.TargetType, r.TargetID)
			if !dismissedKeys[key] {
				filtered = append(filtered, r)
			}
		}
		reminders = filtered
	}

	total := int64(len(reminders))
	start := (query.Page - 1) * query.PageSize
	if start >= int(total) {
		return &response.PageResult{Items: []dto.ReminderResponse{}, Total: total, Page: query.Page, PageSize: query.PageSize}, nil
	}
	end := start + query.PageSize
	if end > int(total) {
		end = int(total)
	}

	return &response.PageResult{Items: reminders[start:end], Total: total, Page: query.Page, PageSize: query.PageSize}, nil
}

func (s *ReminderService) MarkRead(userID uint, id uint) error {
	targetType, targetID := s.decodeID(id)
	if targetType == "" {
		return NewServiceError(400, 40001, "invalid reminder id")
	}
	return s.repo.MarkRead(userID, targetType, targetID)
}

func (s *ReminderService) MarkAllRead(userID uint) error {
	var items []model.ReminderRead
	seen := make(map[string]bool)

	subs, _ := s.repo.GetSubscriptionsRenewingSoon(userID, 7)
	for _, sub := range subs {
		key := fmt.Sprintf("subscription-%d", sub.ID)
		if !seen[key] {
			seen[key] = true
			items = append(items, model.ReminderRead{UserID: userID, TargetType: "subscription", TargetID: sub.ID})
		}
	}

	assets, _ := s.repo.GetAssetsExpiringSoon(userID, 30)
	for _, asset := range assets {
		key := fmt.Sprintf("asset-%d", asset.ID)
		if !seen[key] {
			seen[key] = true
			items = append(items, model.ReminderRead{UserID: userID, TargetType: "asset", TargetID: asset.ID})
		}
	}

	warnings, _ := s.repo.GetWarningAssets(userID)
	for _, asset := range warnings {
		key := fmt.Sprintf("asset-%d", asset.ID)
		if !seen[key] {
			seen[key] = true
			items = append(items, model.ReminderRead{UserID: userID, TargetType: "asset", TargetID: asset.ID})
		}
	}

	if len(items) == 0 {
		return nil
	}
	return s.repo.MarkAllRead(userID, items)
}

func (s *ReminderService) Dismiss(userID uint, id uint) error {
	targetType, targetID := s.decodeID(id)
	if targetType == "" {
		return NewServiceError(400, 40001, "invalid reminder id")
	}
	return s.repo.Dismiss(userID, targetType, targetID)
}

func (s *ReminderService) decodeID(id uint) (string, uint) {
	switch {
	case id >= reminderOffsetWarning:
		return "asset", id - reminderOffsetWarning
	case id >= reminderOffsetExpiry:
		return "asset", id - reminderOffsetExpiry
	default:
		return "subscription", id
	}
}
