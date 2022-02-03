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
	// PrefixUserProfileCacheKey cache prefix
	PrefixUserProfileCacheKey = "UserProfile:%d"
)

// UserProfileCache define a cache struct
type UserProfileCache struct {
	cache cache.Cache
}

// NewUserProfileCache new a cache
func NewUserProfileCache() *UserProfileCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""
	return &UserProfileCache{
		cache: cache.NewRedisCache(redis.RedisClient, cachePrefix, jsonEncoding, func() interface{} {
			return &model.UserProfileModel{}
		}),
	}
}

// GetUserProfileCacheKey get cache key
func (c *UserProfileCache) GetUserProfileCacheKey(id int64) string {
	return fmt.Sprintf(PrefixUserProfileCacheKey, id)
}

// SetUserProfileCache write to cache
func (c *UserProfileCache) SetUserProfileCache(ctx context.Context, id int64, data *model.UserProfileModel, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetUserProfileCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// GetUserProfileCache 获取cache
func (c *UserProfileCache) GetUserProfileCache(ctx context.Context, id int64) (ret *model.UserProfileModel, err error) {
	var data *model.UserProfileModel
	cacheKey := c.GetUserProfileCacheKey(id)
	err = c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		log.WithContext(ctx).Warnf("get err from redis, err: %+v", err)
		return nil, err
	}
	return data, nil
}

// MultiGetUserProfileCache 批量获取cache
func (c *UserProfileCache) MultiGetUserProfileCache(ctx context.Context, ids []int64) (map[string]*model.UserProfileModel, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetUserProfileCacheKey(v)
		keys = append(keys, cacheKey)
	}

	// NOTE: 需要在这里make实例化，如果在返回参数里直接定义会报 nil map
	retMap := make(map[string]*model.UserProfileModel)
	err := c.cache.MultiGet(ctx, keys, retMap)
	if err != nil {
		return nil, err
	}
	return retMap, nil
}

// MultiSetUserProfileCache 批量设置cache
func (c *UserProfileCache) MultiSetUserProfileCache(ctx context.Context, data []*model.UserProfileModel, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetUserProfileCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}
	return nil
}

// DelUserProfileCache 删除cache
func (c *UserProfileCache) DelUserProfileCache(ctx context.Context, id int64) error {
	cacheKey := c.GetUserProfileCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// DelUserProfileCache set empty cache
func (c *UserProfileCache) SetCacheWithNotFound(ctx context.Context, id int64) error {
	cacheKey := c.GetUserProfileCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
