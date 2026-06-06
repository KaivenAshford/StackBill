package service

import (
	"log/slog"

	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/repository"
)

type DashboardService struct {
	subscriptionRepo *repository.SubscriptionRepository
	assetRepo        *repository.AssetRepository
	reminderRepo     *repository.ReminderRepository
	categoryRepo     *repository.CategoryRepository
	subscriptionSvc  *SubscriptionService
}

func NewDashboardService(
	subscriptionRepo *repository.SubscriptionRepository,
	assetRepo *repository.AssetRepository,
	reminderRepo *repository.ReminderRepository,
	categoryRepo *repository.CategoryRepository,
	subscriptionSvc *SubscriptionService,
) *DashboardService {
	return &DashboardService{
		subscriptionRepo: subscriptionRepo,
		assetRepo:        assetRepo,
		reminderRepo:     reminderRepo,
		categoryRepo:     categoryRepo,
		subscriptionSvc:  subscriptionSvc,
	}
}

func (s *DashboardService) GetDashboard(userID uint) (*dto.DashboardResponse, error) {
	monthlyExpense, err := s.subscriptionSvc.CalculateMonthlyExpense(userID)
	if err != nil {
		slog.Error("failed to calculate monthly expense", "user_id", userID, "error", err)
	}
	yearlyExpense, err := s.subscriptionSvc.CalculateYearlyExpense(userID)
	if err != nil {
		slog.Error("failed to calculate yearly expense", "user_id", userID, "error", err)
	}

	subCount := s.subscriptionRepo.CountByUserID(userID)
	assetCount := s.assetRepo.CountByUserID(userID)

	renewingSubs, err := s.reminderRepo.GetSubscriptionsRenewingSoon(userID, 7)
	if err != nil {
		slog.Error("failed to get renewing subscriptions", "user_id", userID, "error", err)
	}
	expiringAssets, err := s.reminderRepo.GetAssetsExpiringSoon(userID, 30)
	if err != nil {
		slog.Error("failed to get expiring assets", "user_id", userID, "error", err)
	}
	warningAssets, err := s.reminderRepo.GetWarningAssets(userID)
	if err != nil {
		slog.Error("failed to get warning assets", "user_id", userID, "error", err)
	}

	recentSubs, err := s.subscriptionRepo.GetRecentByUserID(userID, 5)
	if err != nil {
		slog.Error("failed to get recent subscriptions", "user_id", userID, "error", err)
	}
	recentAssets, err := s.assetRepo.GetRecentByUserID(userID, 5)
	if err != nil {
		slog.Error("failed to get recent assets", "user_id", userID, "error", err)
	}

	categoryExpenseRows, err := s.subscriptionRepo.GetCategoryExpense(userID)
	if err != nil {
		slog.Error("failed to get category expense", "user_id", userID, "error", err)
	}

	recentSubResponses := make([]dto.SubscriptionResponse, len(recentSubs))
	for i := range recentSubs {
		recentSubResponses[i] = SubscriptionToResponse(&recentSubs[i])
	}

	recentAssetResponses := make([]dto.AssetResponse, len(recentAssets))
	for i := range recentAssets {
		recentAssetResponses[i] = AssetToResponse(&recentAssets[i])
	}

	renewalResponses := make([]dto.SubscriptionResponse, len(renewingSubs))
	for i := range renewingSubs {
		renewalResponses[i] = SubscriptionToResponse(&renewingSubs[i])
	}

	expiringResponses := make([]dto.AssetResponse, len(expiringAssets))
	for i := range expiringAssets {
		expiringResponses[i] = AssetToResponse(&expiringAssets[i])
	}

	categoryExpense := make([]dto.CategoryExpenseItem, len(categoryExpenseRows))
	for i, row := range categoryExpenseRows {
		categoryExpense[i] = dto.CategoryExpenseItem{
			CategoryID:   row.CategoryID,
			CategoryName: row.CategoryName,
			Amount:       row.Amount,
			Color:        row.Color,
		}
	}

	return &dto.DashboardResponse{
		MonthlyExpense:      monthlyExpense,
		YearlyExpense:       yearlyExpense,
		SubscriptionCount:   subCount,
		AssetCount:          assetCount,
		UpcomingRenewals:    len(renewingSubs),
		ExpiringAssets:      len(expiringAssets),
		WarningAssets:       len(warningAssets),
		RecentSubscriptions: recentSubResponses,
		RecentAssets:        recentAssetResponses,
		UpcomingRenewalList: renewalResponses,
		ExpiringAssetList:   expiringResponses,
		CategoryExpense:     categoryExpense,
	}, nil
}
