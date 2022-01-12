package repository

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/go-microservice/account-service/internal/model"
)

var (
	_tableUserInfoName   = (&model.UserInfoModel{}).TableName()
	_getUserInfoSQL      = "SELECT * FROM %s WHERE id = ?"
	_batchGetUserInfoSQL = "SELECT * FROM %s WHERE id IN (?)"
)

// CreateUserInfo create a item
func (r *repository) CreateUserInfo(ctx context.Context, data *model.UserInfoModel) (id int64, err error) {
	err = r.db.WithContext(ctx).Create(&data).Error
	if err != nil {
		return 0, errors.Wrap(err, "[repo] create UserInfo err")
	}

	return data.ID, nil
}

// UpdateUserInfo update item
func (r *repository) UpdateUserInfo(ctx context.Context, id int64, data *model.UserInfoModel) error {
	item, err := r.GetUserInfo(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "[repo] update UserInfo err: %v", err)
	}
	err = r.db.Model(&item).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}

// GetUserInfo get a record by primary id
func (r *repository) GetUserInfo(ctx context.Context, id int64) (ret *model.UserInfoModel, err error) {
	item := new(model.UserInfoModel)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_getUserInfoSQL, _tableUserInfoName), id).Scan(&item).Error
	if err != nil {
		return
	}
	return item, nil
}

func (r *repository) GetUserByUsername(ctx context.Context, username string) (*model.UserInfoModel, error) {
	panic("implement me")
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*model.UserInfoModel, error) {
	panic("implement me")
}

// BatchGetUserInfo batch get items by primary id
func (r *repository) BatchGetUserInfo(ctx context.Context, ids int64) (ret []*model.UserInfoModel, err error) {
	items := make([]*model.UserInfoModel, 0)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_batchGetUserInfoSQL, _tableUserInfoName), ids).Scan(&items).Error
	if err != nil {
		return
	}
	return items, nil
}