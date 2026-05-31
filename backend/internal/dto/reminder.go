package dto

type ReminderResponse struct {
	ID          uint     `json:"id"`
	TargetType  string   `json:"target_type"`
	TargetID    uint     `json:"target_id"`
	RemindType  string   `json:"remind_type"`
	RemindDate  string   `json:"remind_date"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	IsRead      bool     `json:"is_read"`
	Amount      *float64 `json:"amount,omitempty"`
	Currency    string   `json:"currency,omitempty"`
	ExpireDate  string   `json:"expire_date,omitempty"`
	AssetStatus string   `json:"asset_status,omitempty"`
}

type ReminderListQuery struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=20"`
	Type     string `form:"type"`
	IsRead   *bool  `form:"is_read"`
}
