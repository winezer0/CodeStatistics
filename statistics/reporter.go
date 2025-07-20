package statistics

import (
	"CodeStatistics/logging"
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
)

// GenerateCSVReport 生成CSV报告
func (cs *CodeStatistics) GenerateCSVReport(outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		logging.Errorf("Unable to create CSV file: %s, error: %v", outputPath, err)
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
		logging.Errorf("write CSV header error: %v", err)
		return err
	}

	summary := cs.calculateSummary()

	// 写入每种文件类型的统计
	for _, ext := range summary.Extensions {
		var record []string

		// 检查是否为代码文件
		if stats, exists := cs.WhitelistStats[ext]; exists {
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
			logging.Errorf("Failed to write CSV record: %v", err)
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
			"Total",
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
			"Total",
			strconv.Itoa(summary.TotalFiles),
			"100.00",
			strconv.Itoa(summary.TotalLines),
			strconv.Itoa(summary.TotalCodeLines),
			strconv.Itoa(summary.TotalBlankLines),
			fmt.Sprintf("%.2f", totalCodeRatio),
		}
	}

	if err := writer.Write(totalRecord); err != nil {
		logging.Errorf("write total records to CSV file error: %v", err)
		return err
	}
	return nil
}

// calculateSummary 计算统计摘要
func (cs *CodeStatistics) calculateSummary() *SummaryData {
	summary := &SummaryData{
		TotalFiles: cs.TotalFiles,
	}

	// 获取所有扩展名（代码文件和黑名单文件）
	allExtensions := make(map[string]bool)
	for ext := range cs.WhitelistStats {
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
		if stats, exists := cs.WhitelistStats[ext]; exists {
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
	if stats, exists := cs.WhitelistStats[ext]; exists {
		return stats.FileCount
	}
	if stats, exists := cs.BlacklistStats[ext]; exists {
		return stats.FileCount
	}
	return 0
}

// CalculateFileRatios 计算文件类型占比
func (cs *CodeStatistics) CalculateFileRatios() {
	if cs.TotalFiles == 0 {
		return
	}

	// 计算代码文件占比
	for _, stats := range cs.WhitelistStats {
		stats.FileRatio = float64(stats.FileCount) / float64(cs.TotalFiles) * 100
	}

	// 计算黑名单文件占比
	for _, stats := range cs.BlacklistStats {
		stats.FileRatio = float64(stats.FileCount) / float64(cs.TotalFiles) * 100
	}
}

// calculateCodeRatio 计算代码占比
func calculateCodeRatio(codeLines, totalLines int) float64 {
	if totalLines == 0 {
		return 0.0
	}
	return float64(codeLines) / float64(totalLines) * 100
}
