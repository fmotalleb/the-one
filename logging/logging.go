package logging

import (
	"go.uber.org/zap"
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

func GetLogger(name string) *zap.Logger {
	logger = logger.Named(name)
	return logger
}
