package model

type User struct {
	Model
	Username string `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email    string `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Password string `gorm:"size:255;not null" json:"-"`
	Nickname string `gorm:"size:50" json:"nickname"`
	Avatar   string `gorm:"size:500" json:"avatar"`
}
