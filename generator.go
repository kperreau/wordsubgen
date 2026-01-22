package wordsubgen

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// GenerateASS generates ASS subtitle content from configuration and text lines
func GenerateASS(cfg *Config, lines []string) (string, error) {
	cfg.Logger.Debug("Starting ASS generation", NewField("lines_count", len(lines)))

	if err := cfg.Validate(); err != nil {
		return "", fmt.Errorf("invalid config: %w", err)
	}

	if len(lines) == 0 {
		return "", errors.New("no lines provided")
	}

	var content strings.Builder

	// Write ASS header
	content.WriteString("[Script Info]\n")
	content.WriteString("Title: Generated Subtitle\n")
	content.WriteString("ScriptType: v4.00+\n")
	content.WriteString("Collisions: Normal\n")
	content.WriteString(fmt.Sprintf("PlayResX: %d\n", cfg.Width))
	content.WriteString(fmt.Sprintf("PlayResY: %d\n", cfg.Height))
	content.WriteString("\n")

	// Write styles section
	content.WriteString("[V4+ Styles]\n")
	content.WriteString("Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding\n")

	styleLine := fmt.Sprintf("Style: Default,%s,%d,%s,%s,%s,%s,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d\n",
		cfg.FontName,
		cfg.FontSize,
		cfg.PrimaryColor,
		cfg.SecondaryColor,
		cfg.OutlineColor,
		cfg.BackColor,
		BoolToInt(cfg.Bold),
		BoolToInt(cfg.Italic),
		BoolToInt(cfg.Underline),
		BoolToInt(cfg.StrikeOut),
		cfg.ScaleX,
		cfg.ScaleY,
		cfg.Spacing,
		cfg.Angle,
		cfg.BorderStyle,
		cfg.Outline,
		cfg.Shadow,
		cfg.Alignment,
		cfg.MarginL,
		cfg.MarginR,
		cfg.MarginV,
		cfg.Encoding,
	)
	content.WriteString(styleLine)
	content.WriteString("\n")

	// Write events section
	content.WriteString("[Events]\n")
	content.WriteString("Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text\n")
	content.WriteString("\n")

	// Generate dialogue lines
	currentTime := time.Duration(cfg.StartDelay) * time.Millisecond
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			cfg.Logger.Debug("Skipping empty line", NewField("line_number", i+1))
			continue
		}

		cfg.Logger.Debug("Generating dialogue line", NewField("line_number", i+1), NewField("text", line))
		dialogue, duration, err := generateDialogueLine(cfg, line, currentTime)
		if err != nil {
			return "", fmt.Errorf("failed to generate dialogue for line %d: %w", i+1, err)
		}

		content.WriteString(dialogue)
		content.WriteString("\n")

		// Update current time for next line
		currentTime += duration + time.Duration(cfg.LineGap)*time.Millisecond
		cfg.Logger.Debug("Generated dialogue line", NewField("line_number", i+1), NewField("duration", duration), NewField("total_time", currentTime))
	}

	cfg.Logger.Info("ASS generation completed successfully", NewField("total_lines", len(lines)), NewField("total_duration", currentTime))
	return content.String(), nil
}

// generateDialogueLine creates a single dialogue line with word-by-word fade effects
func generateDialogueLine(cfg *Config, text string, startTime time.Duration) (string, time.Duration, error) {
	words := strings.Fields(text)
	if len(words) == 0 {
		return "", 0, errors.New("empty line")
	}

	// Calculate timing
	wordDuration := time.Duration(cfg.PerWordDelay) * time.Millisecond
	fadeDuration := time.Duration(cfg.FadeDuration) * time.Millisecond
	holdDuration := time.Duration(cfg.LineHold) * time.Millisecond

	// Total duration: (words-1) * word_delay + fade + hold
	totalDuration := time.Duration(len(words)-1)*wordDuration + fadeDuration + holdDuration
	endTime := startTime + totalDuration

	// Format times in ASS format (h:mm:ss.cs)
	startTimeStr := formatASSTime(startTime)
	endTimeStr := formatASSTime(endTime)

	// Build the dialogue text with fade effects
	var dialogueText strings.Builder

	// Add alignment
	dialogueText.WriteString(fmt.Sprintf("{\\an%d}", cfg.Alignment))

	// Add shadow effect if enabled
	if cfg.ShadowEnabled {
		dialogueText.WriteString(fmt.Sprintf("{\\xshad%d\\yshad%d}", cfg.ShadowX, cfg.ShadowY))
	}

	// Add word-by-word fade effects
	for i, word := range words {
		wordStart := int(i * cfg.PerWordDelay)
		wordEnd := wordStart + cfg.FadeDuration

		// Add fade effect: start invisible, fade in
		dialogueText.WriteString(fmt.Sprintf("{\\alpha&HFF&\\t(%d,%d,\\alpha&H00&)}%s",
			wordStart, wordEnd, word))

		// Add space after word (except last word)
		if i < len(words)-1 {
			dialogueText.WriteString(" ")
		}
	}

	// Add karaoke effect if enabled
	if cfg.Karaoke {
		dialogueText.WriteString(fmt.Sprintf("{\\k%d}", cfg.PerWordDelay/10)) // karaoke duration in centiseconds
	}

	// Format the dialogue line
	dialogue := fmt.Sprintf("Dialogue: 0,%s,%s,Default,,0,0,0,,%s",
		startTimeStr, endTimeStr, dialogueText.String())

	return dialogue, totalDuration, nil
}

// formatASSTime converts a time.Duration to ASS time format (h:mm:ss.cs)
func formatASSTime(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	centiseconds := (d.Milliseconds() / 10) % 100

	return fmt.Sprintf("%d:%02d:%02d.%02d", hours, minutes, seconds, centiseconds)
}

// ParseLinesFromString parses lines from a pipe-separated string
func ParseLinesFromString(input string) []string {
	if input == "" {
		return nil
	}

	lines := strings.Split(input, "|")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// ParseLinesFromFile reads lines from a file (placeholder for file reading)
func ParseLinesFromFile(filename string) ([]string, error) {
	// This would typically read from a file
	// For now, return an error as this should be implemented with actual file reading
	return nil, errors.New("file reading not implemented yet")
}
