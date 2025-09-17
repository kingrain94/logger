# Go Logger

[![Go Version](https://img.shields.io/github/go-mod/go-version/kingrain94/logger)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/kingrain94/logger)](https://goreportcard.com/report/github.com/kingrain94/logger)
[![GoDoc](https://godoc.org/github.com/kingrain94/logger?status.svg)](https://godoc.org/github.com/kingrain94/logger)

A simple, efficient, and production-ready logging library for Go applications, built on top of [Uber's Zap](https://github.com/uber-go/zap). This logger provides structured logging with multiple environments support, thread-safe operations, and high performance.

## Features

- **High Performance**: Built on top of Uber's Zap logger
- **Multiple Log Levels**: DEBUG, INFO, WARN, ERROR, FATAL
- **Environment Support**: Development, Test, Staging, Production configurations
- **Thread-Safe**: Concurrent access with mutex protection
- **Structured Logging**: Support for structured fields and formatted messages
- **Flexible Configuration**: Customizable output paths, encoding, and log levels
- **Zero Dependencies**: Only depends on Zap (and its dependencies)
- **Easy Integration**: Simple API with sensible defaults

## Installation

```bash
go get github.com/kingrain94/logger
```

## Quick Start

```go
package main

import (
    "github.com/kingrain94/logger"
    "go.uber.org/zap"
)

func main() {
    // Basic usage with default configuration
    logger.Info("Application started", zap.String("version", "1.0.0"))
    logger.Debug("Debug information", zap.Int("user_id", 12345))
    logger.Warn("This is a warning", zap.String("component", "auth"))
    logger.Error("An error occurred", zap.Error(fmt.Errorf("connection failed")))
    
    // Formatted logging
    logger.Infof("User %s logged in from %s", "john_doe", "192.168.1.1")
    logger.Errorf("Failed to process request: %v", err)
}
```

## Configuration

### Environment-based Configuration

The logger supports four predefined environments:

```go
// Set environment (automatically configures appropriate settings)
logger.SetEnvironment(logger.Development) // Console output, debug level
logger.SetEnvironment(logger.Test)        // No output by default
logger.SetEnvironment(logger.Staging)     // JSON output, info level  
logger.SetEnvironment(logger.Production)  // JSON output, warn level
```

### Custom Configuration

```go
config := logger.Config{
    Environment: logger.Production,
    Level:       zapcore.InfoLevel,
    OutputPaths: []string{"stdout", "/var/log/app.log"},
    Encoding:    "json", // or "console"
}

err := logger.Initialize(config)
if err != nil {
    log.Fatal("Failed to initialize logger:", err)
}
```

### Dynamic Level Setting

```go
// Change log level at runtime
logger.SetLevel(zapcore.DebugLevel)
```

## API Reference

### Basic Logging Functions

```go
// Structured logging with fields
logger.Debug(msg string, fields ...zap.Field)
logger.Info(msg string, fields ...zap.Field)
logger.Warn(msg string, fields ...zap.Field)
logger.Error(msg string, fields ...zap.Field)
logger.Fatal(msg string, fields ...zap.Field) // Calls os.Exit(1)

// Formatted logging (printf-style)
logger.Debugf(template string, args ...interface{})
logger.Infof(template string, args ...interface{})
logger.Warnf(template string, args ...interface{})
logger.Errorf(template string, args ...interface{})
logger.Fatalf(template string, args ...interface{}) // Calls os.Exit(1)
```

### Advanced Usage

```go
// Create child logger with additional context
serviceLogger := logger.With(zap.String("service", "user-auth"))
serviceLogger.Info("User authenticated", zap.String("user_id", "12345"))

// Get underlying zap logger for advanced usage
zapLogger := logger.GetLogger()
zapLogger.Info("Direct zap usage")

// Get sugared logger for easier usage
sugar := logger.GetSugar()
sugar.Infow("User login", "user_id", 12345, "ip", "192.168.1.1")

// Flush buffered logs (important for graceful shutdowns)
logger.Sync()
```

## Examples

### Web Application

```go
package main

import (
    "net/http"
    "github.com/kingrain94/logger"
    "go.uber.org/zap"
)

func main() {
    // Initialize for production
    logger.SetEnvironment(logger.Production)
    
    // Create request logger with common fields
    requestLogger := logger.With(
        zap.String("service", "web-api"),
        zap.String("version", "1.0.0"),
    )
    
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        requestLogger.Info("Request received",
            zap.String("method", r.Method),
            zap.String("path", r.URL.Path),
            zap.String("remote_addr", r.RemoteAddr),
        )
        
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Hello World"))
    })
    
    logger.Info("Server starting", zap.Int("port", 8080))
    http.ListenAndServe(":8080", nil)
}
```

### Error Handling

```go
func processUser(userID string) error {
    logger.Debug("Processing user", zap.String("user_id", userID))
    
    user, err := getUserFromDB(userID)
    if err != nil {
        logger.Error("Failed to get user from database",
            zap.String("user_id", userID),
            zap.Error(err),
        )
        return err
    }
    
    if user.IsBlocked {
        logger.Warn("Blocked user attempted access",
            zap.String("user_id", userID),
            zap.String("reason", user.BlockReason),
        )
        return errors.New("user is blocked")
    }
    
    logger.Info("User processed successfully", zap.String("user_id", userID))
    return nil
}
```

### Testing

```go
func TestMyFunction(t *testing.T) {
    // Set test environment (minimal logging)
    logger.SetEnvironment(logger.Test)
    
    // Your test code here
    result := myFunction()
    
    // Assertions...
}
```

## Environment Details

| Environment | Log Level | Encoding | Output | Use Case |
|-------------|-----------|----------|---------|----------|
| Development | DEBUG | console | stdout | Local development |
| Test | ERROR | json | none | Unit/integration tests |
| Staging | INFO | json | stdout | Pre-production testing |
| Production | WARN | json | stdout | Production deployment |

## Performance

This logger is built on Zap, which is one of the fastest structured logging libraries for Go:

- **Zero allocation** for most log levels when disabled
- **Structured logging** without reflection
- **Thread-safe** operations with minimal overhead

Benchmark results (on a typical development machine):
```
BenchmarkDebug-8    20000000    85.4 ns/op    0 B/op    0 allocs/op
BenchmarkInfo-8     10000000    152 ns/op     0 B/op    0 allocs/op
BenchmarkError-8    10000000    165 ns/op     0 B/op    0 allocs/op
```

## Thread Safety

All logging functions are thread-safe and can be called concurrently from multiple goroutines without any additional synchronization.

## Best Practices

1. **Use structured logging** with fields instead of formatted strings when possible:
   ```go
   // Good
   logger.Info("User login", zap.String("user_id", userID), zap.String("ip", clientIP))
   
   // Less optimal
   logger.Infof("User %s logged in from %s", userID, clientIP)
   ```

2. **Create child loggers** for components with common context:
   ```go
   dbLogger := logger.With(zap.String("component", "database"))
   ```

3. **Always call Sync()** before application shutdown:
   ```go
   defer logger.Sync()
   ```

4. **Use appropriate log levels**:
   - `DEBUG`: Detailed information for debugging
   - `INFO`: General information about application flow
   - `WARN`: Warning conditions that should be addressed
   - `ERROR`: Error conditions that don't stop the application
   - `FATAL`: Critical errors that cause application termination

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built on top of [Uber's Zap](https://github.com/uber-go/zap)
- Inspired by best practices from the Go community
