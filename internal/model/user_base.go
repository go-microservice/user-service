package model

import (
	"sync"
	"time"
)

// UserBaseModel User represents a registered user.
type UserBaseModel struct {
	ID        int64     `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	Username  string    `gorm:"column:username" json:"username"`
	Password  string    `gorm:"column:password" json:"password"`
	Email     string    `gorm:"column:email" json:"email"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 表名
func (u *UserBaseModel) TableName() string {
	return "user_base"
}