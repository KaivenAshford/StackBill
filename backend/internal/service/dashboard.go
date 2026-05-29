package service

import (
	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/model"
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
	monthlyExpense, _ := s.subscriptionSvc.CalculateMonthlyExpense(userID)
	yearlyExpense, _ := s.subscriptionSvc.CalculateYearlyExpense(userID)

	subCount := s.subscriptionRepo.CountByUserID(userID)
	assetCount := s.assetRepo.CountByUserID(userID)

	renewingSubs, _ := s.reminderRepo.GetSubscriptionsRenewingSoon(userID, 7)
	expiringAssets, _ := s.reminderRepo.GetAssetsExpiringSoon(userID, 30)
	warningAssets, _ := s.reminderRepo.GetWarningAssets(userID)

	recentSubs, _ := s.subscriptionRepo.GetRecentByUserID(userID, 5)
	recentAssets, _ := s.assetRepo.GetRecentByUserID(userID, 5)

	categoryExpenseRows, _ := s.subscriptionRepo.GetCategoryExpense(userID)

	recentSubResponses := make([]dto.SubscriptionResponse, len(recentSubs))
	for i := range recentSubs {
		recentSubResponses[i] = subToResponse(&recentSubs[i])
	}

	recentAssetResponses := make([]dto.AssetResponse, len(recentAssets))
	for i := range recentAssets {
		recentAssetResponses[i] = assetToResponse(&recentAssets[i])
	}

	renewalResponses := make([]dto.SubscriptionResponse, len(renewingSubs))
	for i := range renewingSubs {
		renewalResponses[i] = subToResponse(&renewingSubs[i])
	}

	expiringResponses := make([]dto.AssetResponse, len(expiringAssets))
	for i := range expiringAssets {
		expiringResponses[i] = assetToResponse(&expiringAssets[i])
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

func subToResponse(sub *model.Subscription) dto.SubscriptionResponse {
	var nextPayment, startDate *string
	if sub.NextPaymentDate != nil {
		v := sub.NextPaymentDate.Format("2006-01-02")
		nextPayment = &v
	}
	if sub.StartDate != nil {
		v := sub.StartDate.Format("2006-01-02")
		startDate = &v
	}
	return dto.SubscriptionResponse{
		ID:              sub.ID,
		Name:            sub.Name,
		Description:     sub.Description,
		CategoryID:      sub.CategoryID,
		Amount:          sub.Amount,
		Currency:        sub.Currency,
		BillingCycle:    sub.BillingCycle,
		BillingInterval: sub.BillingInterval,
		NextPaymentDate: nextPayment,
		StartDate:       startDate,
		PaymentMethod:   sub.PaymentMethod,
		AutoRenew:       sub.AutoRenew,
		Status:          sub.Status,
		WebsiteURL:      sub.WebsiteURL,
		Remark:          sub.Remark,
		CreatedAt:       sub.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       sub.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func assetToResponse(asset *model.Asset) dto.AssetResponse {
	var expireDate *string
	if asset.ExpireDate != nil {
		ed := asset.ExpireDate.Format("2006-01-02")
		expireDate = &ed
	}
	return dto.AssetResponse{
		ID:           asset.ID,
		Name:         asset.Name,
		AssetType:    asset.AssetType,
		Provider:     asset.Provider,
		Identifier:   asset.Identifier,
		URL:          asset.URL,
		ExpireDate:   expireDate,
		CostAmount:   asset.CostAmount,
		CostCurrency: asset.CostCurrency,
		BillingCycle: asset.BillingCycle,
		Status:       asset.Status,
		Description:  asset.Description,
		Remark:       asset.Remark,
		CreatedAt:    asset.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    asset.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
