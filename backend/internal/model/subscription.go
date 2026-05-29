package model

import "time"

type Subscription struct {
	Model
	UserID           uint      `gorm:"index;not null" json:"user_id"`
	Name             string    `gorm:"size:100;not null" json:"name"`
	Description      string    `gorm:"size:500" json:"description"`
	CategoryID       uint      `gorm:"index" json:"category_id"`
	Amount           float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Currency         string    `gorm:"size:10;default:USD" json:"currency"`
	BillingCycle     string    `gorm:"size:20;not null" json:"billing_cycle"`
	BillingInterval  int       `gorm:"default:1" json:"billing_interval"`
	NextPaymentDate  *time.Time `json:"next_payment_date"`
	StartDate        *time.Time `json:"start_date"`
	PaymentMethod    string    `gorm:"size:50" json:"payment_method"`
	AutoRenew        bool      `gorm:"default:true" json:"auto_renew"`
	Status           string    `gorm:"size:20;default:active" json:"status"`
	WebsiteURL       string    `gorm:"size:500" json:"website_url"`
	Remark           string    `gorm:"size:500" json:"remark"`
}
