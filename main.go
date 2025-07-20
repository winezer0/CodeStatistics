package main

import (
	"CodeStatistics/logging"
	"CodeStatistics/statistics"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
)

// parseExtensionList parses comma-separated extension list
func parseExtensionList(extensionStr string) []string {
	if extensionStr == "" {
		return nil
	}

	extensions := strings.Split(extensionStr, ",")
	var result []string

	for _, ext := range extensions {
		ext = strings.TrimSpace(ext)
		if ext == "" {
			continue
		}

		// Ensure extension starts with dot
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}

		// Convert to lowercase
		ext = strings.ToLower(ext)
		result = append(result, ext)
	}

	return result
}

// parseDirectoryList parses comma-separated directory list
func parseDirectoryList(dirStr string) []string {
	if dirStr == "" {
		return nil
	}

	directories := strings.Split(dirStr, ",")
	var result []string

	for _, dir := range directories {
		dir = strings.TrimSpace(dir)
		if dir == "" {
			continue
		}

		// Convert to lowercase for consistency
		dir = strings.ToLower(dir)
		result = append(result, dir)
	}

	return result
}

// Options command line options
type Options struct {
	Path     string `short:"p" long:"path" description:"Code directory path to scan (default: null)"`
	Output   string `short:"o" long:"output" description:"Output CSV file path (default: null)"`
	Comments bool   `short:"c" long:"comments" description:"Enable comment line detection (default:false)"`

	// Whitelist and blacklist configuration
	WhiteAdd   string `short:"w" long:"white-add" description:" add ext list to built-in whitelist, comma separated (e.g: .ext1,.ext2)"`
	WhiteCover string `short:"W" long:"white-cover" description:"use ext list override built-in whitelist, comma separated (e.g: .ext1,.ext2)"`
	BlackAdd   string `short:"b" long:"black-add" description:"add ext list to built-in blacklist, comma separated (e.g: .ext1,.ext2)"`
	BlackCover string `short:"B" long:"black-cover" description:"use exy list override built-in blacklist, comma separated (e.g: .ext1,.ext2)"`
	OnlyWhite  bool   `short:"O" long:"only-white" description:"Only Show whitelist files, skip blacklist and unknown file analysis (default: false)"`
	BDirAdd    string `short:"d" long:"bdir-add" description:"add dir list to built-in blacklist dirs, comma separated (e.g: dir1,dir2)"`
	BDirCover  string `short:"D" long:"bdir-cover" description:"use dir list override built-in blacklist dirs, comma separated (e.g: dir1,dir2)"`

	// Information display
	ShowDefaults bool `short:"s" long:"show-builtin" description:"Show built-in default white ext/black ext/black dirs"`

	// Log configuration
	LogFile       string `long:"lf" description:"Log file path (default: null)" `
	LogLevel      string `long:"ll" description:"Log level (debug/info/warn/error)" default:"info"`
	ConsoleFormat string `long:"cf" description:"Console log format (T L C M F combination or off|null to disable)" default:"M"`
}

func main() {
	var opts Options

	parser := flags.NewParser(&opts, flags.Default)
	parser.Usage = "[OPTIONS]"

	// Custom help information
	parser.LongDescription = `Code line count tool`

	_, err := parser.Parse()
	if err != nil {
		var flagsErr *flags.Error
		if errors.As(err, &flagsErr) && errors.Is(flagsErr.Type, flags.ErrHelp) {
			return
		}
		fmt.Fprintf(os.Stderr, "cmd options parsing error: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logCfg := logging.NewLogConfig(opts.LogLevel, opts.LogFile, opts.ConsoleFormat)
	if err := logging.InitLogger(logCfg); err != nil {
		// Cannot use logging here as it's not initialized yet
		fmt.Printf("Failed to initialize the logger: %v\n", err)
		os.Exit(1)
	}
	defer logging.Sync()

	// Check if user wants to show defaults
	if opts.ShowDefaults {
		statistics.ShowDefaults()
		return
	}

	if opts.Path == "" {
		logging.Fatalf("please input valid dir path current path is:[%s]", opts.Path)
	}

	// Check if path exists
	if _, err := os.Stat(opts.Path); os.IsNotExist(err) {
		logging.Fatalf("path does not exist: %s", opts.Path)
	}

	// Parse whitelist configuration
	var whitelistConfig *statistics.WhitelistConfig
	if opts.WhiteAdd != "" || opts.WhiteCover != "" {
		whitelistConfig = &statistics.WhitelistConfig{
			Add:      parseExtensionList(opts.WhiteAdd),
			Override: parseExtensionList(opts.WhiteCover),
		}
	}

	// Parse blacklist configuration
	var blacklistConfig *statistics.BlacklistConfig
	if opts.BlackAdd != "" || opts.BlackCover != "" {
		blacklistConfig = &statistics.BlacklistConfig{
			Add:      parseExtensionList(opts.BlackAdd),
			Override: parseExtensionList(opts.BlackCover),
		}
	}

	// Parse directory blacklist configuration
	var blackDirConfig *statistics.BlackDirConfig
	if opts.BDirAdd != "" || opts.BDirCover != "" {
		blackDirConfig = &statistics.BlackDirConfig{
			Add:      parseDirectoryList(opts.BDirAdd),
			Override: parseDirectoryList(opts.BDirCover),
		}
	}

	// Create statistics analyzer
	statistical := statistics.NewCodeStatistics(opts.Path, opts.Comments, opts.OnlyWhite, whitelistConfig, blacklistConfig, blackDirConfig)
	logging.Infof("Start scanning the path: %s", opts.Path)

	// Scan directory
	if err := statistical.ScanDirFiles(); err != nil {
		logging.Errorf("scanning the dir %s occur error: %v", opts.Path, err)
	}
	// Calculate file type ratios
	statistical.CalculateFileRatios()
	// Print statistics summary
	statistical.PrintSummary()

	// Generate CSV report
	if len(opts.Output) > 0 {
		if err := statistical.GenerateCSVReport(opts.Output); err != nil {
			logging.Errorf("generate csv report occur error: %v", err)
		} else {
			logging.Infof("generate csv report success: %s", opts.Output)
		}
	}
}
