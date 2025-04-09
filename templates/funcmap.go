package templates

import (
	"html/template"
	"strings"
)

// FuncMap 返回所有自定义模板函数
func GetFuncMap() template.FuncMap {
	return template.FuncMap{
		"add":      add,
		"subtract": subtract,
		"truncate": truncate,
		"contains": contains,
	}
}

// add 将两个数相加
func add(a, b int) int {
	return a + b
}

// subtract 将两个数相减
func subtract(a, b int) int {
	return a - b
}

// truncate 截断字符串到指定长度
func truncate(s string, l int) string {
	if len(s) <= l {
		return s
	}
	return s[:l-3] + "..."
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
