package core

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// GenerateCSVReport 生成CSV报告
func (cs *CodeStatistics) GenerateCSVReport(outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入CSV头部
	var headers []string
	if cs.EnableComments {
		headers = []string{
			"文件类型", "文件数量", "文件占比(%)", "总行数", "代码行数", "空行数", "注释行数", "代码占比(%)",
		}
	} else {
		headers = []string{
			"文件类型", "文件数量", "文件占比(%)", "总行数", "代码行数", "空行数", "代码占比(%)",
		}
	}
	if err := writer.Write(headers); err != nil {
		return err
	}

	summary := cs.calculateSummary()

	// 写入每种文件类型的统计
	for _, ext := range summary.Extensions {
		var record []string

		// 检查是否为代码文件
		if stats, exists := cs.Stats[ext]; exists {
			codeRatio := calculateCodeRatio(stats.CodeLines, stats.TotalLines)
			if cs.EnableComments {
				record = []string{
					stats.Extension,
					strconv.Itoa(stats.FileCount),
					fmt.Sprintf("%.2f", stats.FileRatio),
					strconv.Itoa(stats.TotalLines),
					strconv.Itoa(stats.CodeLines),
					strconv.Itoa(stats.BlankLines),
					strconv.Itoa(stats.CommentLines),
					fmt.Sprintf("%.2f", codeRatio),
				}
			} else {
				record = []string{
					stats.Extension,
					strconv.Itoa(stats.FileCount),
					fmt.Sprintf("%.2f", stats.FileRatio),
					strconv.Itoa(stats.TotalLines),
					strconv.Itoa(stats.CodeLines),
					strconv.Itoa(stats.BlankLines),
					fmt.Sprintf("%.2f", codeRatio),
				}
			}
		} else if stats, exists := cs.BlacklistStats[ext]; exists {
			// 黑名单文件：只记录文件数量和占比
			if cs.EnableComments {
				record = []string{
					stats.Extension,
					strconv.Itoa(stats.FileCount),
					fmt.Sprintf("%.2f", stats.FileRatio),
					"-", "-", "-", "-", "-",
				}
			} else {
				record = []string{
					stats.Extension,
					strconv.Itoa(stats.FileCount),
					fmt.Sprintf("%.2f", stats.FileRatio),
					"-", "-", "-", "-",
				}
			}
		}

		if err := writer.Write(record); err != nil {
			return err
		}
	}

	// 写入总计行
	return cs.writeTotalRecord(writer, summary)
}

// writeTotalRecord 写入总计记录
func (cs *CodeStatistics) writeTotalRecord(writer *csv.Writer, summary *SummaryData) error {
	totalCodeRatio := calculateCodeRatio(summary.TotalCodeLines, summary.TotalLines)

	var totalRecord []string
	if cs.EnableComments {
		totalRecord = []string{
			"总计",
			strconv.Itoa(summary.TotalFiles),
			"100.00",
			strconv.Itoa(summary.TotalLines),
			strconv.Itoa(summary.TotalCodeLines),
			strconv.Itoa(summary.TotalBlankLines),
			strconv.Itoa(summary.TotalCommentLines),
			fmt.Sprintf("%.2f", totalCodeRatio),
		}
	} else {
		totalRecord = []string{
			"总计",
			strconv.Itoa(summary.TotalFiles),
			"100.00",
			strconv.Itoa(summary.TotalLines),
			strconv.Itoa(summary.TotalCodeLines),
			strconv.Itoa(summary.TotalBlankLines),
			fmt.Sprintf("%.2f", totalCodeRatio),
		}
	}

	return writer.Write(totalRecord)
}

// calculateSummary 计算统计摘要
func (cs *CodeStatistics) calculateSummary() *SummaryData {
	summary := &SummaryData{
		TotalFiles: cs.TotalFiles,
	}

	// 获取所有扩展名（代码文件和黑名单文件）
	allExtensions := make(map[string]bool)
	for ext := range cs.Stats {
		allExtensions[ext] = true
		summary.Extensions = append(summary.Extensions, ext)
	}
	for ext := range cs.BlacklistStats {
		if !allExtensions[ext] {
			summary.Extensions = append(summary.Extensions, ext)
		}
	}

	// 按文件数量从少到多排序（ASC）
	sort.Slice(summary.Extensions, func(i, j int) bool {
		countI := cs.getFileCount(summary.Extensions[i])
		countJ := cs.getFileCount(summary.Extensions[j])
		return countI < countJ
	})

	// 计算总计（只统计代码文件的行数信息）
	for _, ext := range summary.Extensions {
		if stats, exists := cs.Stats[ext]; exists {
			summary.TotalLines += stats.TotalLines
			summary.TotalCodeLines += stats.CodeLines
			summary.TotalBlankLines += stats.BlankLines
			summary.TotalCommentLines += stats.CommentLines
		}
	}

	return summary
}

// getFileCount 获取文件数量（代码文件或黑名单文件）
func (cs *CodeStatistics) getFileCount(ext string) int {
	if stats, exists := cs.Stats[ext]; exists {
		return stats.FileCount
	}
	if stats, exists := cs.BlacklistStats[ext]; exists {
		return stats.FileCount
	}
	return 0
}

// calculateCodeRatio 计算代码占比
func calculateCodeRatio(codeLines, totalLines int) float64 {
	if totalLines == 0 {
		return 0.0
	}
	return float64(codeLines) / float64(totalLines) * 100
}

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
	fmt.Printf("\n=== 代码统计报告 ===\n")
	fmt.Printf("扫描路径: %s\n", cs.RootPath)
	fmt.Printf("扫描时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("运行时间: %v\n", duration)
	fmt.Printf("总文件数: %d\n\n", cs.TotalFiles)
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
		if stats, exists := cs.Stats[ext]; exists {
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
			"总计", summary.TotalFiles, 100.0, summary.TotalLines,
			summary.TotalCodeLines, summary.TotalBlankLines, summary.TotalCommentLines, totalCodeRatio)
	} else {
		fmt.Println(strings.Repeat("-", 85))
		fmt.Printf("%-15s %8d %9.2f%% %10d %10d %10d %9.2f%%\n",
			"总计", summary.TotalFiles, 100.0, summary.TotalLines,
			summary.TotalCodeLines, summary.TotalBlankLines, totalCodeRatio)
	}
}
