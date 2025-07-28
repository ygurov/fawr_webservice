package storage

import (
	"github.com/fawrwebservice/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDB(path string) *gorm.DB {
	handle, err := gorm.Open(sqlite.Open(path))
	if err != nil {
		panic(err)
	}

	handle.AutoMigrate(&model.Comment{})

	return handle
}
