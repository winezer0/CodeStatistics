package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

func main() {
	var opts Options

	parser := flags.NewParser(&opts, flags.Default)
	parser.Usage = "[OPTIONS]"

	// 自定义帮助信息
	parser.LongDescription = `代码行数统计工具

这个工具可以扫描指定目录下的所有代码文件，统计各种文件类型的行数信息，
包括总行数、代码行数、空行数、注释行数等，并生成详细的CSV报告。

示例:
  ./CodeStatistics -p ./src -o stats.csv
  ./CodeStatistics --path /home/user/project --output report.csv`

	_, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
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
	stats := NewCodeStatistics(opts.Path, opts.EnableComments)

	fmt.Printf("开始扫描目录: %s\n", opts.Path)

	// 扫描目录
	if err := stats.scanDirectory(); err != nil {
		return fmt.Errorf("扫描目录时出错: %v", err)
	}

	// 计算文件类型占比
	stats.calculateFileRatios()

	// 打印统计摘要
	stats.printSummary()

	// 生成CSV报告
	if err := stats.generateCSVReport(opts.Output); err != nil {
		return fmt.Errorf("生成CSV报告时出错: %v", err)
	}

	fmt.Printf("\nCSV报告已生成: %s\n", opts.Output)
	return nil
}
