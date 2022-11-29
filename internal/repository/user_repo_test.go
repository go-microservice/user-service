package repository

import (
	"context"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/go-microservice/user-service/internal/mocks"
	"github.com/go-microservice/user-service/internal/model"
)

func Test_userRepo_CreateUser(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()
	data := &model.UserModel{
		Username:  "test-username",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `user_info` (`username`,`nickname`,`phone`,`email`,`password`,`avatar`,`gender`,`birthday`,`bio`,`login_at`,`status`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)").
		WithArgs(
			data.Username,
			data.Nickname,
			data.Phone,
			data.Email,
			data.Password,
			data.Avatar,
			data.Gender,
			data.Birthday,
			data.Bio,
			data.LoginAt,
			data.Status,
			data.CreatedAt,
			data.UpdatedAt,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	user := NewUser(mockDB, nil)
	ret, err := user.CreateUser(ctx, data)
	assert.NoError(t, err)
	assert.NotNil(t, ret)
}

func Test_userRepo_GetUser(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	var id int64
	ctx := context.Background()
	id = 1

	mockCache := mocks.NewMockUserCache(ctl)
	mockCache.EXPECT().GetUserCache(ctx, id).Return(&model.UserModel{ID: id}, nil).Times(1)

	mock.ExpectQuery("SELECT * FROM user_info WHERE id = ?").
		WithArgs(data.Username).
		WillReturnRows(
			sqlmock.NewRows([]string{`id`, `username`, `nickname`, `phone`, `email`, `password`, `avatar`, `gender`,
				`birthday`, `bio`, `login_at`, `status`, `created_at`, `updated_at`}).
				AddRow(1, data.Username, data.Nickname, data.Phone, data.Email, data.Password, data.Avatar,
					data.Gender, data.Birthday, data.Bio, data.LoginAt, data.Status, data.CreatedAt, data.UpdatedAt,
				),
		)

	user := NewUser(mockDB, mockCache)
	ret, err := user.GetUser(ctx, id)
	assert.NoError(t, err)
	assert.NotNil(t, ret)
	assert.Equal(t, id, ret.ID)
}

func Test_userRepo_GetUserByUsername(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()

	mock.ExpectQuery("SELECT * FROM user_info WHERE username = ?").
		WithArgs(data.Username).
		WillReturnRows(
			sqlmock.NewRows([]string{`id`, `username`, `nickname`, `phone`, `email`, `password`, `avatar`, `gender`,
				`birthday`, `bio`, `login_at`, `status`, `created_at`, `updated_at`}).
				AddRow(1, data.Username, data.Nickname, data.Phone, data.Email, data.Password, data.Avatar,
					data.Gender, data.Birthday, data.Bio, data.LoginAt, data.Status, data.CreatedAt, data.UpdatedAt,
				),
		)

	user := NewUser(mockDB, nil)
	ret, err := user.GetUserByUsername(ctx, data.Username)
	assert.NoError(t, err)
	assert.NotNil(t, ret)
	assert.Equal(t, data.Username, ret.Username)
}
