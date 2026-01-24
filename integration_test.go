package wordsubgen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompleteWorkflow(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.txt")
	outputFile := filepath.Join(tmpDir, "output.ass")
	logger := NewNoOpLogger()

	// Create input file with pipe-separated lines
	inputContent := "Hello world|This is a test|Another subtitle line"
	err := os.WriteFile(inputFile, []byte(inputContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create input file: %v", err)
	}

	// Read content from input file and parse pipe-separated lines
	fileContent, err := os.ReadFile(inputFile)
	if err != nil {
		t.Fatalf("Failed to read input file: %v", err)
	}

	lines := ParseLinesFromString(string(fileContent))
	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}

	// Generate ASS content
	cfg := DefaultConfig()
	assContent, err := GenerateASS(cfg, lines)
	if err != nil {
		t.Fatalf("GenerateASS failed: %v", err)
	}

	// Write ASS file
	err = WriteASS(outputFile, assContent, logger)
	if err != nil {
		t.Fatalf("WriteASS failed: %v", err)
	}

	// Verify the output file
	fileContent, err = os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	content := string(fileContent)

	// Check ASS structure
	if !strings.Contains(content, "[Script Info]") {
		t.Error("Missing [Script Info] section")
	}
	if !strings.Contains(content, "[V4+ Styles]") {
		t.Error("Missing [V4+ Styles] section")
	}
	if !strings.Contains(content, "[Events]") {
		t.Error("Missing [Events] section")
	}

	// Check that all input lines are present (check individual words due to fade effects)
	expectedWords := []string{"Hello", "world", "This", "is", "a", "test", "Another", "subtitle", "line"}
	for _, word := range expectedWords {
		if !strings.Contains(content, word) {
			t.Errorf("Missing word '%s' in output", word)
		}
	}

	// Check dialogue count
	dialogueCount := strings.Count(content, "Dialogue:")
	if dialogueCount != len(lines) {
		t.Errorf("Expected %d dialogue lines, got %d", len(lines), dialogueCount)
	}
}

