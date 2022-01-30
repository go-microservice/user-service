package repository

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"

	"github.com/pkg/errors"

	"github.com/go-microservice/user-service/internal/model"
)

var (
	_tableUserProfileName   = (&model.UserProfileModel{}).TableName()
	_getUserProfileSQL      = "SELECT * FROM %s WHERE id = ?"
	_batchGetUserProfileSQL = "SELECT * FROM %s WHERE id IN (?)"
)

var _ UserProfileRepo = (*userProfileRepo)(nil)

// UserProfileRepo define a repo interface
type UserProfileRepo interface {
	CreateUserProfile(ctx context.Context, data *model.UserProfileModel) (id int64, err error)
	UpdateUserProfile(ctx context.Context, id int64, data *model.UserProfileModel) error
	GetUserProfile(ctx context.Context, id int64) (ret *model.UserProfileModel, err error)
	BatchGetUserProfile(ctx context.Context, ids int64) (ret []*model.UserProfileModel, err error)
}

// userProfileRepo struct
type userProfileRepo struct {
	db     *gorm.DB
	tracer trace.Tracer
}

// New new a repository and return
func NewUserProfile(db *gorm.DB) UserProfileRepo {
	return &userProfileRepo{
		db:     db,
		tracer: otel.Tracer("repository"),
	}
}

// CreateUserProfile create a item
func (r *userProfileRepo) CreateUserProfile(ctx context.Context, data *model.UserProfileModel) (id int64, err error) {
	err = r.db.WithContext(ctx).Create(&data).Error
	if err != nil {
		return 0, errors.Wrap(err, "[repo] create UserProfile err")
	}

	return data.ID, nil
}

// UpdateUserProfile update item
func (r *userProfileRepo) UpdateUserProfile(ctx context.Context, id int64, data *model.UserProfileModel) error {
	item, err := r.GetUserProfile(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "[repo] update UserProfile err: %v", err)
	}
	err = r.db.Model(&item).Updates(data).Error
	if err != nil {
		return err
	}

	return nil
}

// GetUserProfile get a record
func (r *userProfileRepo) GetUserProfile(ctx context.Context, id int64) (ret *model.UserProfileModel, err error) {
	item := new(model.UserProfileModel)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_getUserProfileSQL, _tableUserProfileName), id).Scan(&item).Error
	if err != nil {
		return
	}

	return item, nil
}

// BatchGetUserProfile batch get items
func (r *userProfileRepo) BatchGetUserProfile(ctx context.Context, ids int64) (ret []*model.UserProfileModel, err error) {
	items := make([]*model.UserProfileModel, 0)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_batchGetUserProfileSQL, _tableUserProfileName), ids).Scan(&items).Error
	if err != nil {
		return
	}

	return items, nil
}
