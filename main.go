package main

import (
	"CodeStatistics/pkg/cmdutils"
	"CodeStatistics/pkg/logging"
	"CodeStatistics/statistics"
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

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
	DirsAdd    string `short:"d" long:"dirs-add" description:"add dir list to built-in blacklist dirs, comma separated (e.g: dir1,dir2)"`
	DirsCover  string `short:"D" long:"dirs-cover" description:"use dir list override built-in blacklist dirs, comma separated (e.g: dir1,dir2)"`

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

	if _, err := parser.Parse(); err != nil {
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

	// Create statistics analyzer
	//statistical := InitCodeStatistics(opts.Path, opts.Comments, opts.OnlyWhite, whitelistConfig, blacklistConfig, blackDirConfig)
	statistical := InitCodeStatistics(&opts)
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

// InitCodeStatistics creates a new code statistics analyzer
func InitCodeStatistics(opts *Options) *statistics.CodeStatistics {
	// 构建白名单映射
	// Parse whitelist configuration
	whitelistConfig := &statistics.WhitelistConfig{
		Add:      cmdutils.ParseExtensionList(opts.WhiteAdd, true),
		Override: cmdutils.ParseExtensionList(opts.WhiteCover, true),
	}
	whitelist := cmdutils.BuildAddOrCoverMap(statistics.DefaultWhitelist, whitelistConfig.Add, whitelistConfig.Override)

	// 构建黑名单映射
	blacklistConfig := &statistics.BlacklistConfig{
		Add:      cmdutils.ParseExtensionList(opts.BlackAdd, true),
		Override: cmdutils.ParseExtensionList(opts.BlackCover, true),
	}
	blacklist := cmdutils.BuildAddOrCoverMap(statistics.DefaultBlacklist, blacklistConfig.Add, blacklistConfig.Override)

	// 构建目录黑名单
	// Parse directory blacklist configuration
	blackDirsConfig := &statistics.BlackDirsConfig{
		Add:      cmdutils.ParseCommaStrToList(opts.DirsAdd, true),
		Override: cmdutils.ParseCommaStrToList(opts.DirsCover, true),
	}
	blackDirs := cmdutils.BuildAddOrCoverList(statistics.DefaultBlackDirs, blackDirsConfig.Add, blackDirsConfig.Override)

	return statistics.NewCodeStatistics(opts.Path, opts.Comments, opts.OnlyWhite, whitelist, blacklist, blackDirs)
}
