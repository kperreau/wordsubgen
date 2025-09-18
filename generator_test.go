package wordsubgen

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestGenerateASS(t *testing.T) {
	cfg := DefaultConfig()
	lines := []string{"Hello world", "Test subtitle"}

	content, err := GenerateASS(cfg, lines)
	if err != nil {
		t.Errorf("GenerateASS failed: %v", err)
	}

	// Check that content contains expected sections
	if !strings.Contains(content, "[Script Info]") {
		t.Error("Expected [Script Info] section")
	}
	if !strings.Contains(content, "[V4+ Styles]") {
		t.Error("Expected [V4+ Styles] section")
	}
	if !strings.Contains(content, "[Events]") {
		t.Error("Expected [Events] section")
	}

	// Check that content contains dialogue lines
	if !strings.Contains(content, "Dialogue:") {
		t.Error("Expected Dialogue lines")
	}

	// Check that content contains our test text (with fade effects)
	if !strings.Contains(content, "Hello") {
		t.Error("Expected 'Hello' in content")
	}
	if !strings.Contains(content, "world") {
		t.Error("Expected 'world' in content")
	}
	if !strings.Contains(content, "Test") {
		t.Error("Expected 'Test' in content")
	}
	if !strings.Contains(content, "subtitle") {
		t.Error("Expected 'subtitle' in content")
	}
}

func TestGenerateASSWithInvalidConfig(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Width = 0 // Invalid config
	lines := []string{"Test"}

	_, err := GenerateASS(cfg, lines)
	if err == nil {
		t.Error("Expected error for invalid config")
	}
	if !strings.Contains(err.Error(), "invalid config") {
		t.Errorf("Expected 'invalid config' error, got: %v", err)
	}
}

func TestGenerateASSWithEmptyLines(t *testing.T) {
	cfg := DefaultConfig()
	lines := []string{}

	_, err := GenerateASS(cfg, lines)
	if err == nil {
		t.Error("Expected error for empty lines")
	}
	if !strings.Contains(err.Error(), "no lines provided") {
		t.Errorf("Expected 'no lines provided' error, got: %v", err)
	}
}

func TestGenerateASSWithEmptyLine(t *testing.T) {
	cfg := DefaultConfig()
	lines := []string{"Hello", "", "World"}

	content, err := GenerateASS(cfg, lines)
	if err != nil {
		t.Errorf("GenerateASS failed: %v", err)
	}

	// Should contain both non-empty lines
	if !strings.Contains(content, "Hello") {
		t.Error("Expected 'Hello' in content")
	}
	if !strings.Contains(content, "World") {
		t.Error("Expected 'World' in content")
	}
}

func TestGenerateDialogueLine(t *testing.T) {
	cfg := DefaultConfig()
	text := "Hello world"
	startTime := time.Duration(0)

	dialogue, duration, err := generateDialogueLine(cfg, text, startTime)
	if err != nil {
		t.Errorf("generateDialogueLine failed: %v", err)
	}

	// Check that dialogue contains expected elements
	if !strings.Contains(dialogue, "Dialogue:") {
		t.Error("Expected 'Dialogue:' in dialogue")
	}
	if !strings.Contains(dialogue, "Hello") {
		t.Error("Expected 'Hello' in dialogue")
	}
	if !strings.Contains(dialogue, "world") {
		t.Error("Expected 'world' in dialogue")
	}
	if !strings.Contains(dialogue, "{\\an2}") {
		t.Error("Expected alignment tag")
	}

	// Check duration calculation
	expectedDuration := time.Duration(len(strings.Fields(text))-1)*time.Duration(cfg.PerWordDelay)*time.Millisecond +
		time.Duration(cfg.FadeDuration)*time.Millisecond +
		time.Duration(cfg.LineHold)*time.Millisecond
	if duration != expectedDuration {
		t.Errorf("Expected duration %v, got %v", expectedDuration, duration)
	}
}

func TestGenerateDialogueLineWithEmptyText(t *testing.T) {
	cfg := DefaultConfig()
	text := ""
	startTime := time.Duration(0)

	_, _, err := generateDialogueLine(cfg, text, startTime)
	if err == nil {
		t.Error("Expected error for empty text")
	}
	if !strings.Contains(err.Error(), "empty line") {
		t.Errorf("Expected 'empty line' error, got: %v", err)
	}
}

func TestGenerateDialogueLineWithShadow(t *testing.T) {
	cfg := DefaultConfig()
	cfg.ShadowEnabled = true
	cfg.ShadowX = 3
	cfg.ShadowY = 8
	text := "Test shadow"
	startTime := time.Duration(0)

	dialogue, _, err := generateDialogueLine(cfg, text, startTime)
	if err != nil {
		t.Errorf("generateDialogueLine failed: %v", err)
	}

	// Check that shadow tags are present
	if !strings.Contains(dialogue, "\\xshad3\\yshad8") {
		t.Error("Expected shadow tags in dialogue")
	}
}

func TestGenerateDialogueLineWithKaraoke(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Karaoke = true
	text := "Test karaoke"
	startTime := time.Duration(0)

	dialogue, _, err := generateDialogueLine(cfg, text, startTime)
	if err != nil {
		t.Errorf("generateDialogueLine failed: %v", err)
	}

	// Check that karaoke tags are present
	if !strings.Contains(dialogue, "{\\k") {
		t.Error("Expected karaoke tags in dialogue")
	}
}

