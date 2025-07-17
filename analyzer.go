package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// isCodeFile 判断是否为代码文件
func (cs *CodeStatistics) isCodeFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return !cs.Blacklist[ext]
}

// analyzeFile 分析单个文件
func (cs *CodeStatistics) analyzeFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == "" {
		ext = "no_extension"
	}

	// 初始化统计信息
	if cs.Stats[ext] == nil {
		cs.Stats[ext] = &FileStats{
			Extension: ext,
		}
	}

	stats := cs.Stats[ext]
	stats.FileCount++

	return cs.analyzeFileContent(file, stats, ext)
}

// analyzeFileContent 分析文件内容
func (cs *CodeStatistics) analyzeFileContent(file *os.File, stats *FileStats, ext string) error {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		stats.TotalLines++

		if line == "" {
			stats.BlankLines++
		} else if cs.EnableComments && cs.isCommentLine(line, ext) {
			stats.CommentLines++
		} else {
			stats.CodeLines++
		}
	}
	return scanner.Err()
}

// isCommentLine 简单判断是否为注释行
func (cs *CodeStatistics) isCommentLine(line, ext string) bool {
	line = strings.TrimSpace(line)

	// 根据文件类型判断注释
	switch ext {
	case ".go", ".js", ".ts", ".java", ".c", ".cpp", ".h", ".hpp", ".cs", ".php", ".swift", ".kt":
		return strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") || strings.HasPrefix(line, "*")
	case ".py", ".sh", ".rb", ".pl", ".yaml", ".yml":
		return strings.HasPrefix(line, "#")
	case ".html", ".xml", ".vue":
		return strings.HasPrefix(line, "<!--")
	case ".css", ".scss", ".sass":
		return strings.HasPrefix(line, "/*") || strings.HasPrefix(line, "*")
	case ".sql":
		return strings.HasPrefix(line, "--") || strings.HasPrefix(line, "/*")
	}

	return false
}

// shouldSkipDirectory 判断是否应该跳过目录
func shouldSkipDirectory(dirName string) bool {
	dirName = strings.ToLower(dirName)
	for _, skipDir := range skipDirectories {
		if dirName == skipDir {
			return true
		}
	}
	return false
}

// scanDirectory 扫描目录
func (cs *CodeStatistics) scanDirectory() error {
	return filepath.Walk(cs.RootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			if shouldSkipDirectory(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		cs.TotalFiles++

		// 处理代码文件
		if cs.isCodeFile(info.Name()) {
			return cs.analyzeFile(path)
		} else {
			// 处理黑名单文件（只统计数量，不分析内容）
			return cs.analyzeBlacklistFile(info.Name())
		}
	})
}

// analyzeBlacklistFile 分析黑名单文件（只统计数量，不分析内容）
func (cs *CodeStatistics) analyzeBlacklistFile(filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		ext = "no_extension"
	}

	// 初始化黑名单统计信息
	if cs.BlacklistStats[ext] == nil {
		cs.BlacklistStats[ext] = &FileStats{
			Extension: ext,
		}
	}

	cs.BlacklistStats[ext].FileCount++
	return nil
}

// calculateFileRatios 计算文件类型占比
func (cs *CodeStatistics) calculateFileRatios() {
	if cs.TotalFiles == 0 {
		return
	}

	// 计算代码文件占比
	for _, stats := range cs.Stats {
		stats.FileRatio = float64(stats.FileCount) / float64(cs.TotalFiles) * 100
	}

	// 计算黑名单文件占比
	for _, stats := range cs.BlacklistStats {
		stats.FileRatio = float64(stats.FileCount) / float64(cs.TotalFiles) * 100
	}
}
