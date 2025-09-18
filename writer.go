package wordsubgen

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// WriteASS writes ASS content to a file
func WriteASS(filename string, content string, logger Logger) error {
	logger.Debug("Writing ASS file", NewField("filename", filename), NewField("content_length", len(content)))

	if filename == "" {
		return errors.New("filename cannot be empty")
	}

	if content == "" {
		return errors.New("content cannot be empty")
	}

	// Create or truncate the file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer func() {
		_ = file.Close()
	}()

	// Write content to file
	bytesWritten, err := file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write content to file %s: %w", filename, err)
	}

	// Ensure data is written to disk
	err = file.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync file %s: %w", filename, err)
	}

	logger.Info("Successfully wrote ASS file", NewField("filename", filename), NewField("bytes_written", bytesWritten))
	return nil
}

// ReadLinesFromFile reads lines from a text file
func ReadLinesFromFile(filename string, logger Logger) ([]string, error) {
	logger.Debug("Reading lines from file", NewField("filename", filename))

	if filename == "" {
		return nil, errors.New("filename cannot be empty")
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	// Split by newlines and filter empty lines
	lines := strings.Split(string(content), "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return nil, errors.New("no valid lines found in file")
	}

	logger.Info("Successfully read lines from file", NewField("filename", filename), NewField("lines_read", len(result)))
	return result, nil
}
