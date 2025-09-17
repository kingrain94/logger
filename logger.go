// Package logger provides a simple and efficient logging interface built on top of zap.
// It supports multiple log levels (DEBUG, INFO, WARN, ERROR, FATAL) and different environments.
package logger

import (
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Environment represents different deployment environments
type Environment int

const (
	// Development environment - enables debug logging with human-readable output
	Development Environment = iota
	// Test environment - optimized for testing with minimal output
	Test
	// Staging environment - production-like settings with more verbose logging
	Staging
	// Production environment - optimized for performance with structured logging
	Production
)

// String returns string representation of environment
func (e Environment) String() string {
	switch e {
	case Development:
		return "development"
	case Test:
		return "test"
	case Staging:
		return "staging"
	case Production:
		return "production"
	default:
		return "unknown"
	}
}

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
	mu     sync.RWMutex

	// currentEnv holds the current environment setting
	currentEnv Environment = Development
)

// Config holds logger configuration options
type Config struct {
	Environment Environment
	Level       zapcore.Level
	OutputPaths []string
	Encoding    string // "json" or "console"
}

// DefaultConfig returns a default configuration based on environment
func DefaultConfig(env Environment) Config {
	config := Config{
		Environment: env,
		OutputPaths: []string{"stdout"},
	}

	switch env {
	case Development:
		config.Level = zapcore.DebugLevel
		config.Encoding = "console"
	case Test:
		config.Level = zapcore.ErrorLevel
		config.Encoding = "json"
		config.OutputPaths = []string{}
	case Staging:
		config.Level = zapcore.InfoLevel
		config.Encoding = "json"
	case Production:
		config.Level = zapcore.WarnLevel
		config.Encoding = "json"
	}

	return config
}

func init() {
	if err := Initialize(DefaultConfig(Development)); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
}

// Initialize initializes the logger with the given configuration
func Initialize(config Config) error {
	mu.Lock()
	defer mu.Unlock()

	var zapConfig zap.Config

	switch config.Environment {
	case Development:
		zapConfig = zap.NewDevelopmentConfig()
	case Production, Staging:
		zapConfig = zap.NewProductionConfig()
	case Test:
		zapConfig = zap.NewProductionConfig()
		zapConfig.OutputPaths = []string{} // No output for tests by default
	}

	zapConfig.Level = zap.NewAtomicLevelAt(config.Level)
	zapConfig.Encoding = config.Encoding
	if len(config.OutputPaths) > 0 {
		zapConfig.OutputPaths = config.OutputPaths
	}

	var err error
	logger, err = zapConfig.Build()
	if err != nil {
		return fmt.Errorf("failed to build logger: %w", err)
	}

	sugar = logger.Sugar()
	currentEnv = config.Environment

	return nil
}

// SetEnvironment sets the environment and reinitializes the logger
func SetEnvironment(env Environment) error {
	return Initialize(DefaultConfig(env))
}

// SetLevel sets the log level dynamically
func SetLevel(level zapcore.Level) {
	mu.Lock()
	defer mu.Unlock()
	if logger != nil {
		logger = logger.WithOptions(zap.IncreaseLevel(level))
		sugar = logger.Sugar()
	}
}

// GetLogger returns the underlying zap logger for advanced usage
func GetLogger() *zap.Logger {
	mu.RLock()
	defer mu.RUnlock()
	return logger
}

// GetSugar returns the sugared logger for easier usage
func GetSugar() *zap.SugaredLogger {
	mu.RLock()
	defer mu.RUnlock()
	return sugar
}

// Sync flushes any buffered log entries
func Sync() error {
	mu.RLock()
	defer mu.RUnlock()
	if logger != nil {
		return logger.Sync()
	}
	return nil
}

// Debug logs a message at debug level with optional structured fields
func Debug(msg string, fields ...zap.Field) {
	mu.RLock()
	defer mu.RUnlock()
	if logger != nil {
		logger.Debug(msg, fields...)
	}
}

// Debugf logs a formatted message at debug level
func Debugf(template string, args ...interface{}) {
	mu.RLock()
	defer mu.RUnlock()
	if sugar != nil {
		sugar.Debugf(template, args...)
	}
}

// Info logs a message at info level with optional structured fields
func Info(msg string, fields ...zap.Field) {
	mu.RLock()
	defer mu.RUnlock()
	if logger != nil {
		logger.Info(msg, fields...)
	}
}

// Infof logs a formatted message at info level
func Infof(template string, args ...interface{}) {
	mu.RLock()
	defer mu.RUnlock()
	if sugar != nil {
		sugar.Infof(template, args...)
	}
}

// Warn logs a message at warn level with optional structured fields
func Warn(msg string, fields ...zap.Field) {
	mu.RLock()
	defer mu.RUnlock()
	if logger != nil {
		logger.Warn(msg, fields...)
	}
}

// Warnf logs a formatted message at warn level
func Warnf(template string, args ...interface{}) {
	mu.RLock()
	defer mu.RUnlock()
	if sugar != nil {
		sugar.Warnf(template, args...)
	}
}

// Error logs a message at error level with optional structured fields
func Error(msg string, fields ...zap.Field) {
	mu.RLock()
	defer mu.RUnlock()
	if logger != nil {
		logger.Error(msg, fields...)
	}
}

// Errorf logs a formatted message at error level
func Errorf(template string, args ...interface{}) {
	mu.RLock()
	defer mu.RUnlock()
	if sugar != nil {
		sugar.Errorf(template, args...)
	}
}

// Fatal logs a message at fatal level and calls os.Exit(1)
// Use with caution - this will terminate the program
func Fatal(msg string, fields ...zap.Field) {
	mu.RLock()
	defer mu.RUnlock()
	if logger != nil {
		logger.Fatal(msg, fields...)
	}
}

// Fatalf logs a formatted message at fatal level and calls os.Exit(1)
func Fatalf(template string, args ...interface{}) {
	mu.RLock()
	defer mu.RUnlock()
	if sugar != nil {
		sugar.Fatalf(template, args...)
	}
}

// With creates a child logger with additional structured context
func With(fields ...zap.Field) *zap.Logger {
	mu.RLock()
	defer mu.RUnlock()
	if logger != nil {
		return logger.With(fields...)
	}
	return nil
}

// WithFields is an alias for With for better API compatibility
func WithFields(fields ...zap.Field) *zap.Logger {
	return With(fields...)
}
