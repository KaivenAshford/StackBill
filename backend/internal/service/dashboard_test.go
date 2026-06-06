package service

import (
	"testing"

	"github.com/kingqaquuu/stackbill/internal/dto"
)

func TestDashboardService_GetDashboard(t *testing.T) {
	svc := newTestServices(t)
	userID := registerTestUser(t, svc, "dashuser")

	svc.SubscriptionService.Create(userID, &dto.CreateSubscriptionRequest{
		Name: "Test Sub", Amount: 10.00, Currency: "USD", BillingCycle: "monthly",
	})
	svc.AssetService.Create(userID, &dto.CreateAssetRequest{
		Name: "Test Asset", AssetType: "domain",
	})

	dash, err := svc.DashboardService.GetDashboard(userID)
	if err != nil {
		t.Fatalf("get dashboard: %v", err)
	}
	if dash.SubscriptionCount != 1 {
		t.Errorf("subscription count = %d, want 1", dash.SubscriptionCount)
	}
	if dash.AssetCount != 1 {
		t.Errorf("asset count = %d, want 1", dash.AssetCount)
	}
	if dash.MonthlyExpense != 10.00 {
		t.Errorf("monthly expense = %f, want 10.00", dash.MonthlyExpense)
	}
}

func TestDashboardService_UserIsolation(t *testing.T) {
	svc := newTestServices(t)
	userA := registerTestUser(t, svc, "dashA")
	userB := registerTestUser(t, svc, "dashB")

	svc.SubscriptionService.Create(userA, &dto.CreateSubscriptionRequest{
		Name: "Private", Amount: 100.00, Currency: "USD", BillingCycle: "monthly",
	})

	dashB, _ := svc.DashboardService.GetDashboard(userB)
	if dashB.SubscriptionCount != 0 {
		t.Errorf("user B subscription count = %d, want 0", dashB.SubscriptionCount)
	}
	if dashB.MonthlyExpense != 0 {
		t.Errorf("user B monthly expense = %f, want 0", dashB.MonthlyExpense)
	}
}
