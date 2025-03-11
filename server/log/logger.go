package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

// Level represents the severity level of a log message
type Level int

const (
	// DebugLevel for detailed troubleshooting information
	DebugLevel Level = iota
	// InfoLevel for general operational information
	InfoLevel
	// WarnLevel for potentially harmful situations
	WarnLevel
	// ErrorLevel for error events that might still allow the application to continue
	ErrorLevel
	// FatalLevel for severe error events that will lead the application to abort
	FatalLevel
)

// LevelNames maps log levels to their string representations
var LevelNames = map[Level]string{
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
	FatalLevel: "FATAL",
}

// LevelColors maps log levels to their ANSI color codes
var LevelColors = map[Level]string{
	DebugLevel: colorBlue,
	InfoLevel:  colorGreen,
	WarnLevel:  colorYellow,
	ErrorLevel: colorRed,
	FatalLevel: colorRed,
}

// Logger represents a logger with configurable output and level
type Logger struct {
	level      Level
	logger     *log.Logger
	showCaller bool
	useColors  bool
}

// Global logger instance with default configuration
var defaultLogger = NewLogger(os.Stdout, InfoLevel, true)

// NewLogger creates a new logger with the specified output, level, and caller display option
func NewLogger(out io.Writer, level Level, showCaller bool) *Logger {
	// Enable colors by default if output is a terminal
	useColors := isTerminal(out)

	return &Logger{
		level:      level,
		logger:     log.New(out, "", 0),
		showCaller: showCaller,
		useColors:  useColors,
	}
}

// isTerminal checks if the writer is likely a terminal that supports colors
func isTerminal(w io.Writer) bool {
	// Check if it's stdout or stderr
	return w == os.Stdout || w == os.Stderr
}

// SetLevel sets the minimum log level for the logger
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

// SetOutput sets the output destination for the logger
func (l *Logger) SetOutput(out io.Writer) {
	l.logger = log.New(out, "", 0)
	l.useColors = isTerminal(out)
}

// SetShowCaller sets whether to show the caller information in log messages
func (l *Logger) SetShowCaller(show bool) {
	l.showCaller = show
}

// SetUseColors sets whether to use colors in log messages
func (l *Logger) SetUseColors(use bool) {
	l.useColors = use
}

// formatMessage formats a log message with timestamp, level, and optional caller info
func (l *Logger) formatMessage(level Level, message string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	levelName := LevelNames[level]

	var colorStart, colorEnd string
	if l.useColors {
		colorStart = LevelColors[level]
		colorEnd = colorReset
	}

	var callerInfo string
	if l.showCaller {
		_, file, line, ok := runtime.Caller(3) // Skip 3 frames to get the actual caller
		if ok {
			// Extract just the filename without the full path
			parts := strings.Split(file, "/")
			file = parts[len(parts)-1]
			callerInfo = fmt.Sprintf(" [%s:%d]", file, line)
		}
	}

	return fmt.Sprintf("[%s] %s[%s]%s%s %s", timestamp, colorStart, levelName, colorEnd, callerInfo, message)
}

// log logs a message at the specified level if it's at or above the logger's configured level
func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	var message string
	if len(args) > 0 {
		message = fmt.Sprintf(format, args...)
	} else {
		message = format
	}

	formattedMessage := l.formatMessage(level, message)
	l.logger.Println(formattedMessage)

	// If this is a fatal message, exit the application
	if level == FatalLevel {
		os.Exit(1)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DebugLevel, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(InfoLevel, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WarnLevel, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ErrorLevel, format, args...)
}

// Fatal logs a fatal message and exits the application
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FatalLevel, format, args...)
}

// Global logger functions

// SetLevel sets the log level for the default logger
func SetLevel(level Level) {
	defaultLogger.SetLevel(level)
}

// SetOutput sets the output for the default logger
func SetOutput(out io.Writer) {
	defaultLogger.SetOutput(out)
}

// SetShowCaller sets whether to show caller info for the default logger
func SetShowCaller(show bool) {
	defaultLogger.SetShowCaller(show)
}

// SetUseColors sets whether to use colors for the default logger
func SetUseColors(use bool) {
	defaultLogger.SetUseColors(use)
}

// Debug logs a debug message to the default logger
func Debug(format string, args ...interface{}) {
	defaultLogger.Debug(format, args...)
}

// Info logs an info message to the default logger
func Info(format string, args ...interface{}) {
	defaultLogger.Info(format, args...)
}

// Warn logs a warning message to the default logger
func Warn(format string, args ...interface{}) {
	defaultLogger.Warn(format, args...)
}

// Error logs an error message to the default logger
func Error(format string, args ...interface{}) {
	defaultLogger.Error(format, args...)
}

// Fatal logs a fatal message to the default logger and exits the application
func Fatal(format string, args ...interface{}) {
	defaultLogger.Fatal(format, args...)
}
