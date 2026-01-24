package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kperreau/wordsubgen"
)

// cliFlags holds flag values after parsing.
type cliFlags struct {
	lines, file, jsonArg, jsonFile, output       string
	width, height, fontSize                      int
	fontName, primaryColor, secondaryColor       string
	outlineColor, backColor                      string
	bold, italic, underline, strikeout           bool
	scaleX, scaleY, spacing, angle               int
	borderStyle, outlineWidth, shadow            int
	shadowEnabled                                bool
	shadowX, shadowY, alignment                  int
	marginL, marginR, marginV                    int
	startDelay, perWordDelay, fadeDuration       int
	lineHold, lineGap                            int
	karaoke, help                                bool
}

func parseFlags() *cliFlags {
	def := wordsubgen.DefaultConfig()
	var (
		lines    = flag.String("lines", "", "Input lines separated by | (e.g., 'Hello world|Second line')")
		file     = flag.String("file", "", "Input file with lines (one per line)")
		jsonArg  = flag.String("json", "", "Structured JSON (words with start/end in seconds). Incompatible with --lines/--file")
		jsonFile = flag.String("json-file", "", "Path to structured JSON file. Incompatible with --lines/--file")
		output   = flag.String("out", "output.ass", "Output ASS file path")
		width    = flag.Int("width", def.Width, "Video width")
		height   = flag.Int("height", def.Height, "Video height")
		fontName = flag.String("font", def.FontName, "Font name")
		fontSize = flag.Int("fontsize", def.FontSize, "Font size")
		primary  = flag.String("color", wordsubgen.ColorToHex(def.PrimaryColor), "Primary text color (hex)")
		second   = flag.String("secondary", wordsubgen.ColorToHex(def.SecondaryColor), "Secondary color for karaoke (hex)")
		outline  = flag.String("outline", wordsubgen.ColorToHex(def.OutlineColor), "Outline color (hex)")
		back     = flag.String("background", def.BackColor, "Background color (ASS format)")
		bold     = flag.Bool("bold", def.Bold, "Bold text")
		italic   = flag.Bool("italic", def.Italic, "Italic text")
		ul       = flag.Bool("underline", def.Underline, "Underline text")
		so       = flag.Bool("strikeout", def.StrikeOut, "Strikeout text")
		scaleX   = flag.Int("scalex", def.ScaleX, "Horizontal scale (%)")
		scaleY   = flag.Int("scaley", def.ScaleY, "Vertical scale (%)")
		spacing  = flag.Int("spacing", def.Spacing, "Character spacing")
		angle    = flag.Int("angle", def.Angle, "Text angle")
		bs       = flag.Int("borderstyle", def.BorderStyle, "Border style")
		ow       = flag.Int("outlinewidth", def.Outline, "Outline width")
		sh       = flag.Int("shadow", def.Shadow, "Shadow width")
		shEn     = flag.Bool("shadow-enabled", def.ShadowEnabled, "Enable shadow effect")
		shX      = flag.Int("shadow-x", def.ShadowX, "Horizontal shadow offset")
		shY      = flag.Int("shadow-y", def.ShadowY, "Vertical shadow offset")
		align    = flag.Int("alignment", def.Alignment, "Text alignment (1=bottom-left, 2=bottom-center, 3=bottom-right, etc.)")
		marL     = flag.Int("marginl", def.MarginL, "Left margin")
		marR     = flag.Int("marginr", def.MarginR, "Right margin")
		marV     = flag.Int("marginv", def.MarginV, "Vertical margin")
		startD   = flag.Int("start-delay", def.StartDelay, "Delay before starting subtitles (ms)")
		perWord  = flag.Int("delay", def.PerWordDelay, "Delay between words (ms)")
		fade     = flag.Int("fade", def.FadeDuration, "Fade duration (ms)")
		hold     = flag.Int("hold", def.LineHold, "Line hold duration (ms)")
		gap      = flag.Int("gap", def.LineGap, "Gap between lines (ms)")
		karaoke  = flag.Bool("karaoke", def.Karaoke, "Enable karaoke mode")
		help     = flag.Bool("help", false, "Show help")
	)
	flag.Parse()
	return &cliFlags{
		*lines, *file, *jsonArg, *jsonFile, *output,
		*width, *height, *fontSize, *fontName, *primary, *second, *outline, *back,
		*bold, *italic, *ul, *so, *scaleX, *scaleY, *spacing, *angle,
		*bs, *ow, *sh, *shEn, *shX, *shY, *align, *marL, *marR, *marV,
		*startD, *perWord, *fade, *hold, *gap, *karaoke, *help,
	}
}

func validateInput(f *cliFlags) error {
	textMode := f.lines != "" || f.file != ""
	jsonMode := f.jsonArg != "" || f.jsonFile != ""
	switch {
	case !textMode && !jsonMode:
		return fmt.Errorf("one of --lines, --file, --json, or --json-file must be specified")
	case textMode && jsonMode:
		return fmt.Errorf("cannot mix text input (--lines/--file) with structured JSON (--json/--json-file)")
	case f.lines != "" && f.file != "":
		return fmt.Errorf("cannot specify both --lines and --file")
	case f.jsonArg != "" && f.jsonFile != "":
		return fmt.Errorf("cannot specify both --json and --json-file")
	}
	return nil
}

