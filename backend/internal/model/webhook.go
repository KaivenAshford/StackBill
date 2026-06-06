package model

type Webhook struct {
	Model
	UserID  uint   `gorm:"index;not null" json:"user_id"`
	URL     string `gorm:"size:500;not null" json:"url"`
	Secret  string `gorm:"size:200" json:"secret"`
	Events  string `gorm:"size:500;not null" json:"events"`
	Active  bool   `gorm:"default:true" json:"active"`
}
