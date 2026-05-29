package dto

type DashboardResponse struct {
	MonthlyExpense      float64                `json:"monthly_expense"`
	YearlyExpense       float64                `json:"yearly_expense"`
	SubscriptionCount   int64                  `json:"subscription_count"`
	AssetCount          int64                  `json:"asset_count"`
	UpcomingRenewals    int                    `json:"upcoming_renewals"`
	ExpiringAssets      int                    `json:"expiring_assets"`
	WarningAssets       int                    `json:"warning_assets"`
	RecentSubscriptions []SubscriptionResponse `json:"recent_subscriptions"`
	RecentAssets        []AssetResponse        `json:"recent_assets"`
	UpcomingRenewalList []SubscriptionResponse `json:"upcoming_renewal_list"`
	ExpiringAssetList   []AssetResponse        `json:"expiring_asset_list"`
	CategoryExpense     []CategoryExpenseItem  `json:"category_expense"`
}

type CategoryExpenseItem struct {
	CategoryID   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Amount       float64 `json:"amount"`
	Color        string  `json:"color"`
}
