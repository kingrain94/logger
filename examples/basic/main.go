package main

import (
	"github.com/kingrain94/logger"
	"go.uber.org/zap"
)

func main() {
	// Basic logging examples
	logger.Info("Application started", zap.String("version", "1.0.0"))
	logger.Debug("Debug information", zap.Int("user_id", 12345))
	logger.Warn("This is a warning", zap.String("component", "auth"))
	logger.Error("An error occurred", zap.String("error", "connection failed"))

	// Formatted logging
	logger.Infof("User %s logged in from %s", "john_doe", "192.168.1.1")
	logger.Warnf("Memory usage is at %d%%", 85)
	logger.Errorf("Failed to process request: %v", "timeout")

	// Structured logging with multiple fields
	logger.Info("User action",
		zap.String("user_id", "12345"),
		zap.String("action", "login"),
		zap.String("ip", "192.168.1.1"),
		zap.Int("attempt", 1),
		zap.Bool("success", true),
	)

	// Don't forget to sync logs before exit
	defer logger.Sync()
}
