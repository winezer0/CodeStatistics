package core

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

// SummaryData 统计摘要数据
type SummaryData struct {
	TotalFiles        int
	TotalLines        int
	TotalCodeLines    int
	TotalBlankLines   int
	TotalCommentLines int
	Extensions        []string
}
