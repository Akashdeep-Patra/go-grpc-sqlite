package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Init initializes the logger
func Init(env string) {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var err error
	log, err = config.Build()
	if err != nil {
		panic(err)
	}
}

// GetLogger returns the singleton logger instance
func GetLogger() *zap.Logger {
	if log == nil {
		// Default to development logger if not initialized
		env := os.Getenv("APP_ENV")
		if env == "" {
			env = "development"
		}
		Init(env)
	}
	return log
}

// Info logs a message at info level
func Info(msg string, fields ...zapcore.Field) {
	GetLogger().Info(msg, fields...)
}

// Error logs a message at error level
func Error(msg string, fields ...zapcore.Field) {
	GetLogger().Error(msg, fields...)
}

// Debug logs a message at debug level
func Debug(msg string, fields ...zapcore.Field) {
	GetLogger().Debug(msg, fields...)
}

// Warn logs a message at warn level
func Warn(msg string, fields ...zapcore.Field) {
	GetLogger().Warn(msg, fields...)
}

// Fatal logs a message at fatal level and then calls os.Exit(1)
func Fatal(msg string, fields ...zapcore.Field) {
	GetLogger().Fatal(msg, fields...)
}

// With creates a child logger with the given fields
func With(fields ...zapcore.Field) *zap.Logger {
	return GetLogger().With(fields...)
}

// Sync flushes any buffered log entries
func Sync() error {
	return GetLogger().Sync()
} 