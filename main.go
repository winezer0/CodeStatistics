package main

import (
	"CodeStatistics/core"
	"errors"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

// Options 命令行选项
type Options struct {
	Path           string `short:"p" long:"path" description:"要扫描的代码目录路径" default:"."`
	Output         string `short:"o" long:"output" description:"输出CSV文件路径" default:"code_stats.csv"`
	EnableComments bool   `short:"c" long:"comments" description:"启用注释行判断功能"`
	Help           bool   `short:"h" long:"help" description:"显示帮助信息"`
}

func main() {
	var opts Options

	parser := flags.NewParser(&opts, flags.Default)
	parser.Usage = "[OPTIONS]"

	// 自定义帮助信息
	parser.LongDescription = `代码行数统计工具
  ./CodeStatistics -p ./src -o stats.csv
  ./CodeStatistics --path /home/user/project --output report.csv`

	_, err := parser.Parse()
	if err != nil {
		var flagsErr *flags.Error
		if errors.As(err, &flagsErr) && errors.Is(flagsErr.Type, flags.ErrHelp) {
			return
		}
		fmt.Fprintf(os.Stderr, "参数解析错误: %v\n", err)
		os.Exit(1)
	}

	if err := runStatistics(&opts); err != nil {
		fmt.Fprintf(os.Stderr, "执行统计时出错: %v\n", err)
		os.Exit(1)
	}
}

// runStatistics 执行统计分析
func runStatistics(opts *Options) error {
	// 检查路径是否存在
	if _, err := os.Stat(opts.Path); os.IsNotExist(err) {
		return fmt.Errorf("路径 '%s' 不存在", opts.Path)
	}

	// 创建统计器
	stats := core.NewCodeStatistics(opts.Path, opts.EnableComments)

	fmt.Printf("开始扫描目录: %s\n", opts.Path)

	// 扫描目录
	if err := stats.ScanDirectory(); err != nil {
		return fmt.Errorf("扫描目录时出错: %v", err)
	}

	// 计算文件类型占比
	stats.CalculateFileRatios()

	// 打印统计摘要
	stats.PrintSummary()

	// 生成CSV报告
	if err := stats.GenerateCSVReport(opts.Output); err != nil {
		return fmt.Errorf("生成CSV报告时出错: %v", err)
	}

	fmt.Printf("\nCSV报告已生成: %s\n", opts.Output)
	return nil
}
