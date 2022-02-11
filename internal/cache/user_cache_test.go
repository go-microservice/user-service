package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"

	"github.com/go-microservice/user-service/internal/model"
)

var (
	redisServer *miniredis.Miniredis
	redisClient *redis.Client
	testData    = &model.UserModel{ID: 1, Username: "test"}
	uc          UserCache
)

func setup() {
	redisServer = mockRedis()
	redisClient = redis.NewClient(&redis.Options{Addr: redisServer.Addr()})
	uc = NewUserCache(redisClient)
}

func teardown() {
	redisServer.Close()
}

func mockRedis() *miniredis.Miniredis {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	return s
}

func Test_userCache_SetUserCache(t *testing.T) {
	setup()
	defer teardown()

	var id int64
	ctx := context.Background()
	id = 1
	err := uc.SetUserCache(ctx, id, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_userCache_GetUserCache(t *testing.T) {
	setup()
	defer teardown()

	var id int64
	ctx := context.Background()
	id = 1
	err := uc.SetUserCache(ctx, id, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	act, err := uc.GetUserCache(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, testData, act)
}

func Test_userCache_MultiGetUserCache(t *testing.T) {
	setup()
	defer teardown()

	ctx := context.Background()
	testData := []*model.UserModel{
		{ID: 1},
		{ID: 2},
	}
	err := uc.MultiSetUserCache(ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	expected := make(map[string]*model.UserModel)
	expected["user:1"] = &model.UserModel{ID: 1}
	expected["user:2"] = &model.UserModel{ID: 2}

	act, err := uc.MultiGetUserCache(ctx, []int64{1, 2})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expected, act)
}

func Test_userCache_MultiSetUserCache(t *testing.T) {
	setup()
	defer teardown()

	ctx := context.Background()
	testData := []*model.UserModel{
		{ID: 1},
		{ID: 2},
	}
	err := uc.MultiSetUserCache(ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_userCache_DelUserCache(t *testing.T) {
	setup()
	defer teardown()

	var id int64
	ctx := context.Background()
	id = 1
	err := uc.DelUserCache(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
}
