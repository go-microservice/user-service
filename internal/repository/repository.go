package repository

import (
	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/go-microservice/user-service/internal/model"
)

// ProviderSet is repo providers.
var ProviderSet = wire.NewSet(NewGORMClient, NewUserBase, NewUserProfile)

func NewGORMClient() *gorm.DB {
	return model.GetDB()
}
