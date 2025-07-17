package ezgen

import (
	"slices"
	"strings"
)

type TableFilter struct {
	SkipTables   []string
	SkipPrefixes []string
	SkipSuffixes []string
}

func (f *TableFilter) ShouldSkip(tableName string) bool {
	// 检查完整表名
	if slices.Contains(f.SkipTables, tableName) {
		return true
	}

	// 检查前缀
	for _, prefix := range f.SkipPrefixes {
		if strings.HasPrefix(tableName, prefix) {
			return true
		}
	}

	// 检查后缀
	for _, suffix := range f.SkipSuffixes {
		if strings.HasSuffix(tableName, suffix) {
			return true
		}
	}

	return false
}
