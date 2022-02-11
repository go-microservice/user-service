package repository

import (
	"context"
	"testing"

	"github.com/go-microservice/user-service/internal/model"

	mock_cache "github.com/go-microservice/user-service/internal/mock/cache"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	mockDB *gorm.DB
)

func setup() {
	// mock db
	sqlDB, _, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	mockDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	_ = mockDB
}

func teardown() {
}

func Test_userRepo_CreateUser(t *testing.T) {

}

// see: https://segmentfault.com/a/1190000017132133
func Test_userRepo_GetUser(t *testing.T) {
	setup()
	defer teardown()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	var id int64
	ctx := context.Background()
	id = 1

	mockCache := mock_cache.NewMockUserCache(ctl)
	gomock.InOrder(
		mockCache.EXPECT().GetUserCache(ctx, id).Return(&model.UserModel{ID: 1}, nil).Times(1),
	)

	user := NewUser(mockDB, mockCache)
	ret, err := user.GetUser(ctx, id)
	if err != nil {
		t.Errorf("repo.GetUser err: %v", err)
	}
	t.Log(ret)
}
