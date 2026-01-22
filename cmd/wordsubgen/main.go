package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kperreau/wordsubgen"
)

func main() {
	// Setup logging
	logger := wordsubgen.NewDefaultLogger()

	// Get default configuration to use as flag defaults
	defaultCfg := wordsubgen.DefaultConfig()

	// Command line flags
	var (
		lines  = flag.String("lines", "", "Input lines separated by | (e.g., 'Hello world|Second line')")
		file   = flag.String("file", "", "Input file with lines (one per line)")
		output = flag.String("out", "output.ass", "Output ASS file path")

		// Video settings
		width  = flag.Int("width", defaultCfg.Width, "Video width")
		height = flag.Int("height", defaultCfg.Height, "Video height")

		// Font settings
		fontName = flag.String("font", defaultCfg.FontName, "Font name")
		fontSize = flag.Int("fontsize", defaultCfg.FontSize, "Font size")

		// Colors (hex format) - convert ASS format back to hex for flags
		primaryColor   = flag.String("color", wordsubgen.ColorToHex(defaultCfg.PrimaryColor), "Primary text color (hex)")
		secondaryColor = flag.String("secondary", wordsubgen.ColorToHex(defaultCfg.SecondaryColor), "Secondary color for karaoke (hex)")
		outlineColor   = flag.String("outline", wordsubgen.ColorToHex(defaultCfg.OutlineColor), "Outline color (hex)")
		backColor      = flag.String("background", defaultCfg.BackColor, "Background color (ASS format)")

		// Style settings
		bold      = flag.Bool("bold", defaultCfg.Bold, "Bold text")
		italic    = flag.Bool("italic", defaultCfg.Italic, "Italic text")
		underline = flag.Bool("underline", defaultCfg.Underline, "Underline text")
		strikeout = flag.Bool("strikeout", defaultCfg.StrikeOut, "Strikeout text")

		// Scaling
		scaleX  = flag.Int("scalex", defaultCfg.ScaleX, "Horizontal scale (%)")
		scaleY  = flag.Int("scaley", defaultCfg.ScaleY, "Vertical scale (%)")
		spacing = flag.Int("spacing", defaultCfg.Spacing, "Character spacing")
		angle   = flag.Int("angle", defaultCfg.Angle, "Text angle")

		// Border and shadow
		borderStyle  = flag.Int("borderstyle", defaultCfg.BorderStyle, "Border style")
		outlineWidth = flag.Int("outlinewidth", defaultCfg.Outline, "Outline width")
		shadow       = flag.Int("shadow", defaultCfg.Shadow, "Shadow width")

		// Shadow effects
		shadowEnabled = flag.Bool("shadow-enabled", defaultCfg.ShadowEnabled, "Enable shadow effect")
		shadowX       = flag.Int("shadow-x", defaultCfg.ShadowX, "Horizontal shadow offset")
		shadowY       = flag.Int("shadow-y", defaultCfg.ShadowY, "Vertical shadow offset")

		// Alignment
		alignment = flag.Int("alignment", defaultCfg.Alignment, "Text alignment (1=bottom-left, 2=bottom-center, 3=bottom-right, etc.)")

		// Margins
		marginL = flag.Int("marginl", defaultCfg.MarginL, "Left margin")
		marginR = flag.Int("marginr", defaultCfg.MarginR, "Right margin")
		marginV = flag.Int("marginv", defaultCfg.MarginV, "Vertical margin")

		// Timing
		startDelay   = flag.Int("start-delay", defaultCfg.StartDelay, "Delay before starting subtitles (ms)")
		perWordDelay = flag.Int("delay", defaultCfg.PerWordDelay, "Delay between words (ms)")
		fadeDuration = flag.Int("fade", defaultCfg.FadeDuration, "Fade duration (ms)")
		lineHold     = flag.Int("hold", defaultCfg.LineHold, "Line hold duration (ms)")
		lineGap      = flag.Int("gap", defaultCfg.LineGap, "Gap between lines (ms)")

		// Features
		karaoke = flag.Bool("karaoke", defaultCfg.Karaoke, "Enable karaoke mode")

		// Other
		help = flag.Bool("help", false, "Show help")
	)

	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// Set log level (for now, we'll use the default logger regardless of verbose flag)
	// In the future, we could implement different log levels in DefaultLogger

	// Validate input
	if *lines == "" && *file == "" {
		logger.Error("Either --lines or --file must be specified")
		showHelp()
		os.Exit(1)
	}

	if *lines != "" && *file != "" {
		logger.Error("Cannot specify both --lines and --file")
		os.Exit(1)
	}

	// Create configuration
	cfg := wordsubgen.DefaultConfig()

	// Setup logger for the library
	cfg.Logger = logger

	// Apply command line overrides
	cfg.Width = *width
	cfg.Height = *height
	cfg.FontName = *fontName
	cfg.FontSize = *fontSize
	cfg.PrimaryColor = wordsubgen.ColorToASS(*primaryColor)
	cfg.SecondaryColor = wordsubgen.ColorToASS(*secondaryColor)
	cfg.OutlineColor = wordsubgen.ColorToASS(*outlineColor)
	cfg.BackColor = wordsubgen.ColorToASS(*backColor)
	cfg.Bold = *bold
	cfg.Italic = *italic
	cfg.Underline = *underline
	cfg.StrikeOut = *strikeout
	cfg.ScaleX = *scaleX
	cfg.ScaleY = *scaleY
	cfg.Spacing = *spacing
	cfg.Angle = *angle
	cfg.BorderStyle = *borderStyle
	cfg.Outline = *outlineWidth
	cfg.Shadow = *shadow
	cfg.ShadowEnabled = *shadowEnabled
	cfg.ShadowX = *shadowX
	cfg.ShadowY = *shadowY
	cfg.Alignment = *alignment
	cfg.MarginL = *marginL
	cfg.MarginR = *marginR
	cfg.MarginV = *marginV
	cfg.StartDelay = *startDelay
	cfg.PerWordDelay = *perWordDelay
	cfg.FadeDuration = *fadeDuration
	cfg.LineHold = *lineHold
	cfg.LineGap = *lineGap
	cfg.Karaoke = *karaoke

	// Get input lines
	var inputLines []string
	var err error

	if *lines != "" {
		inputLines = wordsubgen.ParseLinesFromString(*lines)
		logger.Info("Parsed input lines",
			wordsubgen.NewField("input", *lines),
			wordsubgen.NewField("lines", len(inputLines)))
	} else {
		inputLines, err = wordsubgen.ReadLinesFromFile(*file, cfg.Logger)
		if err != nil {
			logger.Error("Failed to read input file",
				wordsubgen.NewField("error", err),
				wordsubgen.NewField("file", *file))
			os.Exit(1)
		}
		logger.Info("Read input lines from file",
			wordsubgen.NewField("file", *file),
			wordsubgen.NewField("lines", len(inputLines)))
	}

	if len(inputLines) == 0 {
		logger.Error("No valid input lines found")
		os.Exit(1)
	}

	// Generate ASS content
	logger.Info("Generating ASS content...")
	content, err := wordsubgen.GenerateASS(cfg, inputLines)
	if err != nil {
		logger.Error("Failed to generate ASS content",
			wordsubgen.NewField("error", err))
		os.Exit(1)
	}

	// Write to file
	logger.Info("Writing ASS file...",
		wordsubgen.NewField("output", *output))
	err = wordsubgen.WriteASS(*output, content, cfg.Logger)
	if err != nil {
		logger.Error("Failed to write ASS file",
			wordsubgen.NewField("error", err),
			wordsubgen.NewField("file", *output))
		os.Exit(1)
	}

	logger.Info("Successfully generated ASS file",
		wordsubgen.NewField("file", *output))
}

