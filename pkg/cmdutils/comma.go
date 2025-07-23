package cmdutils

import "strings"

// ParseCommaStrToList 解析逗号分隔字符串为列表
func ParseCommaStrToList(str string) []string {
	if str == "" {
		return nil
	}

	directories := strings.Split(str, ",")
	var result []string

	for _, dir := range directories {
		dir = strings.TrimSpace(dir)
		if dir == "" {
			continue
		}

		// Convert to lowercase for consistency
		dir = strings.ToLower(dir)
		result = append(result, dir)
	}

	return result
}

// EnsurePrefix 如果数组中的某元素没有前缀，则为其添加前缀
func EnsurePrefix(extList []string, prefix string) []string {
	if len(extList) > 0 {
		for i, ext := range extList {
			if !strings.HasPrefix(ext, prefix) {
				extList[i] = prefix + ext
			}
		}
	}
	return extList
}

// ParseExtensionList 解析逗号分隔的扩展名列表
func ParseExtensionList(extensionStr string) []string {
	result := ParseCommaStrToList(extensionStr)
	result = EnsurePrefix(result, ".")
	return result
}
