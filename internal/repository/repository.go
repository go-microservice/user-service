package repository

import (
	"github.com/go-microservice/user-service/internal/model"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// ProviderSet is repo providers.
var ProviderSet = wire.NewSet(NewGORMClient, NewUserBase, NewUserProfile)

func NewGORMClient() *gorm.DB {
	return model.GetDB()
}
