# Examples

This directory contains practical examples of how to use the Go Logger library in different scenarios.

## Running Examples

Each example is a standalone Go program. To run an example:

```bash
cd examples/basic
go run main.go
```

## Available Examples

### 1. Basic Usage (`basic/`)

Demonstrates fundamental logging operations:
- Basic structured logging with fields
- Formatted logging (printf-style)
- Different log levels
- Multiple fields in a single log entry

**Run:**
```bash
cd basic && go run main.go
```

### 2. Web Server (`web-server/`)

Shows how to integrate logging into a web application:
- Request logging middleware
- Service-specific loggers with common fields
- Different handlers with appropriate log levels
- Production environment configuration

**Run:**
```bash
cd web-server && go run main.go
```

Then visit:
- `http://localhost:8080/` - Normal request
- `http://localhost:8080/health` - Health check
- `http://localhost:8080/error` - Error simulation

### 3. Configuration (`configuration/`)

Explores various configuration options:
- Environment-based configuration
- Custom configuration with different outputs
- Dynamic level changes
- File logging
- Multiple output destinations
- Environment variable based setup

**Run:**
```bash
cd configuration && go run main.go
```

### 4. Advanced Usage (`advanced/`)

Demonstrates sophisticated logging patterns:
- Service-oriented logging with dedicated loggers
- Context-aware logging with trace IDs
- Concurrent logging from multiple goroutines
- Error handling and performance monitoring
- Sugar logger for simplified syntax
- Performance comparisons between different logging styles

**Run:**
```bash
cd advanced && go run main.go
```

## Example Outputs

### Development Environment
```
2024-01-15T10:30:45.123Z    INFO    Application started    {"version": "1.0.0"}
2024-01-15T10:30:45.124Z    DEBUG   Debug information      {"user_id": 12345}
2024-01-15T10:30:45.125Z    WARN    This is a warning      {"component": "auth"}
```

### Production Environment (JSON)
```json
{"level":"info","ts":"2024-01-15T10:30:45.123Z","msg":"Application started","version":"1.0.0"}
{"level":"warn","ts":"2024-01-15T10:30:45.125Z","msg":"This is a warning","component":"auth"}
```

## Best Practices Demonstrated

1. **Structured Logging**: Use fields instead of formatted strings when possible
2. **Context Propagation**: Pass trace IDs and request context through logs
3. **Service Loggers**: Create dedicated loggers for different components
4. **Appropriate Levels**: Use correct log levels for different types of information
5. **Performance Considerations**: Compare structured vs formatted logging performance
6. **Error Handling**: Proper error logging with context
7. **Graceful Shutdown**: Always call `logger.Sync()` before exit

## Integration Patterns

### Middleware Pattern (Web Servers)
```go
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        logger.Info("Request started", /* fields */)
        next(w, r)
        logger.Info("Request completed", zap.Duration("duration", time.Since(start)))
    }
}
```

### Service Pattern (Business Logic)
```go
type UserService struct {
    logger *zap.Logger
}

func NewUserService() *UserService {
    return &UserService{
        logger: logger.With(zap.String("component", "user-service")),
    }
}
```

### Context Pattern (Distributed Systems)
```go
func processRequest(ctx context.Context) {
    traceID := getTraceID(ctx)
    logger.Info("Processing request", zap.String("trace_id", traceID))
}
```

## Performance Notes

The examples include performance comparisons showing:
- Structured logging is generally faster than formatted logging
- Zero allocations for disabled log levels
- Minimal overhead for concurrent logging
- Efficient field serialization

## Environment Variables

Some examples respond to environment variables:
- `APP_ENV`: Set to "development", "test", "staging", or "production"
- `LOG_LEVEL`: Set to "debug", "info", "warn", or "error"

Example:
```bash
APP_ENV=production LOG_LEVEL=info go run main.go
```
