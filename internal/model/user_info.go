package model

// UserInfoModel define a user base info struct.
type UserInfoModel struct {
	ID        int64  `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	Username  string `gorm:"column:username" json:"username"`
	Phone     string `gorm:"column:phone" json:"phone"`
	Email     string `gorm:"column:email" json:"email"`
	Password  string `gorm:"column:password" json:"password"`
	LoginAt   int64  `gorm:"column:login_at" json:"login_at"` // login time for last times
	Status    int32  `gorm:"column:status" json:"status"`
	CreatedAt int64  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt int64  `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 表名
func (u *UserInfoModel) TableName() string {
	return "user_info"
}
