package repository

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"

	"github.com/go-microservice/user-service/internal/model"
)

var (
	_tableUserBaseName        = (&model.UserBaseModel{}).TableName()
	_getUserBaseSQL           = "SELECT * FROM %s WHERE id = ?"
	_getUserBaseByUsernameSQL = "SELECT * FROM %s WHERE username = ?"
	_getUserBaseByEmailSQL    = "SELECT * FROM %s WHERE email = ?"
	_getUserBaseByPhoneSQL    = "SELECT * FROM %s WHERE phone = ?"
	_batchGetUserBaseSQL      = "SELECT * FROM %s WHERE id IN (%s)"
)

var _ UserBaseRepo = (*userBaseRepo)(nil)

// UserBaseRepo define a repo interface
type UserBaseRepo interface {
	CreateUserBase(ctx context.Context, data *model.UserBaseModel) (id int64, err error)
	UpdateUserBase(ctx context.Context, id int64, data *model.UserBaseModel) error
	GetUserBase(ctx context.Context, id int64) (ret *model.UserBaseModel, err error)
	GetUserByUsername(ctx context.Context, username string) (ret *model.UserBaseModel, err error)
	GetUserByEmail(ctx context.Context, email string) (ret *model.UserBaseModel, err error)
	GetUserByPhone(ctx context.Context, phone string) (ret *model.UserBaseModel, err error)
	BatchGetUserBase(ctx context.Context, ids string) (ret []*model.UserBaseModel, err error)
}

type userBaseRepo struct {
	db     *gorm.DB
	tracer trace.Tracer
}

// NewUserBase new a repository and return
func NewUserBase(db *gorm.DB) UserBaseRepo {
	return &userBaseRepo{
		db:     db,
		tracer: otel.Tracer("userBaseRepo"),
	}
}

// CreateUserBase create a item
func (r *userBaseRepo) CreateUserBase(ctx context.Context, data *model.UserBaseModel) (id int64, err error) {
	err = r.db.WithContext(ctx).Create(&data).Error
	if err != nil {
		return 0, errors.Wrap(err, "[repo] create UserBase err")
	}

	return data.ID, nil
}

// UpdateUserBase update item
func (r *userBaseRepo) UpdateUserBase(ctx context.Context, id int64, data *model.UserBaseModel) error {
	item, err := r.GetUserBase(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "[repo] update UserBase err: %v", err)
	}
	err = r.db.Model(&item).Updates(data).Error
	if err != nil {
		return err
	}

	// delete cache

	return nil
}

// GetUserBase get a record by primary id
func (r *userBaseRepo) GetUserBase(ctx context.Context, id int64) (ret *model.UserBaseModel, err error) {
	item := new(model.UserBaseModel)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_getUserBaseSQL, _tableUserBaseName), id).Scan(&item).Error
	if err != nil {
		return
	}

	return item, nil
}

func (r *userBaseRepo) GetUserByUsername(ctx context.Context, username string) (ret *model.UserBaseModel, err error) {
	item := new(model.UserBaseModel)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_getUserBaseByUsernameSQL, _tableUserBaseName), username).Scan(&item).Error
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *userBaseRepo) GetUserByEmail(ctx context.Context, email string) (ret *model.UserBaseModel, err error) {
	item := new(model.UserBaseModel)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_getUserBaseByEmailSQL, _tableUserBaseName), email).Scan(&item).Error
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *userBaseRepo) GetUserByPhone(ctx context.Context, phone string) (ret *model.UserBaseModel, err error) {
	item := new(model.UserBaseModel)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_getUserBaseByPhoneSQL, _tableUserBaseName), phone).Scan(&item).Error
	if err != nil {
		return
	}

	return item, nil
}

// BatchGetUserBase batch get items by primary id
func (r *userBaseRepo) BatchGetUserBase(ctx context.Context, ids string) (ret []*model.UserBaseModel, err error) {
	items := make([]*model.UserBaseModel, 0)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_batchGetUserBaseSQL, _tableUserBaseName, ids)).Scan(&items).Error
	if err != nil {
		return
	}

	return items, nil
}
