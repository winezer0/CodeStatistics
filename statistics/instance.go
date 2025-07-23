package statistics

import (
	"CodeStatistics/pkg/cmdutils"
	"time"
)

// NewCodeStatistics creates a new code statistics analyzer
func NewCodeStatistics(rootPath string, enableComments bool, whitelistOnly bool, whitelistConfig *WhitelistConfig, blacklistConfig *BlacklistConfig, blackDirConfig *BlackDirConfig) *CodeStatistics {
	// 构建白名单映射
	var whitelist map[string]bool
	if whitelistConfig != nil {
		whitelist = cmdutils.BuildAddOrCoverMap(DefaultWhitelist, whitelistConfig.Add, whitelistConfig.Override)
	} else {
		whitelist = cmdutils.BuildAddOrCoverMap(DefaultWhitelist, nil, nil)
	}

	// 构建黑名单映射
	var blacklist map[string]bool
	if blacklistConfig != nil {
		blacklist = cmdutils.BuildAddOrCoverMap(DefaultBlacklist, blacklistConfig.Add, blacklistConfig.Override)
	} else {
		blacklist = cmdutils.BuildAddOrCoverMap(DefaultBlacklist, nil, nil)
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
