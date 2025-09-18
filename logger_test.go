package wordsubgen

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestNewField(t *testing.T) {
	field := NewField("key", "value")
	if field.Key != "key" {
		t.Errorf("Expected Key='key', got %s", field.Key)
	}
	if field.Value != "value" {
		t.Errorf("Expected Value='value', got %v", field.Value)
	}
}

func TestDefaultLogger(t *testing.T) {
	logger := NewDefaultLogger()
	if logger == nil {
		t.Error("NewDefaultLogger returned nil")
		return
	}

	// Test that all loggers are initialized
	if logger.debugLogger == nil {
		t.Error("debugLogger is nil")
	}
	if logger.infoLogger == nil {
		t.Error("infoLogger is nil")
	}
	if logger.errorLogger == nil {
		t.Error("errorLogger is nil")
	}
	if logger.warnLogger == nil {
		t.Error("warnLogger is nil")
	}
}

func TestDefaultLoggerLogging(t *testing.T) {
	logger := NewDefaultLogger()

	// Test that logging methods don't panic
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	// Test with fields
	logger.Debug("debug with fields", NewField("key", "value"))
	logger.Info("info with fields", NewField("key", "value"))
	logger.Warn("warn with fields", NewField("key", "value"))
	logger.Error("error with fields", NewField("key", "value"))
}

func TestDefaultLoggerWithFields(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	defer func() {
		os.Stdout = oldStdout
	}()

	stdoutR, stdoutW, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create stdout pipe: %v", err)
	}
	os.Stdout = stdoutW

	logger := NewDefaultLogger()
	logger.Info("test message", NewField("key1", "value1"), NewField("key2", 42))

	_ = stdoutW.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(stdoutR)
	output := buf.String()

	if !strings.Contains(output, "test message") {
		t.Error("Message not found in output")
	}
	if !strings.Contains(output, "key1=value1") {
		t.Error("Field key1=value1 not found in output")
	}
	if !strings.Contains(output, "key2=42") {
		t.Error("Field key2=42 not found in output")
	}
}

func TestDefaultLoggerWithMultipleFields(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	defer func() {
		os.Stdout = oldStdout
	}()

	stdoutR, stdoutW, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create stdout pipe: %v", err)
	}
	os.Stdout = stdoutW

	logger := NewDefaultLogger()
	logger.Debug("complex message",
		NewField("user", "john"),
		NewField("action", "login"),
		NewField("timestamp", 1234567890),
		NewField("success", true))

	_ = stdoutW.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(stdoutR)
	output := buf.String()

	expectedFields := []string{
		"user=john",
		"action=login",
		"timestamp=1234567890",
		"success=true",
	}

	for _, field := range expectedFields {
		if !strings.Contains(output, field) {
			t.Errorf("Field %s not found in output", field)
		}
	}
}

func TestNoOpLogger(t *testing.T) {
	logger := NewNoOpLogger()
	if logger == nil {
		t.Error("NewNoOpLogger returned nil")
	}

	// These should not panic or cause any issues
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	// Test with fields
	logger.Debug("debug message", NewField("key", "value"))
	logger.Info("info message", NewField("key", "value"))
	logger.Warn("warn message", NewField("key", "value"))
	logger.Error("error message", NewField("key", "value"))
}

func TestLoggerInterface(t *testing.T) {
	// Test that both logger types implement the Logger interface
	var logger1 Logger = NewDefaultLogger()
	var logger2 Logger = NewNoOpLogger()

	// These should compile without issues
	logger1.Debug("test")
	logger1.Info("test")
	logger1.Warn("test")
	logger1.Error("test")

	logger2.Debug("test")
	logger2.Info("test")
	logger2.Warn("test")
	logger2.Error("test")
}

func TestDefaultLoggerTimestamp(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	defer func() {
		os.Stdout = oldStdout
	}()

	stdoutR, stdoutW, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create stdout pipe: %v", err)
	}
	os.Stdout = stdoutW

	logger := NewDefaultLogger()
	logger.Info("timestamp test")

	_ = stdoutW.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(stdoutR)
	output := buf.String()

	// Check that timestamp is present (log.LstdFlags includes date and time)
	// The format should be something like "2023/01/01 12:00:00"
	if !strings.Contains(output, "/") {
		t.Error("Expected timestamp with date format in output")
	}
	if !strings.Contains(output, ":") {
		t.Error("Expected timestamp with time format in output")
	}
}

func TestDefaultLoggerDifferentOutputs(t *testing.T) {
	// Test that error messages go to stderr while others go to stdout
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	defer func() {
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()

	// Create pipes
	stdoutR, stdoutW, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create stdout pipe: %v", err)
	}
	stderrR, stderrW, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create stderr pipe: %v", err)
	}

	os.Stdout = stdoutW
	os.Stderr = stderrW

	logger := NewDefaultLogger()

	// Log to different levels
	logger.Info("info message")
	logger.Debug("debug message")
	logger.Warn("warn message")
	logger.Error("error message")

	// Close write ends
	_ = stdoutW.Close()
	_ = stderrW.Close()

	// Read outputs
	var stdoutBuf, stderrBuf bytes.Buffer
	_, _ = stdoutBuf.ReadFrom(stdoutR)
	_, _ = stderrBuf.ReadFrom(stderrR)

	stdoutOutput := stdoutBuf.String()
	stderrOutput := stderrBuf.String()

	// Check that info, debug, and warn go to stdout
	if !strings.Contains(stdoutOutput, "info message") {
		t.Error("Info message should go to stdout")
	}
	if !strings.Contains(stdoutOutput, "debug message") {
		t.Error("Debug message should go to stdout")
	}
	if !strings.Contains(stdoutOutput, "warn message") {
		t.Error("Warn message should go to stdout")
	}

	// Check that error goes to stderr
	if !strings.Contains(stderrOutput, "error message") {
		t.Error("Error message should go to stderr")
	}

	// Check that error doesn't go to stdout
	if strings.Contains(stdoutOutput, "error message") {
		t.Error("Error message should not go to stdout")
	}
}

func TestFieldTypes(t *testing.T) {
	// Test different field value types
	tests := []struct {
		key   string
		value any
	}{
		{"string", "test"},
		{"int", 42},
		{"float", 3.14},
		{"bool", true},
		{"bool_false", false},
		{"nil", nil},
	}

	for _, test := range tests {
		field := NewField(test.key, test.value)
		if field.Key != test.key {
			t.Errorf("Field key mismatch: expected %s, got %s", test.key, field.Key)
		}
		if field.Value != test.value {
			t.Errorf("Field value mismatch: expected %v, got %v", test.value, field.Value)
		}
	}
}
