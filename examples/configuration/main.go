package main

import (
	"os"

	"github.com/kingrain94/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Example 1: Environment-based configuration
	logger.Info("Using default development environment")
	logger.Debug("This debug message will be shown in development")

	// Example 2: Switching environments
	logger.SetEnvironment(logger.Production)
	logger.Info("Switched to production environment")
	logger.Debug("This debug message will NOT be shown in production")

	// Example 3: Custom configuration
	customConfig := logger.Config{
		Environment: logger.Development,
		Level:       zapcore.InfoLevel, // Only info and above
		OutputPaths: []string{"stdout", "app.log"},
		Encoding:    "json",
	}

	if err := logger.Initialize(customConfig); err != nil {
		logger.Fatal("Failed to initialize custom logger", zap.Error(err))
	}

	logger.Info("Using custom configuration with JSON encoding")
	logger.Debug("This debug message will NOT be shown (level is INFO)")

	// Example 4: Dynamic level changes
	logger.Info("Current level allows INFO messages")
	logger.SetLevel(zapcore.DebugLevel)
	logger.Debug("Now debug messages are enabled!")
	logger.Info("INFO messages still work")

	// Example 5: File logging
	fileConfig := logger.Config{
		Environment: logger.Production,
		Level:       zapcore.InfoLevel,
		OutputPaths: []string{"application.log"},
		Encoding:    "json",
	}

	if err := logger.Initialize(fileConfig); err != nil {
		logger.Fatal("Failed to initialize file logger", zap.Error(err))
	}

	logger.Info("This message goes to application.log file",
		zap.String("config", "file-based"),
		zap.Bool("structured", true),
	)

	// Example 6: Multiple outputs
	multiConfig := logger.Config{
		Environment: logger.Development,
		Level:       zapcore.DebugLevel,
		OutputPaths: []string{"stdout", "debug.log", "app.log"},
		Encoding:    "console",
	}

	if err := logger.Initialize(multiConfig); err != nil {
		logger.Fatal("Failed to initialize multi-output logger", zap.Error(err))
	}

	logger.Info("This message goes to console and multiple log files",
		zap.String("feature", "multi-output"),
		zap.Int("outputs", 3),
	)

	// Example 7: Environment variable based configuration
	env := os.Getenv("APP_ENV")
	switch env {
	case "production":
		logger.SetEnvironment(logger.Production)
	case "staging":
		logger.SetEnvironment(logger.Staging)
	case "test":
		logger.SetEnvironment(logger.Test)
	default:
		logger.SetEnvironment(logger.Development)
	}

	logger.Info("Environment configured from APP_ENV variable",
		zap.String("env", env),
	)

	// Clean up log files created during examples
	defer func() {
		logger.Sync()
		os.Remove("app.log")
		os.Remove("application.log")
		os.Remove("debug.log")
	}()
}
