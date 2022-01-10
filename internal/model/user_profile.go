package model

// UserProfileModel define a user profile struct.
type UserProfileModel struct {
	ID        int64  `gorm:"primary_key;column:id" json:"id"`
	UserID    int64  `gorm:"primary_key;column:user_id" json:"user_id"`
	Nickname  string `gorm:"column:username" json:"nickname"`
	Avatar    string `gorm:"column:avatar" json:"avatar"`
	Gender    string `gorm:"column:gender" json:"gender"`
	Birthday  string `gorm:"column:birthday" json:"birthday"`
	Bio       string `gorm:"column:bio" json:"bio"`
	CreatedAt int64  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt int64  `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 表名
func (u *UserProfileModel) TableName() string {
	return "user_profile"
}
