package logging

import (
	"sync"

	"go.uber.org/zap"
)

var logger *zap.Logger

func SetRootLogger(l *zap.Logger) {
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
