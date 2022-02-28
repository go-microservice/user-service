package repository

import (
	"github.com/go-microservice/user-service/internal/model"
	"github.com/google/wire"
)

// ProviderSet is repo providers.
var ProviderSet = wire.NewSet(model.GetDB, NewUser)
