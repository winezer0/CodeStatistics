package statistics

import (
	"CodeStatistics/logging"
	"fmt"
	"strings"
	"time"
)

// PrintSummary 打印统计摘要
func (cs *CodeStatistics) PrintSummary() {
	cs.printHeader()
	cs.printTableHeader()

	summary := cs.calculateSummary()
	cs.printFileStats(summary)
	cs.printTotalStats(summary)
}

// printHeader 打印报告头部
func (cs *CodeStatistics) printHeader() {
	duration := time.Since(cs.StartTime)
	fmt.Printf("\n=== Code statistics report ===\n")
	fmt.Printf("scan path: %s\n", cs.RootPath)
	fmt.Printf("scan time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("run time: %v\n", duration)
	fmt.Printf("total files: %d\n\n", cs.TotalFiles)
}

// printTableHeader 打印表格头部
func (cs *CodeStatistics) printTableHeader() {
	if cs.EnableComments {
		fmt.Printf("%-15s %8s %10s %10s %10s %10s %10s %10s\n",
			"文件类型", "文件数", "文件占比", "总行数", "代码行数", "空行数", "注释行数", "代码占比")
		fmt.Println(strings.Repeat("-", 95))
	} else {
		fmt.Printf("%-15s %8s %10s %10s %10s %10s %10s\n",
			"文件类型", "文件数", "文件占比", "总行数", "代码行数", "空行数", "代码占比")
		fmt.Println(strings.Repeat("-", 85))
	}
}

// printFileStats 打印文件统计信息
func (cs *CodeStatistics) printFileStats(summary *SummaryData) {
	for _, ext := range summary.Extensions {
		// 检查是否为代码文件
		if stats, exists := cs.WhitelistStats[ext]; exists {
			codeRatio := calculateCodeRatio(stats.CodeLines, stats.TotalLines)
			if cs.EnableComments {
				fmt.Printf("%-15s %8d %9.2f%% %10d %10d %10d %10d %9.2f%%\n",
					stats.Extension, stats.FileCount, stats.FileRatio, stats.TotalLines,
					stats.CodeLines, stats.BlankLines, stats.CommentLines, codeRatio)
			} else {
				fmt.Printf("%-15s %8d %9.2f%% %10d %10d %10d %9.2f%%\n",
					stats.Extension, stats.FileCount, stats.FileRatio, stats.TotalLines,
					stats.CodeLines, stats.BlankLines, codeRatio)
			}
		} else if stats, exists := cs.BlacklistStats[ext]; exists {
			// 黑名单文件：只显示文件数量和占比，其他信息显示为 "-"
			if cs.EnableComments {
				fmt.Printf("%-15s %8d %9.2f%% %10s %10s %10s %10s %9s\n",
					stats.Extension, stats.FileCount, stats.FileRatio,
					"-", "-", "-", "-", "-")
			} else {
				fmt.Printf("%-15s %8d %9.2f%% %10s %10s %10s %9s\n",
					stats.Extension, stats.FileCount, stats.FileRatio,
					"-", "-", "-", "-")
			}
		}
	}
}

// printTotalStats 打印总计统计信息
func (cs *CodeStatistics) printTotalStats(summary *SummaryData) {
	totalCodeRatio := calculateCodeRatio(summary.TotalCodeLines, summary.TotalLines)

	if cs.EnableComments {
		fmt.Println(strings.Repeat("-", 95))
		fmt.Printf("%-15s %8d %9.2f%% %10d %10d %10d %10d %9.2f%%\n",
			"Total", summary.TotalFiles, 100.0, summary.TotalLines,
			summary.TotalCodeLines, summary.TotalBlankLines, summary.TotalCommentLines, totalCodeRatio)
	} else {
		fmt.Println(strings.Repeat("-", 85))
		fmt.Printf("%-15s %8d %9.2f%% %10d %10d %10d %9.2f%%\n",
			"Total", summary.TotalFiles, 100.0, summary.TotalLines,
			summary.TotalCodeLines, summary.TotalBlankLines, totalCodeRatio)
	}
}

// ShowDefaults displays built-in default configurations
func ShowDefaults() {
	logging.Infof("Default White EXT:\n  Total: %d \n  Content: %+v\n", len(DefaultWhitelist), DefaultWhitelist)
	logging.Infof("Default Black EXT:\n  Total: %d \n  Content: %+v\n", len(DefaultBlacklist), DefaultBlacklist)
	logging.Infof("Default Black DIR:\n  Total: %d \n  Content: %+v\n", len(DefaultBlackDirs), DefaultBlackDirs)
}
