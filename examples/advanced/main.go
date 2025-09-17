package main

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/kingrain94/logger"
	"go.uber.org/zap"
)

// User represents a user in our system
type User struct {
	ID       string
	Username string
	Email    string
	IsActive bool
}

// UserService demonstrates advanced logging patterns
type UserService struct {
	logger *zap.Logger
}

// NewUserService creates a new user service with a dedicated logger
func NewUserService() *UserService {
	serviceLogger := logger.With(
		zap.String("component", "user-service"),
		zap.String("version", "1.2.0"),
	)
	return &UserService{logger: serviceLogger}
}

// GetUser demonstrates structured error logging
func (s *UserService) GetUser(ctx context.Context, userID string) (*User, error) {
	// Log method entry with context
	s.logger.Debug("GetUser called",
		zap.String("user_id", userID),
		zap.String("trace_id", getTraceID(ctx)),
	)

	// Simulate database lookup
	if userID == "invalid" {
		err := errors.New("user not found")
		s.logger.Error("Failed to get user",
			zap.String("user_id", userID),
			zap.Error(err),
			zap.String("trace_id", getTraceID(ctx)),
		)
		return nil, err
	}

	// Simulate slow query warning
	if userID == "slow" {
		time.Sleep(100 * time.Millisecond)
		s.logger.Warn("Slow database query detected",
			zap.String("user_id", userID),
			zap.Duration("duration", 100*time.Millisecond),
			zap.String("query", "SELECT * FROM users WHERE id = ?"),
		)
	}

	user := &User{
		ID:       userID,
		Username: "user_" + userID,
		Email:    "user" + userID + "@example.com",
		IsActive: true,
	}

	// Log successful operation
	s.logger.Info("User retrieved successfully",
		zap.String("user_id", userID),
		zap.String("username", user.Username),
		zap.Bool("is_active", user.IsActive),
		zap.String("trace_id", getTraceID(ctx)),
	)

	return user, nil
}

// ProcessUsers demonstrates concurrent logging
func (s *UserService) ProcessUsers(ctx context.Context, userIDs []string) {
	s.logger.Info("Starting batch user processing",
		zap.Int("user_count", len(userIDs)),
		zap.String("trace_id", getTraceID(ctx)),
	)

	var wg sync.WaitGroup
	for i, userID := range userIDs {
		wg.Add(1)
		go func(id string, index int) {
			defer wg.Done()

			// Create a child logger with additional context
			workerLogger := s.logger.With(
				zap.String("worker_id", id),
				zap.Int("worker_index", index),
			)

			workerLogger.Debug("Processing user", zap.String("user_id", id))

			user, err := s.GetUser(ctx, id)
			if err != nil {
				workerLogger.Error("Failed to process user",
					zap.String("user_id", id),
					zap.Error(err),
				)
				return
			}

			// Simulate processing
			time.Sleep(50 * time.Millisecond)

			workerLogger.Info("User processed successfully",
				zap.String("user_id", user.ID),
				zap.String("username", user.Username),
			)
		}(userID, i)
	}

	wg.Wait()
	s.logger.Info("Batch user processing completed",
		zap.Int("user_count", len(userIDs)),
		zap.String("trace_id", getTraceID(ctx)),
	)
}

// getTraceID simulates getting a trace ID from context
func getTraceID(ctx context.Context) string {
	if traceID := ctx.Value("trace_id"); traceID != nil {
		return traceID.(string)
	}
	return "unknown"
}

func main() {
	// Initialize with development settings for detailed output
	logger.SetEnvironment(logger.Development)

	// Example 1: Basic service usage
	userService := NewUserService()

	ctx := context.WithValue(context.Background(), "trace_id", "abc123")

	// Example 2: Successful operation
	user, err := userService.GetUser(ctx, "12345")
	if err != nil {
		logger.Error("Failed to get user", zap.Error(err))
	} else {
		logger.Info("Got user", zap.Any("user", user))
	}

	// Example 3: Error handling
	_, err = userService.GetUser(ctx, "invalid")
	if err != nil {
		logger.Warn("Expected error occurred", zap.Error(err))
	}

	// Example 4: Performance monitoring
	_, _ = userService.GetUser(ctx, "slow")

	// Example 5: Concurrent processing
	userIDs := []string{"user1", "user2", "user3", "invalid", "user5"}
	userService.ProcessUsers(ctx, userIDs)

	// Example 6: Using the sugar logger for simpler syntax
	sugar := logger.GetSugar()
	sugar.Infow("Sugar logger example",
		"feature", "sugared-logging",
		"performance", "high",
		"ease_of_use", true,
	)

	// Example 7: Benchmarking different logging approaches
	logger.Info("Starting logging performance comparison")

	start := time.Now()
	for i := 0; i < 1000; i++ {
		logger.Debug("Structured logging",
			zap.Int("iteration", i),
			zap.String("type", "structured"),
		)
	}
	structuredDuration := time.Since(start)

	start = time.Now()
	for i := 0; i < 1000; i++ {
		logger.Debugf("Formatted logging iteration %d type %s", i, "formatted")
	}
	formattedDuration := time.Since(start)

	logger.Info("Logging performance comparison completed",
		zap.Duration("structured_duration", structuredDuration),
		zap.Duration("formatted_duration", formattedDuration),
		zap.Float64("structured_ns_per_op", float64(structuredDuration.Nanoseconds())/1000),
		zap.Float64("formatted_ns_per_op", float64(formattedDuration.Nanoseconds())/1000),
	)

	// Example 8: Logging with different levels based on conditions
	for i := 0; i < 10; i++ {
		switch {
		case i < 3:
			logger.Debug("Debug level log", zap.Int("iteration", i))
		case i < 6:
			logger.Info("Info level log", zap.Int("iteration", i))
		case i < 8:
			logger.Warn("Warn level log", zap.Int("iteration", i))
		default:
			logger.Error("Error level log", zap.Int("iteration", i))
		}
	}

	// Always sync before exit
	defer logger.Sync()
}