func showHelp() {
	// Get default configuration to show actual defaults
	defaultCfg := wordsubgen.DefaultConfig()

	fmt.Println("wordsubgen - Generate ASS subtitle files with word-by-word fade effects")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  wordsubgen [options]")
	fmt.Println()
	fmt.Println("Input (choose one):")
	fmt.Println("  --lines string    Input lines separated by | (e.g., 'Hello world|Second line')")
	fmt.Println("  --file string     Input file with lines (one per line)")
	fmt.Println()
	fmt.Println("Output:")
	fmt.Println("  --out string      Output ASS file path (default: output.ass)")
	fmt.Println()
	fmt.Printf("Video settings:\n")
	fmt.Printf("  --width int       Video width (default: %d)\n", defaultCfg.Width)
	fmt.Printf("  --height int      Video height (default: %d)\n", defaultCfg.Height)
	fmt.Println()
	fmt.Printf("Font settings:\n")
	fmt.Printf("  --font string     Font name (default: %s)\n", defaultCfg.FontName)
	fmt.Printf("  --fontsize int    Font size (default: %d)\n", defaultCfg.FontSize)
	fmt.Println()
	fmt.Println("Colors (hex format):")
	fmt.Printf("  --color string       Primary text color (default: %s)\n", wordsubgen.ColorToHex(defaultCfg.PrimaryColor))
	fmt.Printf("  --secondary string   Secondary color for karaoke (default: %s)\n", wordsubgen.ColorToHex(defaultCfg.SecondaryColor))
	fmt.Printf("  --outline string     Outline color (default: %s)\n", wordsubgen.ColorToHex(defaultCfg.OutlineColor))
	fmt.Printf("  --background string  Background color in ASS format (default: %s)\n", defaultCfg.BackColor)
	fmt.Println()
	fmt.Println("Style settings:")
	fmt.Printf("  --bold          Bold text (default: %t)\n", defaultCfg.Bold)
	fmt.Printf("  --italic        Italic text (default: %t)\n", defaultCfg.Italic)
	fmt.Printf("  --underline     Underline text (default: %t)\n", defaultCfg.Underline)
	fmt.Printf("  --strikeout     Strikeout text (default: %t)\n", defaultCfg.StrikeOut)
	fmt.Println()
	fmt.Println("Scaling:")
	fmt.Printf("  --scalex int    Horizontal scale %% (default: %d)\n", defaultCfg.ScaleX)
	fmt.Printf("  --scaley int    Vertical scale %% (default: %d)\n", defaultCfg.ScaleY)
	fmt.Printf("  --spacing int   Character spacing (default: %d)\n", defaultCfg.Spacing)
	fmt.Printf("  --angle int     Text angle (default: %d)\n", defaultCfg.Angle)
	fmt.Println()
	fmt.Println("Border and shadow:")
	fmt.Printf("  --borderstyle int  Border style (default: %d)\n", defaultCfg.BorderStyle)
	fmt.Printf("  --outlinewidth int Outline width (default: %d)\n", defaultCfg.Outline)
	fmt.Printf("  --shadow int       Shadow width (default: %d)\n", defaultCfg.Shadow)
	fmt.Println()
	fmt.Println("Shadow effects:")
	fmt.Printf("  --shadow-enabled   Enable shadow effect (default: %t)\n", defaultCfg.ShadowEnabled)
	fmt.Printf("  --shadow-x int      Horizontal shadow offset (default: %d)\n", defaultCfg.ShadowX)
	fmt.Printf("  --shadow-y int      Vertical shadow offset (default: %d)\n", defaultCfg.ShadowY)
	fmt.Println()
	fmt.Printf("Alignment:\n")
	fmt.Printf("  --alignment int    Text alignment (1=bottom-left, 2=bottom-center, 3=bottom-right, etc.) (default: %d)\n", defaultCfg.Alignment)
	fmt.Println()
	fmt.Println("Margins:")
	fmt.Printf("  --marginl int    Left margin (default: %d)\n", defaultCfg.MarginL)
	fmt.Printf("  --marginr int     Right margin (default: %d)\n", defaultCfg.MarginR)
	fmt.Printf("  --marginv int     Vertical margin (default: %d)\n", defaultCfg.MarginV)
	fmt.Println()
	fmt.Println("Timing:")
	fmt.Printf("  --start-delay int Delay before starting subtitles in ms (default: %d)\n", defaultCfg.StartDelay)
	fmt.Printf("  --delay int      Delay between words in ms (default: %d)\n", defaultCfg.PerWordDelay)
	fmt.Printf("  --fade int       Fade duration in ms (default: %d)\n", defaultCfg.FadeDuration)
	fmt.Printf("  --hold int       Line hold duration in ms (default: %d)\n", defaultCfg.LineHold)
	fmt.Printf("  --gap int        Gap between lines in ms (default: %d)\n", defaultCfg.LineGap)
	fmt.Println()
	fmt.Println("Features:")
	fmt.Printf("  --karaoke       Enable karaoke mode (default: %t)\n", defaultCfg.Karaoke)
	fmt.Println()
	fmt.Println("Other:")
	fmt.Println("  --help          Show this help")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Generate from command line")
	fmt.Println("  wordsubgen --lines 'Hello world|This is a test' --out subtitle.ass")
	fmt.Println()
	fmt.Println("  # Generate from file")
	fmt.Println("  wordsubgen --file input.txt --out subtitle.ass")
	fmt.Println()
	fmt.Println("  # Custom styling")
	fmt.Println("  wordsubgen --lines 'Hello world' --fontsize 48 --color '#FF0000' --bold")
	fmt.Println()
	fmt.Println("  # Karaoke mode")
	fmt.Println("  wordsubgen --lines 'Hello world' --karaoke --delay 500")
}
