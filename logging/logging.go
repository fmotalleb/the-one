package logging

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func BootLogger(cfg LogConfig) error {
	constructor := zap.NewProduction
	opts := make([]zap.Option, 0)
	if cfg.Development {
		constructor = zap.NewDevelopment
	}
	opts = append(opts, zap.WithCaller(cfg.ShowCaller))

	l, err := constructor(
		opts...,
	)
	if err != nil {
		return err
	}
	logger = l
	return nil
}

func BootTestLogger() {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel) // or InfoLevel, DebugLevel, etc.

	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	logger = l
}

func GetLogger(name string) *zap.Logger {
	logger = logger.Named(name)
	return logger
}

func LazyLogger(name string) func() *zap.Logger {
	return sync.OnceValue(
		func() *zap.Logger {
			return GetLogger(name)
		},
	)
}
