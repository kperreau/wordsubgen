package wordsubgen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteASS(t *testing.T) {
	// Create a temporary file
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test.ass")
	content := "[Script Info]\nTitle: Test\n\n[V4+ Styles]\nFormat: Name, Fontname\n\n[Events]\nFormat: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text\n"
	logger := NewNoOpLogger()

	err := WriteASS(filename, content, logger)
	if err != nil {
		t.Errorf("WriteASS failed: %v", err)
	}

	// Check that file was created and contains expected content
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		t.Errorf("Failed to read created file: %v", err)
	}

	if string(fileContent) != content {
		t.Errorf("File content mismatch. Expected:\n%s\nGot:\n%s", content, string(fileContent))
	}
}

func TestWriteASSWithEmptyFilename(t *testing.T) {
	logger := NewNoOpLogger()
	err := WriteASS("", "test content", logger)
	if err == nil {
		t.Error("Expected error for empty filename")
	}
	if !strings.Contains(err.Error(), "filename cannot be empty") {
		t.Errorf("Expected 'filename cannot be empty' error, got: %v", err)
	}
}

func TestWriteASSWithEmptyContent(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test.ass")
	logger := NewNoOpLogger()

	err := WriteASS(filename, "", logger)
	if err == nil {
		t.Error("Expected error for empty content")
	}
	if !strings.Contains(err.Error(), "content cannot be empty") {
		t.Errorf("Expected 'content cannot be empty' error, got: %v", err)
	}
}

func TestWriteASSWithInvalidPath(t *testing.T) {
	// Try to write to a directory that doesn't exist
	filename := "/nonexistent/path/test.ass"
	content := "test content"
	logger := NewNoOpLogger()

	err := WriteASS(filename, content, logger)
	if err == nil {
		t.Error("Expected error for invalid path")
	}
	if !strings.Contains(err.Error(), "failed to create file") {
		t.Errorf("Expected 'failed to create file' error, got: %v", err)
	}
}

func TestWriteASSOverwritesExistingFile(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test.ass")
	logger := NewNoOpLogger()

	// Write initial content
	initialContent := "initial content"
	err := WriteASS(filename, initialContent, logger)
	if err != nil {
		t.Errorf("WriteASS failed: %v", err)
	}

	// Write new content to same file
	newContent := "new content"
	err = WriteASS(filename, newContent, logger)
	if err != nil {
		t.Errorf("WriteASS failed: %v", err)
	}

	// Check that file contains new content
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		t.Errorf("Failed to read file: %v", err)
	}

	if string(fileContent) != newContent {
		t.Errorf("File was not overwritten. Expected: %s, Got: %s", newContent, string(fileContent))
	}
}

func TestReadLinesFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test.txt")
	logger := NewNoOpLogger()

	// Create test file with multiple lines
	testContent := "Line 1\nLine 2\nLine 3\n"
	err := os.WriteFile(filename, []byte(testContent), 0o644)
	if err != nil {
		t.Errorf("Failed to create test file: %v", err)
	}

	lines, err := ReadLinesFromFile(filename, logger)
	if err != nil {
		t.Errorf("ReadLinesFromFile failed: %v", err)
	}

	expectedLines := []string{"Line 1", "Line 2", "Line 3"}
	if len(lines) != len(expectedLines) {
		t.Errorf("Expected %d lines, got %d", len(expectedLines), len(lines))
	}

	for i, line := range lines {
		if line != expectedLines[i] {
			t.Errorf("Line %d: expected %q, got %q", i, expectedLines[i], line)
		}
	}
}

func TestReadLinesFromFileWithEmptyLines(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test.txt")
	logger := NewNoOpLogger()

	// Create test file with empty lines and whitespace
	testContent := "Line 1\n\nLine 2\n   \nLine 3\n"
	err := os.WriteFile(filename, []byte(testContent), 0o644)
	if err != nil {
		t.Errorf("Failed to create test file: %v", err)
	}

	lines, err := ReadLinesFromFile(filename, logger)
	if err != nil {
		t.Errorf("ReadLinesFromFile failed: %v", err)
	}

	expectedLines := []string{"Line 1", "Line 2", "Line 3"}
	if len(lines) != len(expectedLines) {
		t.Errorf("Expected %d lines, got %d", len(expectedLines), len(lines))
	}

	for i, line := range lines {
		if line != expectedLines[i] {
			t.Errorf("Line %d: expected %q, got %q", i, expectedLines[i], line)
		}
	}
}

func TestReadLinesFromFileWithEmptyFilename(t *testing.T) {
	logger := NewNoOpLogger()
	_, err := ReadLinesFromFile("", logger)
	if err == nil {
		t.Error("Expected error for empty filename")
	}
	if !strings.Contains(err.Error(), "filename cannot be empty") {
		t.Errorf("Expected 'filename cannot be empty' error, got: %v", err)
	}
}

