package repository

import (
	"context"

	"github.com/go-microservice/user-service/internal/model"
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
	CreateUserBase(ctx context.Context, data *model.UserBaseModel) (id int64, err error)
	UpdateUserBase(ctx context.Context, id int64, data *model.UserBaseModel) error
	GetUserBase(ctx context.Context, id int64) (ret *model.UserBaseModel, err error)
	GetUserByUsername(ctx context.Context, username string) (ret *model.UserBaseModel, err error)
	GetUserByEmail(ctx context.Context, email string) (ret *model.UserBaseModel, err error)
	GetUserByPhone(ctx context.Context, phone string) (ret *model.UserBaseModel, err error)
	BatchGetUserBase(ctx context.Context, ids string) (ret []*model.UserBaseModel, err error)

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
