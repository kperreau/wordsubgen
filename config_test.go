package wordsubgen

import (
	"strings"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	// Test default values
	if cfg.Width != 1080 {
		t.Errorf("Expected Width=1080, got %d", cfg.Width)
	}
	if cfg.Height != 1920 {
		t.Errorf("Expected Height=1920, got %d", cfg.Height)
	}
	if cfg.FontName != "Arial" {
		t.Errorf("Expected FontName=Arial, got %s", cfg.FontName)
	}
	if cfg.FontSize != 64 {
		t.Errorf("Expected FontSize=64, got %d", cfg.FontSize)
	}
	if cfg.PrimaryColor != "&H00FFFFFF" {
		t.Errorf("Expected PrimaryColor=&H00FFFFFF, got %s", cfg.PrimaryColor)
	}
	if cfg.Alignment != 2 {
		t.Errorf("Expected Alignment=2, got %d", cfg.Alignment)
	}
	if cfg.StartDelay != 0 {
		t.Errorf("Expected StartDelay=0, got %d", cfg.StartDelay)
	}
	if cfg.PerWordDelay != 300 {
		t.Errorf("Expected PerWordDelay=300, got %d", cfg.PerWordDelay)
	}
	if cfg.Logger == nil {
		t.Error("Expected Logger to be initialized")
	}
}

func TestConfigValidation(t *testing.T) {
	cfg := DefaultConfig()

	// Test valid config
	err := cfg.Validate()
	if err != nil {
		t.Errorf("Expected no error for default config, got: %v", err)
	}

	// Test invalid width
	cfg.Width = 0
	err = cfg.Validate()
	if err == nil {
		t.Error("Expected error for width <= 0")
	}
	if !strings.Contains(err.Error(), "width and height must be positive") {
		t.Errorf("Expected error about width/height, got: %v", err)
	}

	// Test invalid height
	cfg.Width = 1080
	cfg.Height = -1
	err = cfg.Validate()
	if err == nil {
		t.Error("Expected error for height <= 0")
	}

	// Test invalid font size
	cfg.Height = 1920
	cfg.FontSize = 0
	err = cfg.Validate()
	if err == nil {
		t.Error("Expected error for font size <= 0")
	}
	if !strings.Contains(err.Error(), "font size must be positive") {
		t.Errorf("Expected error about font size, got: %v", err)
	}

	// Test invalid per word delay
	cfg.FontSize = 64
	cfg.PerWordDelay = -1
	err = cfg.Validate()
	if err == nil {
		t.Error("Expected error for negative per word delay")
	}
	if !strings.Contains(err.Error(), "per word delay cannot be negative") {
		t.Errorf("Expected error about per word delay, got: %v", err)
	}

	// Test invalid fade duration
	cfg.PerWordDelay = 300
	cfg.FadeDuration = -1
	err = cfg.Validate()
	if err == nil {
		t.Error("Expected error for negative fade duration")
	}

	// Test invalid line hold
	cfg.FadeDuration = 140
	cfg.LineHold = -1
	err = cfg.Validate()
	if err == nil {
		t.Error("Expected error for negative line hold")
	}

	// Test invalid line gap
	cfg.LineHold = 2000
	cfg.LineGap = -1
	err = cfg.Validate()
	if err == nil {
		t.Error("Expected error for negative line gap")
	}

	// Test invalid alignment
	cfg.LineGap = 0
	cfg.Alignment = 0
	err = cfg.Validate()
	if err == nil {
		t.Error("Expected error for alignment < 1")
	}
	if !strings.Contains(err.Error(), "alignment must be between 1 and 9") {
		t.Errorf("Expected error about alignment, got: %v", err)
	}

	cfg.Alignment = 10
	err = cfg.Validate()
	if err == nil {
		t.Error("Expected error for alignment > 9")
	}

	// Test invalid start delay
	cfg.Alignment = 2
	cfg.StartDelay = -1
	err = cfg.Validate()
	if err == nil {
		t.Error("Expected error for negative start delay")
	}
	if !strings.Contains(err.Error(), "start delay cannot be negative") {
		t.Errorf("Expected error about start delay, got: %v", err)
	}
}

