package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-microservice/user-service/internal/cache"

	"github.com/spf13/cast"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"

	"github.com/pkg/errors"

	"github.com/go-microservice/user-service/internal/model"
)

var (
	_tableUserProfileName   = (&model.UserProfileModel{}).TableName()
	_getUserProfileSQL      = "SELECT * FROM %s WHERE id = ?"
	_batchGetUserProfileSQL = "SELECT * FROM %s WHERE id IN (%s)"
)

var _ UserProfileRepo = (*userProfileRepo)(nil)

// UserProfileRepo define a repo interface
type UserProfileRepo interface {
	CreateUserProfile(ctx context.Context, data *model.UserProfileModel) (id int64, err error)
	UpdateUserProfile(ctx context.Context, id int64, data *model.UserProfileModel) error
	GetUserProfile(ctx context.Context, id int64) (ret *model.UserProfileModel, err error)
	BatchGetUserProfile(ctx context.Context, ids []int64) (ret []*model.UserProfileModel, err error)
}

// userProfileRepo struct
type userProfileRepo struct {
	db     *gorm.DB
	tracer trace.Tracer
	cache  *cache.UserProfileCache
}

// New new a repository and return
func NewUserProfile(db *gorm.DB, cache *cache.UserProfileCache) UserProfileRepo {
	return &userProfileRepo{
		db:     db,
		tracer: otel.Tracer("repository"),
		cache:  cache,
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
func (r *userProfileRepo) BatchGetUserProfile(ctx context.Context, ids []int64) (ret []*model.UserProfileModel, err error) {
	items := make([]*model.UserProfileModel, 0)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_batchGetUserProfileSQL, _tableUserProfileName), ids).Scan(&items).Error
	if err != nil {
		return
	}

	idsStr := cast.ToStringSlice(ids)
	itemMap, err := r.cache.MultiGetUserProfileCache(ctx, ids)
	if err != nil {
		return nil, err
	}
	var missedID []int64
	for _, v := range ids {
		item, ok := itemMap[cast.ToString(v)]
		if !ok {
			missedID = append(missedID, v)
			continue
		}
		ret = append(ret, item)
	}
	// get missed data
	if len(missedID) > 0 {
		var missedData []*model.UserProfileModel
		_sql := fmt.Sprintf(_batchGetUserProfileSQL, _tableUserProfileName, strings.Join(idsStr, ","))
		err = r.db.WithContext(ctx).Raw(_sql).Scan(&missedData).Error
		if err != nil {
			// you can degrade to ignore error
			return nil, err
		}
		if len(missedData) > 0 {
			ret = append(ret, missedData...)
			err = r.cache.MultiSetUserProfileCache(ctx, ret, 5*time.Minute)
			if err != nil {
				// you can degrade to ignore error
				return nil, err
			}
		}
	}
	return ret, nil
}
