package cache

//go:generate mockgen -source=user_cache.go -destination=../../internal/mock/user_cache_mock.go  -package mock

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cast"

	"github.com/go-eagle/eagle/pkg/cache"
	"github.com/go-eagle/eagle/pkg/encoding"
	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-redis/redis/v8"

	"github.com/go-microservice/user-service/internal/model"
)

const (
	// PrefixUserCacheKey cache prefix
	PrefixUserCacheKey = "user:%d"
)

type UserCache interface {
	SetUserCache(ctx context.Context, id int64, data *model.UserModel, duration time.Duration) error
	GetUserCache(ctx context.Context, id int64) (ret *model.UserModel, err error)
	MultiGetUserCache(ctx context.Context, ids []int64) (map[string]*model.UserModel, error)
	MultiSetUserCache(ctx context.Context, data []*model.UserModel, duration time.Duration) error
	DelUserCache(ctx context.Context, id int64) error
	SetCacheWithNotFound(ctx context.Context, id int64) error
}

// userCache define a cache struct
type userCache struct {
	cache cache.Cache
}

// NewUserCache new a cache
func NewUserCache(rdb *redis.Client) UserCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""
	return &userCache{
		cache: cache.NewRedisCache(rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.UserModel{}
		}),
	}
}

// GetUserCacheKey get cache key
func (c *userCache) GetUserCacheKey(id int64) string {
	return fmt.Sprintf(PrefixUserCacheKey, id)
}

// SetUserCache write to cache
func (c *userCache) SetUserCache(ctx context.Context, id int64, data *model.UserModel, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetUserCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// GetUserCache 获取cache
func (c *userCache) GetUserCache(ctx context.Context, id int64) (ret *model.UserModel, err error) {
	var data *model.UserModel
	cacheKey := c.GetUserCacheKey(id)
	err = c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		log.WithContext(ctx).Warnf("get err from redis, err: %+v", err)
		return nil, err
	}
	return data, nil
}

// MultiGetUserCache 批量获取cache
func (c *userCache) MultiGetUserCache(ctx context.Context, ids []int64) (map[string]*model.UserModel, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetUserCacheKey(v)
		keys = append(keys, cacheKey)
	}

	// NOTE: 需要在这里make实例化，如果在返回参数里直接定义会报 nil map
	itemMap := make(map[string]*model.UserModel)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[string]*model.UserModel)
	for _, v := range ids {
		val, ok := itemMap[c.GetUserCacheKey(v)]
		if ok {
			retMap[cast.ToString(v)] = val
		}
	}
	return retMap, nil
}

// MultiSetUserCache 批量设置cache
func (c *userCache) MultiSetUserCache(ctx context.Context, data []*model.UserModel, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetUserCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}
	return nil
}

// DelUserCache 删除cache
func (c *userCache) DelUserCache(ctx context.Context, id int64) error {
	cacheKey := c.GetUserCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// DelUserCache set empty cache
func (c *userCache) SetCacheWithNotFound(ctx context.Context, id int64) error {
	cacheKey := c.GetUserCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