func TestWorkflowWithCustomConfig(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.txt")
	outputFile := filepath.Join(tmpDir, "output.ass")
	logger := NewNoOpLogger()

	// Create input file
	inputContent := "Custom styled text|With shadow effects"
	err := os.WriteFile(inputFile, []byte(inputContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create input file: %v", err)
	}

	// Read lines
	lines, err := ReadLinesFromFile(inputFile, logger)
	if err != nil {
		t.Fatalf("ReadLinesFromFile failed: %v", err)
	}

	// Create custom config
	cfg := DefaultConfig()
	cfg.FontName = "Times New Roman"
	cfg.FontSize = 72
	cfg.PrimaryColor = "&H00FF0000" // Red
	cfg.ShadowEnabled = true
	cfg.ShadowX = 5
	cfg.ShadowY = 10
	cfg.Karaoke = true
	cfg.PerWordDelay = 500
	cfg.FadeDuration = 200
	cfg.LineHold = 3000

	// Generate ASS content
	assContent, err := GenerateASS(cfg, lines)
	if err != nil {
		t.Fatalf("GenerateASS failed: %v", err)
	}

	// Write ASS file
	err = WriteASS(outputFile, assContent, logger)
	if err != nil {
		t.Fatalf("WriteASS failed: %v", err)
	}

	// Verify custom settings in output
	fileContent, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	content := string(fileContent)

	// Check custom font
	if !strings.Contains(content, "Times New Roman") {
		t.Error("Custom font name not found in output")
	}

	// Check custom color
	if !strings.Contains(content, "&H00FF0000") {
		t.Error("Custom primary color not found in output")
	}

	// Check shadow effects
	if !strings.Contains(content, "\\xshad5\\yshad10") {
		t.Error("Custom shadow effects not found in output")
	}

	// Check karaoke effects
	if !strings.Contains(content, "{\\k") {
		t.Error("Karaoke effects not found in output")
	}
}

func TestWorkflowWithEmptyLines(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.txt")
	outputFile := filepath.Join(tmpDir, "output.ass")
	logger := NewNoOpLogger()

	// Create input file with empty lines
	inputContent := "First line\n\nSecond line\n   \nThird line\n"
	err := os.WriteFile(inputFile, []byte(inputContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create input file: %v", err)
	}

	// Read lines (should filter out empty lines)
	lines, err := ReadLinesFromFile(inputFile, logger)
	if err != nil {
		t.Fatalf("ReadLinesFromFile failed: %v", err)
	}

	if len(lines) != 3 {
		t.Errorf("Expected 3 non-empty lines, got %d", len(lines))
	}

	// Generate ASS content
	cfg := DefaultConfig()
	assContent, err := GenerateASS(cfg, lines)
	if err != nil {
		t.Fatalf("GenerateASS failed: %v", err)
	}

	// Write ASS file
	err = WriteASS(outputFile, assContent, logger)
	if err != nil {
		t.Fatalf("WriteASS failed: %v", err)
	}

	// Verify output
	fileContent, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	content := string(fileContent)

	// Should contain all non-empty lines (check individual words due to fade effects)
	expectedWords := []string{"First", "line", "Second", "line", "Third", "line"}
	for _, word := range expectedWords {
		if !strings.Contains(content, word) {
			t.Errorf("Missing word '%s' in output", word)
		}
	}
}

func TestWorkflowWithDifferentAlignments(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.txt")
	logger := NewNoOpLogger()

	// Create input file
	inputContent := "Test alignment"
	err := os.WriteFile(inputFile, []byte(inputContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create input file: %v", err)
	}

	// Read lines
	lines, err := ReadLinesFromFile(inputFile, logger)
	if err != nil {
		t.Fatalf("ReadLinesFromFile failed: %v", err)
	}

	// Test different alignments
	alignments := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for _, alignment := range alignments {
		cfg := DefaultConfig()
		cfg.Alignment = alignment

		// Generate ASS content
		assContent, err := GenerateASS(cfg, lines)
		if err != nil {
			t.Fatalf("GenerateASS failed for alignment %d: %v", alignment, err)
		}

		// Check that alignment tag is present
		expectedTag := fmt.Sprintf("{\\an%d}", alignment)
		if !strings.Contains(assContent, expectedTag) {
			t.Errorf("Alignment tag %s not found for alignment %d", expectedTag, alignment)
		}
	}
}

func TestWorkflowWithLineGap(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.txt")
	outputFile := filepath.Join(tmpDir, "output.ass")
	logger := NewNoOpLogger()

	// Create input file with multiple lines
	inputContent := "First line\nSecond line\nThird line"
	err := os.WriteFile(inputFile, []byte(inputContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create input file: %v", err)
	}

	// Read lines
	lines, err := ReadLinesFromFile(inputFile, logger)
	if err != nil {
		t.Fatalf("ReadLinesFromFile failed: %v", err)
	}

	// Create config with line gap
	cfg := DefaultConfig()
	cfg.LineGap = 1000 // 1 second gap between lines

	// Generate ASS content
	assContent, err := GenerateASS(cfg, lines)
	if err != nil {
		t.Fatalf("GenerateASS failed: %v", err)
	}

	// Write ASS file
	err = WriteASS(outputFile, assContent, logger)
	if err != nil {
		t.Fatalf("WriteASS failed: %v", err)
	}

	// Verify output
	fileContent, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	content := string(fileContent)

	// Should contain all lines (check individual words due to fade effects)
	expectedWords := []string{"First", "line", "Second", "line", "Third", "line"}
	for _, word := range expectedWords {
		if !strings.Contains(content, word) {
			t.Errorf("Missing word '%s' in output", word)
		}
	}

	// Check dialogue count
	dialogueCount := strings.Count(content, "Dialogue:")
	if dialogueCount != len(lines) {
		t.Errorf("Expected %d dialogue lines, got %d", len(lines), dialogueCount)
	}
}

func TestWorkflowErrorHandling(t *testing.T) {
	logger := NewNoOpLogger()

	// Test with invalid config
	cfg := DefaultConfig()
	cfg.Width = 0 // Invalid
	lines := []string{"Test"}

	_, err := GenerateASS(cfg, lines)
	if err == nil {
		t.Error("Expected error for invalid config")
	}

	// Test with empty lines
	cfg = DefaultConfig()
	emptyLines := []string{}

	_, err = GenerateASS(cfg, emptyLines)
	if err == nil {
		t.Error("Expected error for empty lines")
	}

	// Test WriteASS with empty filename
	err = WriteASS("", "content", logger)
	if err == nil {
		t.Error("Expected error for empty filename")
	}

	// Test WriteASS with empty content
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.ass")
	err = WriteASS(outputFile, "", logger)
	if err == nil {
		t.Error("Expected error for empty content")
	}

	// Test ReadLinesFromFile with empty filename
	_, err = ReadLinesFromFile("", logger)
	if err == nil {
		t.Error("Expected error for empty filename")
	}

	// Test ReadLinesFromFile with nonexistent file
	_, err = ReadLinesFromFile("nonexistent.txt", logger)
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestWorkflowWithStructuredJSON(t *testing.T) {
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "words.json")
	outputFile := filepath.Join(tmpDir, "output.ass")
	logger := NewNoOpLogger()

	jsonContent := `{"segments":[{"words":[{"word":"Découvrez","start":0.162,"end":1.024},{"word":"cette","start":1.064,"end":1.344}]}]}`
	err := os.WriteFile(jsonFile, []byte(jsonContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create JSON file: %v", err)
	}

	data, err := os.ReadFile(jsonFile)
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	phrases, err := ParseStructuredJSON(data)
	if err != nil {
		t.Fatalf("ParseStructuredJSON failed: %v", err)
	}
	if len(phrases) != 1 || len(phrases[0].Words) != 2 {
		t.Errorf("Expected 1 phrase with 2 words, got %d phrases", len(phrases))
	}

	cfg := DefaultConfig()
	cfg.StartDelay = 1500 // 1.5 s shift

	content, err := GenerateASSFromStructured(cfg, phrases)
	if err != nil {
		t.Fatalf("GenerateASSFromStructured failed: %v", err)
	}

	err = WriteASS(outputFile, content, logger)
	if err != nil {
		t.Fatalf("WriteASS failed: %v", err)
	}

	out, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}
	str := string(out)

	if !strings.Contains(str, "Découvrez") || !strings.Contains(str, "cette") {
		t.Error("Expected words in output")
	}
	// 0.162 + 1.5 = 1.662s -> 0:00:01.66
	if !strings.Contains(str, "0:00:01.66") && !strings.Contains(str, "0:00:01.67") {
		// Allow for rounding: 1.662 could be 1.66 or 1.67 in ASS
		if !strings.Contains(str, "0:00:01.6") {
			t.Error("Expected delayed start time (~1.66s) in output")
		}
	}
}

func TestWorkflowWithSpecialCharacters(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.txt")
	outputFile := filepath.Join(tmpDir, "output.ass")
	logger := NewNoOpLogger()

	// Create input file with special characters
	inputContent := "Hello, world!|This has \"quotes\"|And 'apostrophes'|Numbers: 123|Symbols: @#$%"
	err := os.WriteFile(inputFile, []byte(inputContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create input file: %v", err)
	}

	// Read content from input file and parse pipe-separated lines
	fileContent, err := os.ReadFile(inputFile)
	if err != nil {
		t.Fatalf("Failed to read input file: %v", err)
	}

	lines := ParseLinesFromString(string(fileContent))

	// Generate ASS content
	cfg := DefaultConfig()
	assContent, err := GenerateASS(cfg, lines)
	if err != nil {
		t.Fatalf("GenerateASS failed: %v", err)
	}

	// Write ASS file
	err = WriteASS(outputFile, assContent, logger)
	if err != nil {
		t.Fatalf("WriteASS failed: %v", err)
	}

	// Verify output
	fileContent, err = os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	content := string(fileContent)

	// Should contain all lines with special characters (check individual words due to fade effects)
	expectedWords := []string{"Hello,", "world!", "This", "has", "\"quotes\"", "And", "'apostrophes'", "Numbers:", "123", "Symbols:", "@#$%"}
	for _, word := range expectedWords {
		if !strings.Contains(content, word) {
			t.Errorf("Missing word '%s' in output", word)
		}
	}
}
