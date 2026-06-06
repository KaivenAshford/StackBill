package service

import (
	"testing"

	"github.com/kingqaquuu/stackbill/internal/dto"
)

func TestSubscriptionService_CreateAndGet(t *testing.T) {
	svc := newTestServices(t)
	userID := registerTestUser(t, svc, "subuser")

	startDate := "2026-01-01"
	resp, err := svc.SubscriptionService.Create(userID, &dto.CreateSubscriptionRequest{
		Name: "Netflix", Amount: 15.99, Currency: "USD",
		BillingCycle: "monthly", StartDate: &startDate,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if resp.Name != "Netflix" {
		t.Errorf("name = %q, want %q", resp.Name, "Netflix")
	}
	if resp.Status != "active" {
		t.Errorf("status = %q, want %q", resp.Status, "active")
	}

	got, err := svc.SubscriptionService.GetByID(userID, resp.ID)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if got.Name != "Netflix" {
		t.Errorf("got name = %q, want %q", got.Name, "Netflix")
	}
}

func TestSubscriptionService_Update(t *testing.T) {
	svc := newTestServices(t)
	userID := registerTestUser(t, svc, "subuser2")

	created, _ := svc.SubscriptionService.Create(userID, &dto.CreateSubscriptionRequest{
		Name: "Spotify", Amount: 9.99, Currency: "USD", BillingCycle: "monthly",
	})

	updated, err := svc.SubscriptionService.Update(userID, created.ID, &dto.UpdateSubscriptionRequest{
		Name: "Spotify Premium", Amount: 14.99,
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Name != "Spotify Premium" {
		t.Errorf("name = %q, want %q", updated.Name, "Spotify Premium")
	}
	if updated.Amount != 14.99 {
		t.Errorf("amount = %f, want %f", updated.Amount, 14.99)
	}
}

func TestSubscriptionService_Delete(t *testing.T) {
	svc := newTestServices(t)
	userID := registerTestUser(t, svc, "subuser3")

	created, _ := svc.SubscriptionService.Create(userID, &dto.CreateSubscriptionRequest{
		Name: "ToDelete", Amount: 5.00, Currency: "USD", BillingCycle: "yearly",
	})

	err := svc.SubscriptionService.Delete(userID, created.ID)
	if err != nil {
		t.Fatalf("delete: %v", err)
	}

	_, err = svc.SubscriptionService.GetByID(userID, created.ID)
	if err == nil {
		t.Error("expected error getting deleted subscription")
	}
}

func TestSubscriptionService_UserIsolation(t *testing.T) {
	svc := newTestServices(t)
	userA := registerTestUser(t, svc, "userA")
	userB := registerTestUser(t, svc, "userB")

	created, _ := svc.SubscriptionService.Create(userA, &dto.CreateSubscriptionRequest{
		Name: "Private Sub", Amount: 10.00, Currency: "USD", BillingCycle: "monthly",
	})

	_, err := svc.SubscriptionService.GetByID(userB, created.ID)
	if err == nil {
		t.Error("user B should not access user A's subscription")
	}

	err = svc.SubscriptionService.Delete(userB, created.ID)
	if err == nil {
		t.Error("user B should not delete user A's subscription")
	}
}

func TestSubscriptionService_CalculateMonthlyExpense(t *testing.T) {
	svc := newTestServices(t)
	userID := registerTestUser(t, svc, "expenseuser")

	svc.SubscriptionService.Create(userID, &dto.CreateSubscriptionRequest{
		Name: "Monthly Sub", Amount: 10.00, Currency: "USD", BillingCycle: "monthly",
	})
	svc.SubscriptionService.Create(userID, &dto.CreateSubscriptionRequest{
		Name: "Yearly Sub", Amount: 120.00, Currency: "USD", BillingCycle: "yearly",
	})

	monthly, err := svc.SubscriptionService.CalculateMonthlyExpense(userID)
	if err != nil {
		t.Fatalf("monthly expense: %v", err)
	}
	// 10.00 (monthly) + 10.00 (120/12 yearly) = 20.00
	if monthly != 20.00 {
		t.Errorf("monthly expense = %f, want %f", monthly, 20.00)
	}
}
