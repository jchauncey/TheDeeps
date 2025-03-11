package log

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Contains(t, output, "[DEBUG]", "Debug message not logged")
	assert.Contains(t, output, "[INFO]", "Info message not logged")
	assert.Contains(t, output, "[WARN]", "Warning message not logged")
	assert.Contains(t, output, "[ERROR]", "Error message not logged")

	// Test log level filtering
	buf.Reset()
	logger.SetLevel(InfoLevel)

	logger.Debug("This debug message should be filtered out")
	logger.Info("This info message should be logged")

	output = buf.String()
	assert.NotContains(t, output, "[DEBUG]", "Debug message was logged when it should have been filtered")
	assert.Contains(t, output, "[INFO]", "Info message not logged")

	// Test formatted messages
	buf.Reset()
	logger.Info("Formatted message: %s, %d", "test", 123)

	output = buf.String()
	assert.Contains(t, output, "Formatted message: test, 123", "Formatted message not logged correctly")

	// Test caller info
	buf.Reset()
	logger.SetShowCaller(true)
	logger.Info("Message with caller info")

	output = buf.String()
	assert.Contains(t, output, "[INFO]", "INFO level tag not included in log message")
	assert.Contains(t, output, "[logger_test.go:", "Caller info not included in log message")

	// Test color output
	buf.Reset()
	logger.SetUseColors(true)
	logger.Info("Colored message")

	output = buf.String()
	assert.Contains(t, output, colorGreen, "Green color code not included in log message")
	assert.Contains(t, output, colorReset, "Color reset code not included in log message")
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
	assert.Contains(t, output, "[DEBUG]", "Debug message not logged")
	assert.Contains(t, output, "[INFO]", "Info message not logged")
	assert.Contains(t, output, "[WARN]", "Warning message not logged")
	assert.Contains(t, output, "[ERROR]", "Error message not logged")

	// Reset the default logger to stdout
	SetOutput(nil)
}
