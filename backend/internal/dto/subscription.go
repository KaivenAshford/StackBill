package dto

type SubscriptionResponse struct {
	ID              uint    `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	CategoryID      uint    `json:"category_id"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	BillingCycle    string  `json:"billing_cycle"`
	BillingInterval int     `json:"billing_interval"`
	NextPaymentDate *string `json:"next_payment_date"`
	StartDate       *string `json:"start_date"`
	PaymentMethod   string  `json:"payment_method"`
	AutoRenew       bool    `json:"auto_renew"`
	Status          string  `json:"status"`
	WebsiteURL      string  `json:"website_url"`
	Remark          string  `json:"remark"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

type CreateSubscriptionRequest struct {
	Name            string  `json:"name" binding:"required,max=100"`
	Description     string  `json:"description" binding:"max=500"`
	CategoryID      uint    `json:"category_id"`
	Amount          float64 `json:"amount" binding:"required"`
	Currency        string  `json:"currency" binding:"max=10"`
	BillingCycle    string  `json:"billing_cycle" binding:"required,oneof=weekly monthly quarterly yearly custom one_time"`
	BillingInterval int     `json:"billing_interval"`
	StartDate       *string `json:"start_date"`
	PaymentMethod   string  `json:"payment_method" binding:"max=50"`
	AutoRenew       *bool   `json:"auto_renew"`
	Status          string  `json:"status" binding:"omitempty,oneof=active paused cancelled expired"`
	WebsiteURL      string  `json:"website_url" binding:"max=500"`
	Remark          string  `json:"remark" binding:"max=500"`
}

type UpdateSubscriptionRequest struct {
	Name            string  `json:"name" binding:"max=100"`
	Description     string  `json:"description" binding:"max=500"`
	CategoryID      uint    `json:"category_id"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency" binding:"max=10"`
	BillingCycle    string  `json:"billing_cycle" binding:"oneof=weekly monthly quarterly yearly custom one_time"`
	BillingInterval int     `json:"billing_interval"`
	StartDate       *string `json:"start_date"`
	PaymentMethod   string  `json:"payment_method" binding:"max=50"`
	AutoRenew       *bool   `json:"auto_renew"`
	Status          string  `json:"status" binding:"omitempty,oneof=active paused cancelled expired"`
	WebsiteURL      string  `json:"website_url" binding:"max=500"`
	Remark          string  `json:"remark" binding:"max=500"`
}

type SubscriptionListQuery struct {
	Page            int    `form:"page,default=1"`
	PageSize        int    `form:"page_size,default=20"`
	CategoryID      *uint  `form:"category_id"`
	Status          string `form:"status"`
	UpcomingRenewal bool   `form:"upcoming_renewal"`
}
