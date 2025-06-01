package logging

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapWriter struct {
	output *zap.Logger
	level  zapcore.Level
}

func NewZapWriter(name string, level ...zapcore.Level) io.Writer {
	log := GetLogger("service." + name)
	lvl := zapcore.InfoLevel
	if len(level) != 0 {
		lvl = level[0]
	}
	return &ZapWriter{
		output: log,
		level:  lvl,
	}
}

func (b *ZapWriter) Write(p []byte) (n int, err error) {
	b.output.Log(b.level, string(p))
	return len(p), nil
}
