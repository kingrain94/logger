package logger

import (
	"testing"
)

func TestLogging(t *testing.T) {
	Debug("This is a debug level log.")
	Info("This is an information level log.")
	Warn("This is a warning level log.")
}
