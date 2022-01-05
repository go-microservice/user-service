package model

import (
	"sync"
	"time"
)

// UserProfileModel User represents a registered user.
type UserProfileModel struct {
	UserID    int64     `gorm:"column:user_id" json:"user_id"`
	Nickname  string    `gorm:"column:username" json:"nickname"`
	Avatar    string    `gorm:"column:avatar" json:"avatar"`
	Gender    string    `gorm:"column:gender" json:"gender"`
	Phone     string    `gorm:"column:phone" json:"phone"`
	Birthday  string    `gorm:"column:birthday" json:"birthday"`
	Bio       string    `gorm:"column:bio" json:"bio"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 表名
func (u *UserProfileModel) TableName() string {
	return "user_profile"
}