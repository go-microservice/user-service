package model

import (
	"gorm.io/gorm"

	"github.com/go-eagle/eagle/pkg/storage/orm"
)

var (
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

var (
	// DB define a gloabl db
	DB *gorm.DB
)

// Init init db
func Init() (*gorm.DB, func(), error) {
	err := orm.New([]string{"default"}...)
	if err != nil {
		return nil, nil, err
	}

	// get first db
	DB, err := orm.GetDB("default")
	if err != nil {
		return nil, nil, err
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, nil, err
	}

	// here you can add second or more db, and remember to add close to below cleanFunc
	// ...

	cleanFunc := func() {
		sqlDB.Close()
	}
	return DB, cleanFunc, nil
}
