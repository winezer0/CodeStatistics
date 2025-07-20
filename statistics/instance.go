package statistics

import (
	"time"
)

// NewCodeStatistics creates a new code statistics analyzer
func NewCodeStatistics(rootPath string, enableComments bool, whitelistOnly bool, whitelistConfig *WhitelistConfig, blacklistConfig *BlacklistConfig, blackDirConfig *BlackDirConfig) *CodeStatistics {
	// 构建白名单映射
	var whitelist map[string]bool
	if whitelistConfig != nil {
		whitelist = buildAddOrCoverMap(DefaultWhitelist, whitelistConfig.Add, whitelistConfig.Override)
	} else {
		whitelist = buildAddOrCoverMap(DefaultWhitelist, nil, nil)
	}

	// 构建黑名单映射
	var blacklist map[string]bool
	if blacklistConfig != nil {
		blacklist = buildAddOrCoverMap(DefaultBlacklist, blacklistConfig.Add, blacklistConfig.Override)
	} else {
		blacklist = buildAddOrCoverMap(DefaultBlacklist, nil, nil)
	}

	// 构建目录黑名单
	var skipBlackDirs []string
	if blackDirConfig != nil && len(blackDirConfig.Override) > 0 {
		// 如果有覆盖配置，使用覆盖配置
		skipBlackDirs = make([]string, len(blackDirConfig.Override))
		copy(skipBlackDirs, blackDirConfig.Override)
	} else {
		// 使用默认配置
		skipBlackDirs = make([]string, len(DefaultBlackDirs))
		copy(skipBlackDirs, DefaultBlackDirs)

		// 添加额外的目录黑名单项
		if blackDirConfig != nil && len(blackDirConfig.Add) > 0 {
			skipBlackDirs = append(skipBlackDirs, blackDirConfig.Add...)
		}
	}

	return &CodeStatistics{
		RootPath:        rootPath,
		WhitelistStats:  make(map[string]*FileStats),
		BlacklistStats:  make(map[string]*FileStats),
		Whitelist:       whitelist,
		Blacklist:       blacklist,
		SkipDirectories: skipBlackDirs,
		EnableComments:  enableComments,
		OnlyWhite:       whitelistOnly,
		StartTime:       time.Now(),
	}
}

// 构建映射表（用于白名单或黑名单）
func buildAddOrCoverMap(defaultList, addList, overrideList []string) map[string]bool {
	result := make(map[string]bool)

	// 如果有覆盖配置，优先使用覆盖配置
	if len(overrideList) > 0 {
		for _, ext := range overrideList {
			result[ext] = true
		}
	} else {
		// 否则使用默认 + 追加
		for _, ext := range defaultList {
			result[ext] = true
		}
		for _, ext := range addList {
			result[ext] = true
		}
	}

	return result
}
