package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kingrain94/logger"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger for production
	logger.SetEnvironment(logger.Production)

	// Create a service logger with common fields
	serviceLogger := logger.With(
		zap.String("service", "web-api"),
		zap.String("version", "1.0.0"),
	)

	// Middleware for request logging
	loggingMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Log incoming request
			serviceLogger.Info("Request started",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
			)

			// Call the actual handler
			next(w, r)

			// Log request completion
			duration := time.Since(start)
			serviceLogger.Info("Request completed",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Duration("duration", duration),
			)
		}
	}

	// Home handler
	http.HandleFunc("/", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Processing home request")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hello, World! Time: %s", time.Now().Format(time.RFC3339))
	}))

	// Health check handler
	http.HandleFunc("/health", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Health check requested")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
	}))

	// Error simulation handler
	http.HandleFunc("/error", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		logger.Error("Simulated error occurred",
			zap.String("path", r.URL.Path),
			zap.String("error", "simulated database connection failed"),
		)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Internal server error"}`))
	}))

	// Start server
	port := 8080
	logger.Info("Starting web server",
		zap.Int("port", port),
		zap.String("environment", "production"),
	)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		logger.Fatal("Failed to start server",
			zap.Error(err),
			zap.Int("port", port),
		)
	}

	// Ensure logs are flushed on exit
	defer logger.Sync()
}
