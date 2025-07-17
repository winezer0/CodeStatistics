package main

import (
	"time"
)

// FileStats 文件统计信息
type FileStats struct {
	Extension    string  // 文件扩展名
	FileCount    int     // 文件数量
	TotalLines   int     // 总行数
	CodeLines    int     // 非空行数
	BlankLines   int     // 空行数
	CommentLines int     // 注释行数
	FileRatio    float64 // 文件类型占比
}

// CodeStatistics 代码统计器
type CodeStatistics struct {
	RootPath       string
	Stats          map[string]*FileStats
	BlacklistStats map[string]*FileStats // 黑名单文件统计
	TotalFiles     int
	Blacklist      map[string]bool
	EnableComments bool // 是否启用注释行判断
	StartTime      time.Time
}

// Options 命令行选项
type Options struct {
	Path           string `short:"p" long:"path" description:"要扫描的代码目录路径" default:"."`
	Output         string `short:"o" long:"output" description:"输出CSV文件路径" default:"code_statistics.csv"`
	EnableComments bool   `short:"c" long:"comments" description:"启用注释行判断功能"`
	Help           bool   `short:"h" long:"help" description:"显示帮助信息"`
}

// SummaryData 统计摘要数据
type SummaryData struct {
	TotalFiles        int
	TotalLines        int
	TotalCodeLines    int
	TotalBlankLines   int
	TotalCommentLines int
	Extensions        []string
}
