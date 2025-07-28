package storage

import (
	"github.com/fawrwebservice/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDB() *gorm.DB {
	handle, err := gorm.Open(sqlite.Open("test.db"))
	if err != nil {
		panic(err)
	}

	handle.AutoMigrate(&model.Comment{})

	return handle
}
