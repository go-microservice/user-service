package repository

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-microservice/user-service/internal/model"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/go-microservice/user-service/internal/cache"
)

var (
	mockDB    *gorm.DB
	mock      sqlmock.Sqlmock
	sqlDB     *sql.DB
	testCache cache.UserCache
	testRepo  UserRepo

	data *model.UserModel
)

func TestMain(m *testing.M) {
	var err error
	// mocks db
	sqlDB, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalf("mock db, err: %v", err)
	}
	_ = mock

	mockDB, err = gorm.Open(mysql.New(
		mysql.Config{
			SkipInitializeWithVersion: true,
			DriverName:                "mysql",
			Conn:                      sqlDB,
		}), &gorm.Config{})
	if err != nil {
		log.Fatalf("cannot connect to gorm, err: %v", err)
	}

	redisServer := mockRedis()
	redisClient := redis.NewClient(&redis.Options{Addr: redisServer.Addr()})

	testCache = cache.NewUserCache(redisClient)

	testRepo = NewUser(mockDB, nil)

	data = &model.UserModel{
		ID:        1,
		Username:  "test-username",
		Nickname:  "nickname",
		Phone:     "12345678",
		Email:     "test@test.com",
		Password:  "123456",
		Avatar:    "",
		Gender:    "",
		Birthday:  "2022-11-10",
		Bio:       "ok",
		LoginAt:   time.Now().Unix(),
		Status:    0,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	os.Exit(m.Run())
}

func mockRedis() *miniredis.Miniredis {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	return s
}
