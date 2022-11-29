package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-microservice/user-service/internal/model"

	"github.com/go-microservice/user-service/internal/mocks"
	"github.com/golang/mock/gomock"
)

func Test_userRepo_CreateUser(t *testing.T) {

}

func Test_userRepo_GetUser(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	var id int64
	ctx := context.Background()
	id = 1

	mockCache := mocks.NewMockUserCache(ctl)
	mockCache.EXPECT().GetUserCache(ctx, id).Return(&model.UserModel{ID: id}, nil).Times(1)

	user := NewUser(mockDB, mockCache)
	ret, err := user.GetUser(ctx, id)
	assert.NoError(t, err)
	assert.NotNil(t, ret)
	assert.Equal(t, id, ret.ID)
}

func Test_userRepo_GetUserByUsername(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	var (
		id       int64
		username string
	)
	ctx := context.Background()
	id = 1
	username = "test-username"

	// todo: mock expectQuery

	user := NewUser(mockDB, nil)
	ret, err := user.GetUserByUsername(ctx, username)
	assert.NoError(t, err)
	assert.NotNil(t, ret)
	assert.Equal(t, id, ret.ID)
	assert.Equal(t, username, ret.Username)
}
