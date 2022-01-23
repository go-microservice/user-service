package repository

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

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

// CreateUserBase create a item
func (r *repository) CreateUserBase(ctx context.Context, data *model.UserBaseModel) (id int64, err error) {
	err = r.db.WithContext(ctx).Create(&data).Error
	if err != nil {
		return 0, errors.Wrap(err, "[repo] create UserBase err")
	}

	return data.ID, nil
}

// UpdateUserBase update item
func (r *repository) UpdateUserBase(ctx context.Context, id int64, data *model.UserBaseModel) error {
	item, err := r.GetUserBase(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "[repo] update UserBase err: %v", err)
	}
	err = r.db.Model(&item).Updates(data).Error
	if err != nil {
		return err
	}

	return nil
}

// GetUserBase get a record by primary id
func (r *repository) GetUserBase(ctx context.Context, id int64) (ret *model.UserBaseModel, err error) {
	item := new(model.UserBaseModel)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_getUserBaseSQL, _tableUserBaseName), id).Scan(&item).Error
	if err != nil {
		return
	}

	return item, nil
}

func (r *repository) GetUserByUsername(ctx context.Context, username string) (ret *model.UserBaseModel, err error) {
	item := new(model.UserBaseModel)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_getUserBaseByUsernameSQL, _tableUserBaseName), username).Scan(&item).Error
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (ret *model.UserBaseModel, err error) {
	item := new(model.UserBaseModel)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_getUserBaseByEmailSQL, _tableUserBaseName), email).Scan(&item).Error
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *repository) GetUserByPhone(ctx context.Context, phone string) (ret *model.UserBaseModel, err error) {
	item := new(model.UserBaseModel)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_getUserBaseByPhoneSQL, _tableUserBaseName), phone).Scan(&item).Error
	if err != nil {
		return
	}

	return item, nil
}

// BatchGetUserBase batch get items by primary id
func (r *repository) BatchGetUserBase(ctx context.Context, ids string) (ret []*model.UserBaseModel, err error) {
	items := make([]*model.UserBaseModel, 0)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_batchGetUserBaseSQL, _tableUserBaseName, ids)).Scan(&items).Error
	if err != nil {
		return
	}

	return items, nil
}