func TestShadowValidation(t *testing.T) {
	// Test valid shadow values
	cfg := DefaultConfig()
	cfg.ShadowEnabled = true
	cfg.ShadowX = 3
	cfg.ShadowY = 8

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Expected no error for valid shadow values, got: %v", err)
	}

	// Test invalid shadow X values
	cfg.ShadowX = 51
	err = cfg.Validate()
	if err == nil {
		t.Error("Expected error for shadow X > 50, got nil")
	}
	if !strings.Contains(err.Error(), "shadow X offset must be between -50 and 50") {
		t.Errorf("Expected error about shadow X offset, got: %v", err)
	}

	cfg.ShadowX = -51
	err = cfg.Validate()
	if err == nil {
		t.Error("Expected error for shadow X < -50, got nil")
	}

	// Test invalid shadow Y values
	cfg.ShadowX = 3 // Reset to valid value
	cfg.ShadowY = 51
	err = cfg.Validate()
	if err == nil {
		t.Error("Expected error for shadow Y > 50, got nil")
	}
	if !strings.Contains(err.Error(), "shadow Y offset must be between -50 and 50") {
		t.Errorf("Expected error about shadow Y offset, got: %v", err)
	}

	cfg.ShadowY = -51
	err = cfg.Validate()
	if err == nil {
		t.Error("Expected error for shadow Y < -50, got nil")
	}
}

func TestColorToASS(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"#FFFFFF", "&H00FFFFFF"},
		{"#FF0000", "&H000000FF"},
		{"#00FF00", "&H0000FF00"},
		{"#0000FF", "&H00FF0000"},
		{"#123456", "&H00563412"},
		{"#FF0000FF", "&HFFFF0000"},  // AARRGGBB -> AABBGGRR
		{"#80FF0000", "&H800000FF"},  // AARRGGBB -> AABBGGRR
		{"&H00FFFFFF", "&H00FFFFFF"}, // Already in ASS format
		{"FFFFFF", "&H00FFFFFF"},     // Without #
		{"invalid", "&H00FFFFFF"},    // Invalid input defaults to white
		{"", "&H00FFFFFF"},           // Empty input defaults to white
	}

	for _, test := range tests {
		result := ColorToASS(test.input)
		if result != test.expected {
			t.Errorf("ColorToASS(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestColorToHex(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"&H00FFFFFF", "#00FFFFFF"}, // BBGGRR -> RRGGBB with alpha
		{"&H000000FF", "#00FF0000"}, // BBGGRR -> RRGGBB with alpha
		{"&H0000FF00", "#0000FF00"}, // BBGGRR -> RRGGBB with alpha
		{"&H00FF0000", "#000000FF"}, // BBGGRR -> RRGGBB with alpha
		{"&H00563412", "#00123456"}, // BBGGRR -> RRGGBB with alpha
		{"&HFF0000FF", "#FFFF0000"}, // AABBGGRR -> AARRGGBB
		{"&H800000FF", "#80FF0000"}, // AABBGGRR -> AARRGGBB
		{"#FFFFFF", "#FFFFFF"},      // Already in hex format
		{"FFFFFF", "#FFFFFF"},       // Without #
		{"invalid", "#FFFFFF"},      // Invalid input defaults to white
		{"", "#FFFFFF"},             // Empty input defaults to white
	}

	for _, test := range tests {
		result := ColorToHex(test.input)
		if result != test.expected {
			t.Errorf("ColorToHex(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestBoolToInt(t *testing.T) {
	if BoolToInt(true) != 1 {
		t.Error("BoolToInt(true) should return 1")
	}
	if BoolToInt(false) != 0 {
		t.Error("BoolToInt(false) should return 0")
	}
}

func TestShadowGeneration(t *testing.T) {
	// Test shadow enabled
	cfg := DefaultConfig()
	cfg.ShadowEnabled = true
	cfg.ShadowX = 3
	cfg.ShadowY = 8

	lines := []string{"Test shadow"}
	content, err := GenerateASS(cfg, lines)
	if err != nil {
		t.Errorf("Failed to generate ASS: %v", err)
	}

	// Check that shadow tags are present
	if !strings.Contains(content, "\\xshad3\\yshad8") {
		t.Error("Expected shadow tags \\xshad3\\yshad8 in generated content")
	}

	// Test shadow disabled
	cfg.ShadowEnabled = false
	content2, err := GenerateASS(cfg, lines)
	if err != nil {
		t.Errorf("Failed to generate ASS: %v", err)
	}

	// Check that shadow tags are not present
	if strings.Contains(content2, "\\xshad") {
		t.Error("Expected no shadow tags when shadow is disabled")
	}
}
