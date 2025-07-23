package cmdutils

// BuildAddOrCoverMap 构建映射表（用于白名单或黑名单）
func BuildAddOrCoverMap(defaultList, addList, coverList []string) map[string]bool {
	result := make(map[string]bool)

	// 如果有覆盖配置，优先使用覆盖配置
	if len(coverList) > 0 {
		for _, ext := range coverList {
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
