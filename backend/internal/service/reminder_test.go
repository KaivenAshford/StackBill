package service

import (
	"testing"

	"github.com/kingqaquuu/stackbill/internal/dto"
)

func TestReminderService_ListEmpty(t *testing.T) {
	svc := newTestServices(t)
	userID := registerTestUser(t, svc, "remuser")

	result, err := svc.ReminderService.List(userID, &dto.ReminderListQuery{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if result.Total != 0 {
		t.Errorf("total = %d, want 0", result.Total)
	}
}

func TestReminderService_ListWithRenewals(t *testing.T) {
	svc := newTestServices(t)
	userID := registerTestUser(t, svc, "remuser2")

	startDate := "2026-06-01"
	svc.SubscriptionService.Create(userID, &dto.CreateSubscriptionRequest{
		Name: "Renewing Soon", Amount: 10.00, Currency: "USD",
		BillingCycle: "weekly", StartDate: &startDate,
	})

	result, err := svc.ReminderService.List(userID, &dto.ReminderListQuery{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if result.Total < 1 {
		t.Error("expected at least 1 reminder")
	}
}

func TestReminderService_MarkRead(t *testing.T) {
	svc := newTestServices(t)
	userID := registerTestUser(t, svc, "remuser3")

	startDate := "2026-06-01"
	sub, _ := svc.SubscriptionService.Create(userID, &dto.CreateSubscriptionRequest{
		Name: "ReadTest", Amount: 5.00, Currency: "USD",
		BillingCycle: "weekly", StartDate: &startDate,
	})

	// Mark the renewal reminder as read (ID = sub.ID + reminderOffsetRenewal = sub.ID)
	err := svc.ReminderService.MarkRead(userID, sub.ID)
	if err != nil {
		t.Fatalf("mark read: %v", err)
	}
}
