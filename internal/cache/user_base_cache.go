package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-eagle/eagle/pkg/cache"
	"github.com/go-eagle/eagle/pkg/encoding"
	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/redis"

	"github.com/go-microservice/user-service/internal/model"
)

const (
	// PrefixUserBaseCacheKey cache prefix
	PrefixUserBaseCacheKey = "UserBase:%d"
)

// UserBaseCache define a cache struct
type UserBaseCache struct {
	cache cache.Cache
}

// NewUserBaseCache new a cache
func NewUserBaseCache() *UserBaseCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""
	return &UserBaseCache{
		cache: cache.NewRedisCache(redis.RedisClient, cachePrefix, jsonEncoding, func() interface{} {
			return &model.UserBaseModel{}
		}),
	}
}

// GetUserBaseCacheKey get cache key
func (c *UserBaseCache) GetUserBaseCacheKey(id int64) string {
	return fmt.Sprintf(PrefixUserBaseCacheKey, id)
}

// SetUserBaseCache write to cache
func (c *UserBaseCache) SetUserBaseCache(ctx context.Context, id int64, data *model.UserBaseModel, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetUserBaseCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// GetUserBaseCache 获取cache
func (c *UserBaseCache) GetUserBaseCache(ctx context.Context, id int64) (data *model.UserBaseModel, err error) {
	cacheKey := c.GetUserBaseCacheKey(id)
	err = c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		log.WithContext(ctx).Warnf("get err from redis, err: %+v", err)
		return nil, err
	}
	return data, nil
}

// MultiGetUserBaseCache 批量获取cache
func (c *UserBaseCache) MultiGetUserBaseCache(ctx context.Context, ids []int64) (map[string]*model.UserBaseModel, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetUserBaseCacheKey(v)
		keys = append(keys, cacheKey)
	}

	// NOTE: 需要在这里make实例化，如果在返回参数里直接定义会报 nil map
	retMap := make(map[string]*model.UserBaseModel)
	err := c.cache.MultiGet(ctx, keys, retMap)
	if err != nil {
		return nil, err
	}
	return retMap, nil
}

// DelUserBaseCache 删除cache
func (c *UserBaseCache) DelUserBaseCache(ctx context.Context, id int64) error {
	cacheKey := c.GetUserBaseCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// DelUserBaseCache set empty cache
func (c *UserBaseCache) SetCacheWithNotFound(ctx context.Context, id int64) error {
	cacheKey := c.GetUserBaseCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
