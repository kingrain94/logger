package logger

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestEnvironmentString(t *testing.T) {
	tests := []struct {
		env      Environment
		expected string
	}{
		{Development, "development"},
		{Test, "test"},
		{Staging, "staging"},
		{Production, "production"},
		{Environment(99), "unknown"},
	}

	for _, tt := range tests {
		if got := tt.env.String(); got != tt.expected {
			t.Errorf("Environment.String() = %v, want %v", got, tt.expected)
		}
	}
}

func TestDefaultConfig(t *testing.T) {
	tests := []struct {
		env      Environment
		expected Config
	}{
		{
			Development,
			Config{
				Environment: Development,
				Level:       zapcore.DebugLevel,
				OutputPaths: []string{"stdout"},
				Encoding:    "console",
			},
		},
		{
			Test,
			Config{
				Environment: Test,
				Level:       zapcore.ErrorLevel,
				OutputPaths: []string{},
				Encoding:    "json",
			},
		},
		{
			Staging,
			Config{
				Environment: Staging,
				Level:       zapcore.InfoLevel,
				OutputPaths: []string{"stdout"},
				Encoding:    "json",
			},
		},
		{
			Production,
			Config{
				Environment: Production,
				Level:       zapcore.WarnLevel,
				OutputPaths: []string{"stdout"},
				Encoding:    "json",
			},
		},
	}

	for _, tt := range tests {
		got := DefaultConfig(tt.env)
		if got.Environment != tt.expected.Environment {
			t.Errorf("DefaultConfig(%v).Environment = %v, want %v", tt.env, got.Environment, tt.expected.Environment)
		}
		if got.Level != tt.expected.Level {
			t.Errorf("DefaultConfig(%v).Level = %v, want %v", tt.env, got.Level, tt.expected.Level)
		}
		if got.Encoding != tt.expected.Encoding {
			t.Errorf("DefaultConfig(%v).Encoding = %v, want %v", tt.env, got.Encoding, tt.expected.Encoding)
		}
	}
}

func TestInitialize(t *testing.T) {
	// Test successful initialization
	config := Config{
		Environment: Test,
		Level:       zapcore.InfoLevel,
		OutputPaths: []string{},
		Encoding:    "json",
	}

	err := Initialize(config)
	if err != nil {
		t.Fatalf("Initialize() error = %v, want nil", err)
	}

	// Verify logger is not nil
	if GetLogger() == nil {
		t.Error("GetLogger() returned nil after initialization")
	}

	// Verify sugar logger is not nil
	if GetSugar() == nil {
		t.Error("GetSugar() returned nil after initialization")
	}
}

func TestSetEnvironment(t *testing.T) {
	originalEnv := currentEnv

	// Test setting different environments
	environments := []Environment{Development, Test, Staging, Production}

	for _, env := range environments {
		err := SetEnvironment(env)
		if err != nil {
			t.Errorf("SetEnvironment(%v) error = %v, want nil", env, err)
		}

		if currentEnv != env {
			t.Errorf("currentEnv = %v, want %v", currentEnv, env)
		}
	}

	// Restore original environment
	SetEnvironment(originalEnv)
}

func TestSetLevel(t *testing.T) {
	// Initialize with a known configuration
	Initialize(DefaultConfig(Test))

	// Test setting different levels
	levels := []zapcore.Level{
		zapcore.DebugLevel,
		zapcore.InfoLevel,
		zapcore.WarnLevel,
		zapcore.ErrorLevel,
	}

	for _, level := range levels {
		SetLevel(level)
		// Note: We can't easily test if the level was actually set without
		// checking internal state, but we can at least verify no panic occurs
	}
}