func buildConfig(logger wordsubgen.Logger, f *cliFlags) *wordsubgen.Config {
	cfg := wordsubgen.DefaultConfig()
	cfg.Logger = logger
	cfg.Width, cfg.Height = f.width, f.height
	cfg.FontName, cfg.FontSize = f.fontName, f.fontSize
	cfg.PrimaryColor = wordsubgen.ColorToASS(f.primaryColor)
	cfg.SecondaryColor = wordsubgen.ColorToASS(f.secondaryColor)
	cfg.OutlineColor = wordsubgen.ColorToASS(f.outlineColor)
	cfg.BackColor = f.backColor
	cfg.Bold, cfg.Italic, cfg.Underline, cfg.StrikeOut = f.bold, f.italic, f.underline, f.strikeout
	cfg.ScaleX, cfg.ScaleY, cfg.Spacing, cfg.Angle = f.scaleX, f.scaleY, f.spacing, f.angle
	cfg.BorderStyle, cfg.Outline, cfg.Shadow = f.borderStyle, f.outlineWidth, f.shadow
	cfg.ShadowEnabled, cfg.ShadowX, cfg.ShadowY = f.shadowEnabled, f.shadowX, f.shadowY
	cfg.Alignment, cfg.MarginL, cfg.MarginR, cfg.MarginV = f.alignment, f.marginL, f.marginR, f.marginV
	cfg.StartDelay, cfg.PerWordDelay = f.startDelay, f.perWordDelay
	cfg.FadeDuration, cfg.LineHold, cfg.LineGap = f.fadeDuration, f.lineHold, f.lineGap
	cfg.Karaoke = f.karaoke
	return cfg
}

func runGenerate(cfg *wordsubgen.Config, f *cliFlags) (string, error) {
	jsonMode := f.jsonArg != "" || f.jsonFile != ""

	if jsonMode {
		var jsonBytes []byte
		if f.jsonArg != "" {
			jsonBytes = []byte(f.jsonArg)
			cfg.Logger.Info("Using structured JSON from --json", wordsubgen.NewField("length", len(jsonBytes)))
		} else {
			var err error
			jsonBytes, err = os.ReadFile(f.jsonFile)
			if err != nil {
				return "", fmt.Errorf("read JSON file %s: %w", f.jsonFile, err)
			}
			cfg.Logger.Info("Read structured JSON from file", wordsubgen.NewField("file", f.jsonFile), wordsubgen.NewField("size", len(jsonBytes)))
		}
		phrases, err := wordsubgen.ParseStructuredJSON(jsonBytes)
		if err != nil {
			return "", fmt.Errorf("parse structured JSON: %w", err)
		}
		cfg.Logger.Info("Generating ASS from structured JSON...", wordsubgen.NewField("phrases", len(phrases)))
		return wordsubgen.GenerateASSFromStructured(cfg, phrases)
	}

	var lines []string
	if f.lines != "" {
		lines = wordsubgen.ParseLinesFromString(f.lines)
		cfg.Logger.Info("Parsed input lines", wordsubgen.NewField("input", f.lines), wordsubgen.NewField("lines", len(lines)))
	} else {
		var err error
		lines, err = wordsubgen.ReadLinesFromFile(f.file, cfg.Logger)
		if err != nil {
			return "", fmt.Errorf("read input file: %w", err)
		}
		cfg.Logger.Info("Read input lines from file", wordsubgen.NewField("file", f.file), wordsubgen.NewField("lines", len(lines)))
	}
	if len(lines) == 0 {
		return "", fmt.Errorf("no valid input lines found")
	}
	cfg.Logger.Info("Generating ASS content...")
	return wordsubgen.GenerateASS(cfg, lines)
}

func main() {
	logger := wordsubgen.NewDefaultLogger()
	f := parseFlags()

	if f.help {
		showHelp()
		return
	}
	if err := validateInput(f); err != nil {
		logger.Error(err.Error())
		showHelp()
		os.Exit(1)
	}

	cfg := buildConfig(logger, f)
	content, err := runGenerate(cfg, f)
	if err != nil {
		logger.Error("Failed", wordsubgen.NewField("error", err))
		os.Exit(1)
	}

	logger.Info("Writing ASS file...", wordsubgen.NewField("output", f.output))
	if err := wordsubgen.WriteASS(f.output, content, cfg.Logger); err != nil {
		logger.Error("Failed to write ASS file", wordsubgen.NewField("error", err), wordsubgen.NewField("file", f.output))
		os.Exit(1)
	}
	logger.Info("Successfully generated ASS file", wordsubgen.NewField("file", f.output))
}

func showHelp() {
	// Get default configuration to show actual defaults
	defaultCfg := wordsubgen.DefaultConfig()

	fmt.Println("wordsubgen - Generate ASS subtitle files with word-by-word fade effects")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  wordsubgen [options]")
	fmt.Println()
	fmt.Println("Input (choose one, mutually exclusive):")
	fmt.Println("  --lines string      Input lines separated by | (e.g., 'Hello world|Second line')")
	fmt.Println("  --file string       Input file with lines (one per line)")
	fmt.Println("  --json string       Structured JSON: {\"segments\":[{\"words\":[{\"word\":\"...\",\"start\":0.1,\"end\":0.5},...]},...]}")
	fmt.Println("  --json-file string  Path to structured JSON file (same format as --json)")
	fmt.Println("                      With JSON: --delay is ignored; --start-delay shifts all timings (e.g. 1500ms = 1.5s).")
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
	fmt.Println()
	fmt.Println("  # From structured JSON (words with start/end in seconds)")
	fmt.Println("  wordsubgen --json-file words.json --start-delay 1500 --out subtitle.ass")
	fmt.Println("  wordsubgen --json '{\"segments\":[{\"words\":[{\"word\":\"Hello\",\"start\":0,\"end\":0.5},{\"word\":\"world\",\"start\":0.6,\"end\":1.2}]}]}' --out subtitle.ass")
}
