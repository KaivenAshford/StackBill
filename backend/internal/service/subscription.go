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

type SubscriptionService struct {
	repo *repository.SubscriptionRepository
}

func NewSubscriptionService(repo *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) List(userID uint, query *dto.SubscriptionListQuery) (*response.PageResult, error) {
	subs, total, err := s.repo.List(userID, query.Page, query.PageSize, query.CategoryID, query.Status, query.UpcomingRenewal)
	if err != nil {
		return nil, err
	}
	items := make([]dto.SubscriptionResponse, len(subs))
	for i, sub := range subs {
		items[i] = s.toResponse(&sub)
	}
	return &response.PageResult{Items: items, Total: total, Page: query.Page, PageSize: query.PageSize}, nil
}

func (s *SubscriptionService) GetByID(userID uint, id uint) (*dto.SubscriptionResponse, error) {
	sub, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewServiceError(404, 40400, "subscription not found")
		}
		return nil, err
	}
	if sub.UserID != userID {
		return nil, NewServiceError(403, 40301, "forbidden")
	}
	resp := s.toResponse(sub)
	return &resp, nil
}

func (s *SubscriptionService) Create(userID uint, req *dto.CreateSubscriptionRequest) (*dto.SubscriptionResponse, error) {
	sub := &model.Subscription{
		UserID:          userID,
		Name:            req.Name,
		Description:     req.Description,
		CategoryID:      req.CategoryID,
		Amount:          req.Amount,
		Currency:        req.Currency,
		BillingCycle:    req.BillingCycle,
		BillingInterval: req.BillingInterval,
		PaymentMethod:   req.PaymentMethod,
		Status:          "active",
		WebsiteURL:      req.WebsiteURL,
		Remark:          req.Remark,
		AutoRenew:       true,
	}
	if req.Status != "" {
		sub.Status = req.Status
	}
	if req.AutoRenew != nil {
		sub.AutoRenew = *req.AutoRenew
	}
	if req.StartDate != nil {
		t, _ := time.Parse("2006-01-02", *req.StartDate)
		sub.StartDate = &t
	}
	if sub.BillingInterval <= 0 {
		sub.BillingInterval = 1
	}
	sub.NextPaymentDate = s.calculateNextPayment(sub.StartDate, sub.BillingCycle, sub.BillingInterval)

	if err := s.repo.Create(sub); err != nil {
		return nil, err
	}
	resp := s.toResponse(sub)
	return &resp, nil
}

func (s *SubscriptionService) Update(userID uint, id uint, req *dto.UpdateSubscriptionRequest) (*dto.SubscriptionResponse, error) {
	sub, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewServiceError(404, 40400, "subscription not found")
		}
		return nil, err
	}
	if sub.UserID != userID {
		return nil, NewServiceError(403, 40301, "forbidden")
	}

	if req.Name != "" {
		sub.Name = req.Name
	}
	sub.Description = req.Description
	sub.CategoryID = req.CategoryID
	if req.Amount != 0 {
		sub.Amount = req.Amount
	}
	if req.Currency != "" {
		sub.Currency = req.Currency
	}
	if req.BillingCycle != "" {
		sub.BillingCycle = req.BillingCycle
	}
	if req.BillingInterval != 0 {
		sub.BillingInterval = req.BillingInterval
	}
	if req.StartDate != nil {
		t, _ := time.Parse("2006-01-02", *req.StartDate)
		sub.StartDate = &t
	}
	sub.PaymentMethod = req.PaymentMethod
	if req.AutoRenew != nil {
		sub.AutoRenew = *req.AutoRenew
	}
	if req.Status != "" {
		sub.Status = req.Status
	}
	sub.WebsiteURL = req.WebsiteURL
	sub.Remark = req.Remark

	sub.NextPaymentDate = s.calculateNextPayment(sub.StartDate, sub.BillingCycle, sub.BillingInterval)

	if err := s.repo.Update(sub); err != nil {
		return nil, err
	}
	resp := s.toResponse(sub)
	return &resp, nil
}

func (s *SubscriptionService) Delete(userID uint, id uint) error {
	sub, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewServiceError(404, 40400, "subscription not found")
		}
		return err
	}
	if sub.UserID != userID {
		return NewServiceError(403, 40301, "forbidden")
	}
	return s.repo.Delete(id)
}

func (s *SubscriptionService) CalculateMonthlyExpense(userID uint) (float64, error) {
	subs, err := s.repo.GetActiveByUserID(userID)
	if err != nil {
		return 0, err
	}
	var total float64
	for _, sub := range subs {
		total += s.normalizeToMonthly(sub.Amount, sub.BillingCycle, sub.BillingInterval)
	}
	return total, nil
}

func (s *SubscriptionService) CalculateYearlyExpense(userID uint) (float64, error) {
	subs, err := s.repo.GetActiveByUserID(userID)
	if err != nil {
		return 0, err
	}
	var total float64
	for _, sub := range subs {
		total += s.normalizeToYearly(sub.Amount, sub.BillingCycle, sub.BillingInterval)
	}
	return total, nil
}

func (s *SubscriptionService) normalizeToMonthly(amount float64, cycle string, interval int) float64 {
	if interval <= 0 {
		interval = 1
	}
	switch cycle {
	case "weekly":
		return amount * 4.33 / float64(interval)
	case "monthly":
		return amount / float64(interval)
	case "quarterly":
		return amount / 3.0 / float64(interval)
	case "yearly":
		return amount / 12.0 / float64(interval)
	case "one_time":
		return 0
	default:
		return amount / float64(interval)
	}
}

func (s *SubscriptionService) normalizeToYearly(amount float64, cycle string, interval int) float64 {
	if interval <= 0 {
		interval = 1
	}
	switch cycle {
	case "weekly":
		return amount * 52.0 / float64(interval)
	case "monthly":
		return amount * 12.0 / float64(interval)
	case "quarterly":
		return amount * 4.0 / float64(interval)
	case "yearly":
		return amount / float64(interval)
	case "one_time":
		return 0
	default:
		return amount * 12.0 / float64(interval)
	}
}

func (s *SubscriptionService) calculateNextPayment(startDate *time.Time, cycle string, interval int) *time.Time {
	if startDate == nil || cycle == "one_time" {
		return nil
	}
	now := time.Now()
	next := *startDate
	for next.Before(now) || next.Equal(now) {
		next = s.addCycle(next, cycle, interval)
		if next.Before(*startDate) || next.Equal(*startDate) {
			break
		}
	}
	return &next
}

func (s *SubscriptionService) addCycle(t time.Time, cycle string, interval int) time.Time {
	if interval <= 0 {
		interval = 1
	}
	switch cycle {
	case "weekly":
		return t.AddDate(0, 0, 7*interval)
	case "monthly":
		return t.AddDate(0, interval, 0)
	case "quarterly":
		return t.AddDate(0, 3*interval, 0)
	case "yearly":
		return t.AddDate(interval, 0, 0)
	case "custom":
		return t.AddDate(0, 0, interval)
	default:
		return t.AddDate(0, interval, 0)
	}
}

func (s *SubscriptionService) toResponse(sub *model.Subscription) dto.SubscriptionResponse {
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
