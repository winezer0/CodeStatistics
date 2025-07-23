package statistics

import (
	"CodeStatistics/pkg/logging"
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// isCodeFile 判断是否为代码文件
func isCodeFile(filename string, whitelist map[string]bool, blacklist map[string]bool, onlyWhite bool) bool {
	// 判断逻辑：
	//  1. 如果文件扩展名在白名单中，则认为是代码文件
	//  2. 如果文件扩展名在黑名单中，则认为不是代码文件
	//  3. 如果文件扩展名既不在白名单也不在黑名单中，则提示并作为代码文件处理
	ext := strings.ToLower(filepath.Ext(filename))

	// 处理无扩展名文件
	if ext == "" {
		ext = "NONE"
	}

	// 优先检查白名单
	if whitelist[ext] {
		return true
	}

	// 检查黑名单
	if blacklist[ext] {
		return false
	}

	//当启用了仅白名单模式时,将所有未知后缀文件作为黑名单处理
	if onlyWhite {
		return false
	} else {
		return true
	}
}

// ScanDirFiles 扫描目录
func (cs *CodeStatistics) ScanDirFiles() error {
	return filepath.Walk(cs.RootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logging.Errorf("access path: %s occur error: %v", path, err)
			// 对于访问权限等问题，记录错误但继续处理其他文件
			return nil
		}

		// 跳过目录
		if info.IsDir() {
			if isBlackDirs(info.Name(), cs.BlackDirs) {
				logging.Debugf("Skip this dir: %s", path)
				return filepath.SkipDir
			}
			return nil
		}

		cs.TotalFiles++

		// 处理代码文件
		if isCodeFile(info.Name(), cs.Whitelist, cs.Blacklist, cs.OnlyWhite) {
			if err := cs.analyzeCodeFile(path); err != nil {
				logging.Warnf("analyze code file: %s occur error: %v", path, err)
				// 分析单个文件失败时，记录警告但继续处理其他文件
				return nil
			}
		} else if !cs.OnlyWhite {
			// 如果不是仅显示白名单模式，则处理黑名单文件（只统计数量，不分析内容）
			if err := cs.analyzeBlackFile(info.Name()); err != nil {
				logging.Warnf("process black list file: %s occur error: %v", path, err)
				return nil
			}
		}
		// 如果是仅显示白名单模式且文件不是代码文件，则跳过处理

		return nil
	})
}

// isBlackDirs 判断是否应该跳过目录
func isBlackDirs(dirName string, skipDirectories []string) bool {
	dirName = strings.ToLower(dirName)
	for _, skipDir := range skipDirectories {
		if dirName == skipDir {
			return true
		}
	}
	return false
}

// analyze Black file 分析黑名单文件（只统计数量，不分析内容）
func (cs *CodeStatistics) analyzeBlackFile(filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		ext = "NONE"
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

// analyzeCodeFile 分析单个文件
func (cs *CodeStatistics) analyzeCodeFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		logging.Errorf("open the file: %s occur error: %v", filePath, err)
		return err
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == "" {
		ext = "NONE"
	}

	// 初始化统计信息
	if cs.WhitelistStats[ext] == nil {
		cs.WhitelistStats[ext] = &FileStats{
			Extension: ext,
		}
	}

	stats := cs.WhitelistStats[ext]
	stats.FileCount++

	if err := cs.analyzeFileContent(file, stats, ext, filePath); err != nil {
		logging.Errorf("analyze the file content : %s occur error: %v", filePath, err)
		return err
	}

	return nil
}

// analyzeFileContent 分析文件内容
func (cs *CodeStatistics) analyzeFileContent(file *os.File, stats *FileStats, ext string, filePath string) error {
	scanner := bufio.NewScanner(file)

	// 设置更大的缓冲区来处理长行，最大64MB
	const maxCapacity = 64 * 1024 * 1024 // 64MB
	buf := make([]byte, 0, 64*1024)      // 初始64KB
	scanner.Buffer(buf, maxCapacity)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
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

	if err := scanner.Err(); err != nil {
		logging.Errorf("scanning the file: %s line num: %d, occur error: %v", filePath, lineNum, err)
		return err
	}

	return nil
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
