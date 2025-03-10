package log

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer

	// Create a new logger that writes to the buffer
	logger := NewLogger(&buf, DebugLevel, false)

	// Disable colors for testing
	logger.SetUseColors(false)

	// Test each log level
	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")

	// Check that all messages were logged
	output := buf.String()
	if !strings.Contains(output, "[DEBUG]") {
		t.Error("Debug message not logged")
	}
	if !strings.Contains(output, "[INFO]") {
		t.Error("Info message not logged")
	}
	if !strings.Contains(output, "[WARN]") {
		t.Error("Warning message not logged")
	}
	if !strings.Contains(output, "[ERROR]") {
		t.Error("Error message not logged")
	}

	// Test log level filtering
	buf.Reset()
	logger.SetLevel(InfoLevel)

	logger.Debug("This debug message should be filtered out")
	logger.Info("This info message should be logged")

	output = buf.String()
	if strings.Contains(output, "[DEBUG]") {
		t.Error("Debug message was logged when it should have been filtered")
	}
	if !strings.Contains(output, "[INFO]") {
		t.Error("Info message not logged")
	}

	// Test formatted messages
	buf.Reset()
	logger.Info("Formatted message: %s, %d", "test", 123)

	output = buf.String()
	if !strings.Contains(output, "Formatted message: test, 123") {
		t.Error("Formatted message not logged correctly")
	}

	// Test caller info
	buf.Reset()
	logger.SetShowCaller(true)
	logger.Info("Message with caller info")

	output = buf.String()
	if !strings.Contains(output, "[INFO]") || !strings.Contains(output, "[logger_test.go:") {
		t.Error("Caller info not included in log message")
	}

	// Test color output
	buf.Reset()
	logger.SetUseColors(true)
	logger.Info("Colored message")

	output = buf.String()
	if !strings.Contains(output, colorGreen) || !strings.Contains(output, colorReset) {
		t.Error("Color codes not included in log message")
	}
}

func TestGlobalLogger(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer

	// Set the default logger to write to the buffer
	SetOutput(&buf)
	SetLevel(DebugLevel)
	SetUseColors(false)

	// Test each log level
	Debug("This is a debug message")
	Info("This is an info message")
	Warn("This is a warning message")
	Error("This is an error message")

	// Check that all messages were logged
	output := buf.String()
	if !strings.Contains(output, "[DEBUG]") {
		t.Error("Debug message not logged")
	}
	if !strings.Contains(output, "[INFO]") {
		t.Error("Info message not logged")
	}
	if !strings.Contains(output, "[WARN]") {
		t.Error("Warning message not logged")
	}
	if !strings.Contains(output, "[ERROR]") {
		t.Error("Error message not logged")
	}

	// Reset the default logger to stdout
	SetOutput(nil)
}
