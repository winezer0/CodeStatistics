package cmdutils

import "strings"

// ParseCommaStrToList 解析逗号分隔字符串为列表
func ParseCommaStrToList(CommaStr string, toLower bool) []string {
	var result []string
	if CommaStr != "" {
		strList := strings.Split(CommaStr, ",")
		for _, str := range strList {
			str = strings.TrimSpace(str)
			if str == "" {
				continue
			}
			// Convert to lowercase for consistency
			if toLower {
				str = strings.ToLower(str)
			}
			result = append(result, str)
		}
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
func ParseExtensionList(extensionStr string, toLower bool) []string {
	result := ParseCommaStrToList(extensionStr, toLower)
	result = EnsurePrefix(result, ".")
	return result
}
