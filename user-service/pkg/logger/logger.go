package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func Init(logLevel string) error {
	var cfg zap.Config

	// Set the log level
	level := zapcore.InfoLevel
	err := level.UnmarshalText([]byte(logLevel))
	if err != nil {
		level = zapcore.InfoLevel
	}

	// Configure the logger
	cfg = zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      false,
		Encoding:         "json", // Use "json" or "console"
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder, // e.g. "info"
			EncodeTime:     zapcore.ISO8601TimeEncoder,    // e.g. "2021-01-01T12:34:56.789Z"
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder, // e.g. "pkg/file.go:123"
		},
	}

	// Build the logger
	Logger, err = cfg.Build()
	if err != nil {
		return err
	}

	return nil
}

// Sync flushes any buffered log entries
func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}
