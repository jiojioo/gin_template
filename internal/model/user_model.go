package model

type User struct {
	BaseModel
	Username string `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Password string `gorm:"size:255;not null" json:"-"`
	Nickname string `gorm:"size:64" json:"nickname"`
	Status   int    `gorm:"not null;default:1" json:"status"`
}

func (User) TableName() string {
	return "users"
}
