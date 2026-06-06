package model

import "time"

type Asset struct {
	Model
	UserID         uint       `gorm:"index;not null" json:"user_id"`
	Name           string     `gorm:"size:100;not null" json:"name"`
	AssetType      string     `gorm:"size:30;not null" json:"asset_type"`
	Provider       string     `gorm:"size:100" json:"provider"`
	Identifier     string     `gorm:"size:200" json:"identifier"`
	URL            string     `gorm:"size:500" json:"url"`
	ExpireDate     *time.Time `json:"expire_date"`
	CostAmount     float64    `gorm:"type:decimal(10,2)" json:"cost_amount"`
	CostCurrency   string     `gorm:"size:10;default:USD" json:"cost_currency"`
	BillingCycle   string     `gorm:"size:20" json:"billing_cycle"`
	Status         string     `gorm:"size:20;default:active" json:"status"`
	SubscriptionID uint       `gorm:"index" json:"subscription_id"`
	Description    string     `gorm:"size:500" json:"description"`
	Remark         string     `gorm:"size:500" json:"remark"`
}
