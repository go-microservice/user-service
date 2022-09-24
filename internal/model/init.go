package model

import (
	"gorm.io/gorm"

	"github.com/go-eagle/eagle/pkg/config"
	"github.com/go-eagle/eagle/pkg/storage/orm"
)

var (
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

var (
	DB *gorm.DB
)

// Init 初始化数据库
func Init() (*gorm.DB, func(), error) {
	cfg, err := loadConf()
	if err != nil {
		return nil, nil, err
	}

	DB = orm.NewMySQL(cfg)
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, nil, err
	}
	cleanFunc := func() {
		sqlDB.Close()
	}

	return DB, cleanFunc, nil
}

// loadConf load gorm config
func loadConf() (ret *orm.Config, err error) {
	var cfg orm.Config
	if err := config.Load("database", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