func TestFormatASSTime(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{0, "0:00:00.00"},
		{time.Second, "0:00:01.00"},
		{time.Minute, "0:01:00.00"},
		{time.Hour, "1:00:00.00"},
		{time.Hour + time.Minute + time.Second, "1:01:01.00"},
		{time.Millisecond * 150, "0:00:00.15"},
		{time.Millisecond * 1234, "0:00:01.23"},
		{time.Hour + time.Minute*30 + time.Second*45 + time.Millisecond*500, "1:30:45.50"},
	}

	for _, test := range tests {
		result := formatASSTime(test.duration)
		if result != test.expected {
			t.Errorf("formatASSTime(%v) = %s, expected %s", test.duration, result, test.expected)
		}
	}
}

func TestParseLinesFromString(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"", nil},
		{"Hello", []string{"Hello"}},
		{"Hello|World", []string{"Hello", "World"}},
		{"Hello|World|Test", []string{"Hello", "World", "Test"}},
		{" Hello | World ", []string{"Hello", "World"}},
		{"Hello||World", []string{"Hello", "World"}},
		{"|Hello|World|", []string{"Hello", "World"}},
		{"  |  Hello  |  World  |  ", []string{"Hello", "World"}},
	}

	for _, test := range tests {
		result := ParseLinesFromString(test.input)
		if len(result) != len(test.expected) {
			t.Errorf("ParseLinesFromString(%q) returned %d lines, expected %d", test.input, len(result), len(test.expected))
			continue
		}
		for i, line := range result {
			if line != test.expected[i] {
				t.Errorf("ParseLinesFromString(%q)[%d] = %q, expected %q", test.input, i, line, test.expected[i])
			}
		}
	}
}

func TestParseLinesFromFile(t *testing.T) {
	_, err := ParseLinesFromFile("nonexistent.txt")
	if err == nil {
		t.Error("Expected error for ParseLinesFromFile")
	}
	if !strings.Contains(err.Error(), "not implemented") {
		t.Errorf("Expected 'not implemented' error, got: %v", err)
	}
}

func TestGenerateASSWithMultipleLines(t *testing.T) {
	cfg := DefaultConfig()
	cfg.LineGap = 500 // 500ms gap between lines
	lines := []string{"First line", "Second line", "Third line"}

	content, err := GenerateASS(cfg, lines)
	if err != nil {
		t.Errorf("GenerateASS failed: %v", err)
	}

	// Count dialogue lines
	dialogueCount := strings.Count(content, "Dialogue:")
	if dialogueCount != 3 {
		t.Errorf("Expected 3 dialogue lines, got %d", dialogueCount)
	}

	// Check that all lines are present (check individual words due to fade effects)
	expectedWords := []string{"First", "line", "Second", "line", "Third", "line"}
	for _, word := range expectedWords {
		if !strings.Contains(content, word) {
			t.Errorf("Expected word '%s' in content", word)
		}
	}
}

func TestGenerateASSWithDifferentAlignments(t *testing.T) {
	lines := []string{"Test alignment"}

	alignments := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for _, alignment := range alignments {
		cfg := DefaultConfig()
		cfg.Alignment = alignment

		content, err := GenerateASS(cfg, lines)
		if err != nil {
			t.Errorf("GenerateASS failed for alignment %d: %v", alignment, err)
		}

		expectedTag := fmt.Sprintf("{\\an%d}", alignment)
		if !strings.Contains(content, expectedTag) {
			t.Errorf("Expected alignment tag %s in content for alignment %d", expectedTag, alignment)
		}
	}
}

func TestGenerateASSWithFadeEffects(t *testing.T) {
	cfg := DefaultConfig()
	cfg.PerWordDelay = 200
	cfg.FadeDuration = 100
	lines := []string{"Hello world"}

	content, err := GenerateASS(cfg, lines)
	if err != nil {
		t.Errorf("GenerateASS failed: %v", err)
	}

	// Check that fade effects are present
	if !strings.Contains(content, "{\\alpha&HFF&\\t(") {
		t.Error("Expected fade effect tags in content")
	}
	if !strings.Contains(content, "\\alpha&H00&") {
		t.Error("Expected fade-in effect in content")
	}
}

func TestGenerateASSWithStyleSettings(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Bold = true
	cfg.Italic = true
	cfg.Underline = true
	cfg.StrikeOut = true
	cfg.FontName = "Times New Roman"
	cfg.FontSize = 48
	cfg.PrimaryColor = "&H00FF0000"
	cfg.SecondaryColor = "&H0000FF00"
	cfg.OutlineColor = "&H000000FF"
	cfg.BackColor = "&H80FFFFFF"

	lines := []string{"Styled text"}
	content, err := GenerateASS(cfg, lines)
	if err != nil {
		t.Errorf("GenerateASS failed: %v", err)
	}

	// Check that style line contains our settings
	styleLine := fmt.Sprintf("Style: Default,%s,%d,%s,%s,%s,%s,%d,%d,%d,%d",
		cfg.FontName, cfg.FontSize, cfg.PrimaryColor, cfg.SecondaryColor,
		cfg.OutlineColor, cfg.BackColor, BoolToInt(cfg.Bold), BoolToInt(cfg.Italic),
		BoolToInt(cfg.Underline), BoolToInt(cfg.StrikeOut))

	if !strings.Contains(content, styleLine) {
		t.Error("Expected style line with custom settings in content")
	}
}
