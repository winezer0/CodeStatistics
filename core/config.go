package core

import (
	"time"
)

// 默认黑名单 - 排除非代码文件
var defaultBlacklist = []string{
	// 图片文件
	".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".ico", ".webp",
	// 音视频文件
	".mp3", ".mp4", ".avi", ".mov", ".wmv", ".flv", ".wav", ".ogg",
	// 压缩文件
	".zip", ".rar", ".7z", ".tar", ".gz", ".bz2",
	// 文档文件
	".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
	// 字体文件
	".ttf", ".otf", ".woff", ".woff2", ".eot",
	// 其他二进制文件
	".exe", ".dll", ".so", ".dylib", ".bin", ".dat",
	// 临时文件
	".tmp", ".temp", ".cache", ".log",
	// 版本控制和配置目录相关
	".git", ".svn", ".hg",
	// 包管理器文件
	".lock",
}

// 需要跳过的目录
var skipDirectories = []string{
	".git", "node_modules", "vendor", ".svn", ".hg",
	"target", "build", "dist", ".idea", ".vscode",
}

// NewCodeStatistics 创建新的代码统计器
func NewCodeStatistics(rootPath string, enableComments bool) *CodeStatistics {
	blacklist := make(map[string]bool)
	for _, ext := range defaultBlacklist {
		blacklist[ext] = true
	}

	return &CodeStatistics{
		RootPath:       rootPath,
		Stats:          make(map[string]*FileStats),
		BlacklistStats: make(map[string]*FileStats),
		Blacklist:      blacklist,
		EnableComments: enableComments,
		StartTime:      time.Now(),
	}
}