func TestLoggingFunctions(t *testing.T) {
	// Capture output for testing
	var buf bytes.Buffer

	// Create a custom logger that writes to our buffer
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	zapConfig.Encoding = "json"

	// Create a core that writes to our buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapConfig.EncoderConfig),
		zapcore.AddSync(&buf),
		zapConfig.Level,
	)
	testLogger := zap.New(core)

	// Replace the global logger temporarily
	mu.Lock()
	oldLogger := logger
	oldSugar := sugar
	logger = testLogger
	sugar = testLogger.Sugar()
	mu.Unlock()

	// Test all logging functions
	Debug("debug message", zap.String("key", "value"))
	Debugf("debug formatted %s", "message")
	Info("info message", zap.String("key", "value"))
	Infof("info formatted %s", "message")
	Warn("warn message", zap.String("key", "value"))
	Warnf("warn formatted %s", "message")
	Error("error message", zap.String("key", "value"))
	Errorf("error formatted %s", "message")

	// Restore original logger
	mu.Lock()
	logger = oldLogger
	sugar = oldSugar
	mu.Unlock()

	// Check that logs were written
	output := buf.String()
	if output == "" {
		t.Error("No log output captured")
	}

	// Verify JSON structure of logs
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}
		var logEntry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &logEntry); err != nil {
			t.Errorf("Line %d: Invalid JSON log entry: %v", i+1, err)
		}

		// Check required fields
		if _, ok := logEntry["level"]; !ok {
			t.Errorf("Line %d: Missing 'level' field", i+1)
		}
		if _, ok := logEntry["msg"]; !ok {
			t.Errorf("Line %d: Missing 'msg' field", i+1)
		}
		if _, ok := logEntry["ts"]; !ok {
			t.Errorf("Line %d: Missing 'ts' field", i+1)
		}
	}
}

func TestWith(t *testing.T) {
	// Initialize logger
	Initialize(DefaultConfig(Test))

	// Test With function
	childLogger := With(zap.String("service", "test"))
	if childLogger == nil {
		t.Error("With() returned nil")
	}

	// Test WithFields function (alias)
	childLogger2 := WithFields(zap.String("service", "test"))
	if childLogger2 == nil {
		t.Error("WithFields() returned nil")
	}
}

func TestSync(t *testing.T) {
	// Initialize logger
	Initialize(DefaultConfig(Test))

	// Test Sync function
	err := Sync()
	if err != nil {
		// Sync can return an error for stdout/stderr, which is expected in tests
		// We just want to make sure it doesn't panic
		t.Logf("Sync() returned error (expected): %v", err)
	}
}

func TestConcurrentAccess(t *testing.T) {
	// Test concurrent access to logger functions
	Initialize(DefaultConfig(Test))

	done := make(chan bool, 10)

	// Start multiple goroutines that use the logger
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				Debug("concurrent debug", zap.Int("goroutine", id), zap.Int("iteration", j))
				Info("concurrent info", zap.Int("goroutine", id), zap.Int("iteration", j))
				Warn("concurrent warn", zap.Int("goroutine", id), zap.Int("iteration", j))
				Error("concurrent error", zap.Int("goroutine", id), zap.Int("iteration", j))
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestNilLoggerHandling(t *testing.T) {
	// Temporarily set logger to nil to test nil handling
	mu.Lock()
	oldLogger := logger
	oldSugar := sugar
	logger = nil
	sugar = nil
	mu.Unlock()

	// These should not panic
	Debug("test")
	Info("test")
	Warn("test")
	Error("test")
	Debugf("test %s", "formatted")
	Infof("test %s", "formatted")
	Warnf("test %s", "formatted")
	Errorf("test %s", "formatted")

	// Test other functions
	if GetLogger() != nil {
		t.Error("GetLogger() should return nil when logger is nil")
	}
	if GetSugar() != nil {
		t.Error("GetSugar() should return nil when sugar is nil")
	}
	if With(zap.String("test", "value")) != nil {
		t.Error("With() should return nil when logger is nil")
	}

	// Restore logger
	mu.Lock()
	logger = oldLogger
	sugar = oldSugar
	mu.Unlock()
}

// Benchmark tests
func BenchmarkDebug(b *testing.B) {
	Initialize(DefaultConfig(Test))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Debug("benchmark debug message", zap.Int("iteration", i))
	}
}

func BenchmarkInfo(b *testing.B) {
	Initialize(DefaultConfig(Test))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info("benchmark info message", zap.Int("iteration", i))
	}
}

func BenchmarkError(b *testing.B) {
	Initialize(DefaultConfig(Test))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Error("benchmark error message", zap.Int("iteration", i))
	}
}

func BenchmarkDebugf(b *testing.B) {
	Initialize(DefaultConfig(Test))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Debugf("benchmark debug message %d", i)
	}
}

// Example tests
func ExampleDebug() {
	Initialize(DefaultConfig(Development))
	Debug("This is a debug message", zap.String("key", "value"))
}

func ExampleInfo() {
	Initialize(DefaultConfig(Development))
	Info("Application started", zap.String("version", "1.0.0"))
}

func ExampleWith() {
	Initialize(DefaultConfig(Development))
	serviceLogger := With(zap.String("service", "user-auth"))
	serviceLogger.Info("User logged in", zap.String("user_id", "12345"))
}
