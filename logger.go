package wordsubgen

import (
	"fmt"
	"log"
	"os"
)

// Logger interface defines the logging methods that can be implemented by different loggers
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
}

// Field represents a key-value pair for structured logging
type Field struct {
	Key   string
	Value any
}

// NewField creates a new Field
func NewField(key string, value any) Field {
	return Field{Key: key, Value: value}
}

// DefaultLogger implements Logger using Go's standard log package
type DefaultLogger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	errorLogger *log.Logger
	warnLogger  *log.Logger
}

// NewDefaultLogger creates a new DefaultLogger
func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		debugLogger: log.New(os.Stdout, "[DEBUG] ", log.LstdFlags),
		infoLogger:  log.New(os.Stdout, "[INFO] ", log.LstdFlags),
		errorLogger: log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
		warnLogger:  log.New(os.Stdout, "[WARN] ", log.LstdFlags),
	}
}

// Debug logs a debug message
func (l *DefaultLogger) Debug(msg string, fields ...Field) {
	l.log(l.debugLogger, "DEBUG", msg, fields...)
}

// Info logs an info message
func (l *DefaultLogger) Info(msg string, fields ...Field) {
	l.log(l.infoLogger, "INFO", msg, fields...)
}

// Error logs an error message
func (l *DefaultLogger) Error(msg string, fields ...Field) {
	l.log(l.errorLogger, "ERROR", msg, fields...)
}

// Warn logs a warning message
func (l *DefaultLogger) Warn(msg string, fields ...Field) {
	l.log(l.warnLogger, "WARN", msg, fields...)
}

// log formats and logs a message with fields
func (l *DefaultLogger) log(logger *log.Logger, _, msg string, fields ...Field) {
	if len(fields) == 0 {
		logger.Printf("%s", msg)
		return
	}

	// Format fields as key=value pairs
	fieldStr := ""
	for i, field := range fields {
		if i > 0 {
			fieldStr += " "
		}
		fieldStr += fmt.Sprintf("%s=%v", field.Key, field.Value)
	}

	logger.Printf("%s %s", msg, fieldStr)
}

// NoOpLogger is a no-operation logger that discards all log messages
type NoOpLogger struct{}

// NewNoOpLogger creates a new NoOpLogger
func NewNoOpLogger() *NoOpLogger {
	return &NoOpLogger{}
}

// Debug does nothing
func (n *NoOpLogger) Debug(msg string, fields ...Field) {}

// Info does nothing
func (n *NoOpLogger) Info(msg string, fields ...Field) {}

// Error does nothing
func (n *NoOpLogger) Error(msg string, fields ...Field) {}

// Warn does nothing
func (n *NoOpLogger) Warn(msg string, fields ...Field) {}
