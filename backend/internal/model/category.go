package model

type Category struct {
	Model
	UserID     uint   `gorm:"index;not null" json:"user_id"`
	Name       string `gorm:"size:50;not null" json:"name"`
	Type       string `gorm:"size:20;not null" json:"type"`
	Color      string `gorm:"size:20" json:"color"`
	Icon       string `gorm:"size:50" json:"icon"`
	SortOrder  int    `gorm:"default:0" json:"sort_order"`
}
