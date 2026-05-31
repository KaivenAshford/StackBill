# StackBill Phase 4: Dashboard Design

## Overview

Single dashboard API that aggregates statistics from existing services, plus a frontend dashboard page with cards, ECharts chart, and quick-glance lists.

## API

**GET `/api/v1/dashboard`** (JWT protected)

Response:
```json
{
  "code": 0,
  "data": {
    "monthly_expense": 120.50,
    "yearly_expense": 1446.00,
    "subscription_count": 12,
    "asset_count": 8,
    "upcoming_renewals": 3,
    "expiring_assets": 2,
    "warning_assets": 1,
    "recent_subscriptions": [...],
    "recent_assets": [...],
    "upcoming_renewal_list": [...],
    "expiring_asset_list": [...],
    "category_expense": [
      {"category_id": 1, "category_name": "AI 工具", "amount": 50.00, "color": "#7c3aed"},
      {"category_id": 2, "category_name": "开发工具", "amount": 30.00, "color": "#2563eb"}
    ]
  }
}
```

`recent_subscriptions` and `recent_assets` return top 5 newest items (same DTOs as existing list responses).
`upcoming_renewal_list` returns top 5 subscriptions with next_payment_date within 7 days.
`expiring_asset_list` returns top 5 assets with expire_date within 30 days.
`category_expense` groups active subscriptions by category_id, sums normalized monthly amounts.

## Backend Files

**Create:**
- `internal/service/dashboard.go` — DashboardService
- `internal/api/dashboard.go` — DashboardHandler
- `internal/dto/dashboard.go` — DashboardResponse DTO

**Modify:**
- `internal/router/router.go` — add dashboard route
- `internal/repository/subscription.go` — add CountByUserID, GetRecentByUserID, GetCategoryExpense methods
- `internal/repository/asset.go` — add CountByUserID, GetRecentByUserID methods

## DTO

```go
type DashboardResponse struct {
    MonthlyExpense     float64                  `json:"monthly_expense"`
    YearlyExpense      float64                  `json:"yearly_expense"`
    SubscriptionCount  int64                    `json:"subscription_count"`
    AssetCount         int64                    `json:"asset_count"`
    UpcomingRenewals   int                      `json:"upcoming_renewals"`
    ExpiringAssets     int                      `json:"expiring_assets"`
    WarningAssets      int                      `json:"warning_assets"`
    RecentSubscriptions []SubscriptionResponse  `json:"recent_subscriptions"`
    RecentAssets       []AssetResponse          `json:"recent_assets"`
    UpcomingRenewalList []SubscriptionResponse  `json:"upcoming_renewal_list"`
    ExpiringAssetList  []AssetResponse          `json:"expiring_asset_list"`
    CategoryExpense    []CategoryExpenseItem    `json:"category_expense"`
}

type CategoryExpenseItem struct {
    CategoryID   uint    `json:"category_id"`
    CategoryName string  `json:"category_name"`
    Amount       float64 `json:"amount"`
    Color        string  `json:"color"`
}
```

## Service

DashboardService depends on SubscriptionService, AssetService, SubscriptionRepository, AssetRepository, ReminderRepository.

```go
func NewDashboardService(...) *DashboardService
func (s *DashboardService) GetDashboard(userID uint) (*dto.DashboardResponse, error)
```

Implementation aggregates calls to existing services/repos.

## Frontend

**Create:**
- `src/api/dashboard.ts` — getDashboard() function

**Modify:**
- `src/views/dashboard/Index.vue` — replace with stats cards + ECharts pie chart + upcoming lists
- `src/layouts/MainLayout.vue` — add onMounted(() => store.fetchUser())
- `src/types/index.ts` — expand DashboardStats to match full DTO

## Dashboard Page Layout

```
┌─────────┬─────────┬─────────┬─────────┐
│ Monthly │ Yearly  │ Subs    │ Assets  │
│ Expense │ Expense │ Count   │ Count   │
└─────────┴─────────┴─────────┴─────────┘
┌──────────────────┬──────────────────────┐
│ Category Expense │ Upcoming Renewals    │
│ (ECharts Pie)    │ (list)               │
├──────────────────┼───────────────────── │
│ Warning Assets   │ Expiring Assets      │
│ (count badge)    │ (list)               │
└──────────────────┴──────────────────────┘
```

PC: 4 stat cards on top, 2 columns below (chart left, lists right).
Mobile layout is out of scope for the current plan and should be designed only when requested.
