package repository

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/go-microservice/account-service/internal/model"
)

var (
	_tableName           = (&model.UserBaseModel{}).TableName()
	_getUserBaseSQL      = "SELECT * FROM %s WHERE id = ?"
	_batchGetUserBaseSQL = "SELECT * FROM %s WHERE id IN (?)"
)

// CreateUserBase create a item
func (r *repository) CreateUserBase(ctx context.Context, data *model.UserBaseModel) (id int64, err error) {
	err = r.db.WithContext(ctx).Create(&data).Error
	if err != nil {
		return 0, errors.Wrap(err, "[repo.user_base] create user base err")
	}

	return data.ID, nil
}

// UpdateUserBase update user base
func (r *repository) UpdateUserBase(ctx context.Context, id int64, data *model.UserBaseModel) error {
	user, err := r.GetUserBase(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "[repo.user_base] update user base err: %v", err)
	}
	err = r.db.Model(&user).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}

// GetUserBase get a user
func (r *repository) GetUserBase(ctx context.Context, uid int64) (ret *model.UserBaseModel, err error) {
	item := new(model.UserBaseModel)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_getUserBaseSQL, _tableName), uid).Scan(&item).Error
	if err != nil {
		return
	}
	return item, nil
}

// BatchGetUserBase batch get items
func (r *repository) BatchGetUserBase(ctx context.Context, ids int64) (ret []*model.UserBaseModel, err error) {
	items := make([]*model.UserBaseModel, 0)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_batchGetUserBaseSQL, _tableName), ids).Scan(&items).Error
	if err != nil {
		return
	}
	return items, nil
}
