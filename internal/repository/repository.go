package repository

import (
	"context"

	"github.com/go-microservice/account-service/internal/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

var (
	// ErrNotFound data is not exist
	ErrNotFound = gorm.ErrRecordNotFound
)

var _ Repository = (*repository)(nil)

// Repository define a repo interface
type Repository interface {
	// user info
	CreateUserInfo(ctx context.Context, data *model.UserInfoModel) (id int64, err error)
	UpdateUserInfo(ctx context.Context, id int64, data *model.UserInfoModel) error
	GetUserInfo(ctx context.Context, id int64) (ret *model.UserInfoModel, err error)
	GetUserByUsername(ctx context.Context, username string) (ret *model.UserInfoModel, err error)
	GetUserByEmail(ctx context.Context, email string) (ret *model.UserInfoModel, err error)
	GetUserByPhone(ctx context.Context, phone string) (ret *model.UserInfoModel, err error)
	BatchGetUserInfo(ctx context.Context, ids int64) (ret []*model.UserInfoModel, err error)

	// user profile
	CreateUserProfile(ctx context.Context, data *model.UserProfileModel) (id int64, err error)
	UpdateUserProfile(ctx context.Context, id int64, data *model.UserProfileModel) error
	GetUserProfile(ctx context.Context, id int64) (ret *model.UserProfileModel, err error)
	BatchGetUserProfile(ctx context.Context, ids int64) (ret []*model.UserProfileModel, err error)
}

// repository mysql struct
type repository struct {
	db     *gorm.DB
	tracer trace.Tracer
}

// New new a repository and return
func New(db *gorm.DB) Repository {
	return &repository{
		db:     db,
		tracer: otel.Tracer("repository"),
	}
}

// Close release mysql connection
func (r *repository) Close() {

}