func TestReadLinesFromFileWithNonexistentFile(t *testing.T) {
	logger := NewNoOpLogger()
	_, err := ReadLinesFromFile("nonexistent.txt", logger)
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
	if !strings.Contains(err.Error(), "failed to read file") {
		t.Errorf("Expected 'failed to read file' error, got: %v", err)
	}
}

func TestReadLinesFromFileWithEmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "empty.txt")
	logger := NewNoOpLogger()

	// Create empty file
	err := os.WriteFile(filename, []byte(""), 0o644)
	if err != nil {
		t.Errorf("Failed to create empty file: %v", err)
	}

	_, err = ReadLinesFromFile(filename, logger)
	if err == nil {
		t.Error("Expected error for empty file")
	}
	if !strings.Contains(err.Error(), "no valid lines found") {
		t.Errorf("Expected 'no valid lines found' error, got: %v", err)
	}
}

func TestReadLinesFromFileWithOnlyWhitespace(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "whitespace.txt")
	logger := NewNoOpLogger()

	// Create file with only whitespace
	testContent := "   \n\n\t\n   \n"
	err := os.WriteFile(filename, []byte(testContent), 0o644)
	if err != nil {
		t.Errorf("Failed to create whitespace file: %v", err)
	}

	_, err = ReadLinesFromFile(filename, logger)
	if err == nil {
		t.Error("Expected error for file with only whitespace")
	}
	if !strings.Contains(err.Error(), "no valid lines found") {
		t.Errorf("Expected 'no valid lines found' error, got: %v", err)
	}
}

func TestReadLinesFromFileWithTrailingNewline(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test.txt")
	logger := NewNoOpLogger()

	// Create test file with trailing newline
	testContent := "Line 1\nLine 2\nLine 3\n"
	err := os.WriteFile(filename, []byte(testContent), 0o644)
	if err != nil {
		t.Errorf("Failed to create test file: %v", err)
	}

	lines, err := ReadLinesFromFile(filename, logger)
	if err != nil {
		t.Errorf("ReadLinesFromFile failed: %v", err)
	}

	expectedLines := []string{"Line 1", "Line 2", "Line 3"}
	if len(lines) != len(expectedLines) {
		t.Errorf("Expected %d lines, got %d", len(expectedLines), len(lines))
	}
}

func TestReadLinesFromFileWithWindowsLineEndings(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test.txt")
	logger := NewNoOpLogger()

	// Create test file with Windows line endings
	testContent := "Line 1\r\nLine 2\r\nLine 3\r\n"
	err := os.WriteFile(filename, []byte(testContent), 0o644)
	if err != nil {
		t.Errorf("Failed to create test file: %v", err)
	}

	lines, err := ReadLinesFromFile(filename, logger)
	if err != nil {
		t.Errorf("ReadLinesFromFile failed: %v", err)
	}

	expectedLines := []string{"Line 1", "Line 2", "Line 3"}
	if len(lines) != len(expectedLines) {
		t.Errorf("Expected %d lines, got %d", len(expectedLines), len(lines))
	}

	for i, line := range lines {
		if line != expectedLines[i] {
			t.Errorf("Line %d: expected %q, got %q", i, expectedLines[i], line)
		}
	}
}

func TestWriteASSAndReadLinesIntegration(t *testing.T) {
	tmpDir := t.TempDir()
	assFile := filepath.Join(tmpDir, "test.ass")
	txtFile := filepath.Join(tmpDir, "input.txt")
	logger := NewNoOpLogger()

	// Create input file
	inputContent := "Hello world|Test subtitle|Another line"
	err := os.WriteFile(txtFile, []byte(inputContent), 0o644)
	if err != nil {
		t.Errorf("Failed to create input file: %v", err)
	}

	// Read lines from input file
	lines, err := ReadLinesFromFile(txtFile, logger)
	if err != nil {
		t.Errorf("ReadLinesFromFile failed: %v", err)
	}

	// Generate ASS content
	cfg := DefaultConfig()
	assContent, err := GenerateASS(cfg, lines)
	if err != nil {
		t.Errorf("GenerateASS failed: %v", err)
	}

	// Write ASS file
	err = WriteASS(assFile, assContent, logger)
	if err != nil {
		t.Errorf("WriteASS failed: %v", err)
	}

	// Verify the ASS file was created and contains expected content
	fileContent, err := os.ReadFile(assFile)
	if err != nil {
		t.Errorf("Failed to read ASS file: %v", err)
	}

	if !strings.Contains(string(fileContent), "Hello") {
		t.Error("ASS file should contain 'Hello'")
	}
	if !strings.Contains(string(fileContent), "world") {
		t.Error("ASS file should contain 'world'")
	}
	if !strings.Contains(string(fileContent), "Test") {
		t.Error("ASS file should contain 'Test'")
	}
	if !strings.Contains(string(fileContent), "subtitle") {
		t.Error("ASS file should contain 'subtitle'")
	}
	if !strings.Contains(string(fileContent), "Another") {
		t.Error("ASS file should contain 'Another'")
	}
	if !strings.Contains(string(fileContent), "line") {
		t.Error("ASS file should contain 'line'")
	}
}
