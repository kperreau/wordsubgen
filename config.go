package wordsubgen

import (
	"errors"
	"fmt"
	"strings"
)

// Config holds all configuration options for generating ASS subtitle files
type Config struct {
	// Logger for the application
	Logger Logger `json:"-"`

	// Video dimensions
	Width  int `json:"width"`
	Height int `json:"height"`

	// Font settings
	FontName string `json:"font_name"`
	FontSize int    `json:"font_size"`

	// Colors (in ASS format: &H00BBGGRR)
	PrimaryColor   string `json:"primary_color"`   // Text color
	SecondaryColor string `json:"secondary_color"` // Karaoke color
	OutlineColor   string `json:"outline_color"`   // Outline color
	BackColor      string `json:"back_color"`      // Background color

	// Style settings
	Bold      bool `json:"bold"`
	Italic    bool `json:"italic"`
	Underline bool `json:"underline"`
	StrikeOut bool `json:"strike_out"`

	// Scaling and spacing
	ScaleX  int `json:"scale_x"`
	ScaleY  int `json:"scale_y"`
	Spacing int `json:"spacing"`
	Angle   int `json:"angle"`

	// Border and shadow
	BorderStyle int `json:"border_style"`
	Outline     int `json:"outline"`
	Shadow      int `json:"shadow"`

	// Shadow effects
	ShadowEnabled bool `json:"shadow_enabled"` // Enable/disable shadow effect
	ShadowX       int  `json:"shadow_x"`       // Horizontal shadow offset
	ShadowY       int  `json:"shadow_y"`       // Vertical shadow offset

	// Alignment (1=bottom-left, 2=bottom-center, 3=bottom-right, etc.)
	Alignment int `json:"alignment"`

	// Margins
	MarginL int `json:"margin_l"`
	MarginR int `json:"margin_r"`
	MarginV int `json:"margin_v"`

	// Timing settings
	StartDelay   int `json:"start_delay"`    // Milliseconds to delay start of subtitles
	PerWordDelay int `json:"per_word_delay"` // Milliseconds between words
	FadeDuration int `json:"fade_duration"`  // Milliseconds for fade effect
	LineHold     int `json:"line_hold"`      // Milliseconds to hold line after last word
	LineGap      int `json:"line_gap"`       // Milliseconds between lines

	// Karaoke mode
	Karaoke bool `json:"karaoke"`

	// Encoding
	Encoding int `json:"encoding"`
}

// DefaultConfig returns a config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Logger: NewDefaultLogger(),
		Width:  1080,
		Height: 1920,

		FontName: "Arial",
		FontSize: 64,

		PrimaryColor:   "&H00FFFFFF", // White
		SecondaryColor: "&H0000FFFF", // Yellow
		OutlineColor:   "&H00000000", // Black
		BackColor:      "&H64000000", // Semi-transparent black

		Bold:      true,
		Italic:    false,
		Underline: false,
		StrikeOut: false,

		ScaleX:  100,
		ScaleY:  100,
		Spacing: 0,
		Angle:   0,

		BorderStyle: 1,
		Outline:     4,
		Shadow:      0,

		ShadowEnabled: false,
		ShadowX:       3, // Default horizontal offset
		ShadowY:       8, // Default vertical offset (more downward)

		Alignment: 2, // Bottom-left

		MarginL: 40,
		MarginR: 40,
		MarginV: 120,

		StartDelay:   0,    // No delay by default
		PerWordDelay: 300,  // 300ms between words
		FadeDuration: 140,  // 140ms fade
		LineHold:     2000, // 2s hold
		LineGap:      0,    // No gap between lines

		Karaoke: false,

		Encoding: 1,
	}
}

// Validate checks if the config is valid and returns any errors
func (c *Config) Validate() error {
	c.Logger.Debug("Validating configuration")

	if c.Width <= 0 || c.Height <= 0 {
		return errors.New("width and height must be positive")
	}
	if c.FontSize <= 0 {
		return errors.New("font size must be positive")
	}
	if c.StartDelay < 0 {
		return errors.New("start delay cannot be negative")
	}
	if c.PerWordDelay < 0 {
		return errors.New("per word delay cannot be negative")
	}
	if c.FadeDuration < 0 {
		return errors.New("fade duration cannot be negative")
	}
	if c.LineHold < 0 {
		return errors.New("line hold cannot be negative")
	}
	if c.LineGap < 0 {
		return errors.New("line gap cannot be negative")
	}
	if c.Alignment < 1 || c.Alignment > 9 {
		return errors.New("alignment must be between 1 and 9")
	}
	if c.ShadowX < -50 || c.ShadowX > 50 {
		return errors.New("shadow X offset must be between -50 and 50")
	}
	if c.ShadowY < -50 || c.ShadowY > 50 {
		return errors.New("shadow Y offset must be between -50 and 50")
	}

	c.Logger.Debug("Configuration validation successful")
	return nil
}

// ColorToASS converts a hex color string to ASS format
func ColorToASS(color string) string {
	// Remove # if present
	color = strings.TrimPrefix(color, "#")

	// If it's a 6-digit hex, convert to ASS format
	if len(color) == 6 {
		// ASS format is &H00BBGGRR
		// Convert from RRGGBB to BBGGRR
		rr := color[0:2]
		gg := color[2:4]
		bb := color[4:6]
		return fmt.Sprintf("&H00%s%s%s", bb, gg, rr)
	}

	// If it's an 8-digit hex with alpha, convert to ASS format
	if len(color) == 8 {
		// ASS format is &HAABBGGRR
		// Convert from AARRGGBB to AABBGGRR
		aa := color[0:2]
		rr := color[2:4]
		gg := color[4:6]
		bb := color[6:8]
		return fmt.Sprintf("&H%s%s%s%s", aa, bb, gg, rr)
	}

	// If it's already in ASS format, return as is
	if strings.HasPrefix(color, "&H") {
		return color
	}

	// Default to white if invalid
	return "&H00FFFFFF"
}

// ColorToHex converts an ASS color string to hex format
func ColorToHex(color string) string {
	// Remove &H prefix if present
	color = strings.TrimPrefix(color, "&H")

	// If it's already a hex color without &H, return with #
	if len(color) == 6 && !strings.HasPrefix(color, "#") {
		return "#" + color
	}

	// If it's already a hex color with #, return as is
	if strings.HasPrefix(color, "#") {
		return color
	}

	// If it's an 8-digit ASS color, convert from AABBGGRR to #AARRGGBB
	if len(color) == 8 {
		aa := color[0:2]
		bb := color[2:4]
		gg := color[4:6]
		rr := color[6:8]
		return fmt.Sprintf("#%s%s%s%s", aa, rr, gg, bb)
	}

	// If it's a 6-digit ASS color, convert from BBGGRR to #RRGGBB
	if len(color) == 6 {
		bb := color[0:2]
		gg := color[2:4]
		rr := color[4:6]
		return fmt.Sprintf("#%s%s%s", rr, gg, bb)
	}

	// Default to white if invalid
	return "#FFFFFF"
}

// BoolToInt converts boolean to int (0 or 1)
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
