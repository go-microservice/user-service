package repository

import (
	"log"
	"os"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/go-microservice/user-service/internal/cache"
)

var (
	mockDB    *gorm.DB
	testCache cache.UserCache
	testRepo  UserRepo
)

func TestMain(m *testing.M) {
	// mocks db
	sqlDB, _, err := sqlmock.New()
	if err != nil {
		log.Fatalf("mock db, err: %v", err)
	}

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

	os.Exit(m.Run())
}

func mockRedis() *miniredis.Miniredis {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	return s
}
