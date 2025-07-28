package logger

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func Init() error {
	var err error
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	
	Logger, err = config.Build()
	if err != nil {
		return err
	}
	
	return nil
}

func Close() {
	if Logger != nil {
		Logger.Sync()
	}
}