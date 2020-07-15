package logger

import (
	"go.uber.org/zap"
)

type environment int

const (
	// DEVELOPMENT environment
	DEVELOPMENT environment = iota
	// TEST environment
	TEST
	// STAGING environment
	STAGING
	// PRODUCTION environment
	PRODUCTION
)

var (
	logger *zap.Logger

	// default is development
	env environment = DEVELOPMENT
)

func init() {
	SetDebugMode(env == DEVELOPMENT)
}

// SetEnvironment set environment that use logger
func SetEnvironment(e environment) {
	env = e
}

// SetDebugMode enable debug mode for logging
func SetDebugMode(isDebug bool) {
	if isDebug {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
}

// Debug Wrapper Function
// If there is information you want to output in the application code, use Debug in principle.
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Info Wrapper Function
// Use Info for all basic information that should be output, such as environment information.
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Warn Wrapper Function
// It is not necessary to stop the relevant process, but it is used when there is an event that requires some investigation or countermeasure in the related function.
// Also, in the case that there is a process that needs to be improved and it cannot be improved immediately,
// Output information at Warn level (abolition of calling previous version due to api version upgrade, etc.).
// In principle, warn is associated with alert notifications and other operations.
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Fatal Wrapper Function
// Do not use for online processing. Available only when an error occurs in batch processing. Be careful to determine the cause of the message.
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}
