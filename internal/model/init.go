package model

import (
	"fmt"
	"sync"

	"gorm.io/gorm"

	"github.com/go-eagle/eagle/pkg/config"
	"github.com/go-eagle/eagle/pkg/storage/orm"
)

var (
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

var (
	DB   *gorm.DB
	Once sync.Once
)

// Init 初始化数据库
func Init() *gorm.DB {
	cfg, err := loadConf()
	if err != nil {
		panic(fmt.Sprintf("load orm conf err: %v", err))
	}

	DB = orm.NewMySQL(cfg)
	return DB
}

// GetDB 返回默认的数据库
func GetDB() *gorm.DB {
	Once.Do(func() {
		DB = Init()
	})
	return DB
}

// loadConf load gorm config
func loadConf() (ret *orm.Config, err error) {
	var cfg orm.Config
	if err := config.Load("database", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
