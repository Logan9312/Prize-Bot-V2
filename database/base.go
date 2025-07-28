package database

import (
	"github.com/Logan9312/Prize-Bot-V2/logger"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Init() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		logger.Logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	db.AutoMigrate(&User{})

}
