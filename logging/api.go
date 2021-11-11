package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

var DefaultLogger *Logger

func init() {
	DefaultLogger = new(Logger)
	cfg := zap.NewProductionConfig()
	cfg.Encoding = "console"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.DisableStacktrace = true
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	DefaultLogger.Logger = logger
}
