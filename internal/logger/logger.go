package logger

import (
	"log"

	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger(env string) {
	var err error

	if env == "production" {
		Log, err = zap.NewProduction()
	} else {
		Log, err = zap.NewDevelopment()
	}

	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
}

func GetLogger() *zap.Logger {
	return Log
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
