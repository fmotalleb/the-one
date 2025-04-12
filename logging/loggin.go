package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func GetLogger(name string) *zap.Logger {
	logger, err := zap.NewProduction(
		zap.WithCaller(true),
		zap.WithClock(zapcore.DefaultClock),
	)
	if err != nil {
		panic(err)
	}

	logger = logger.Named(name)

	return logger
}
