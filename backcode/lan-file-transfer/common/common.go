package common

import (
	"path/filepath"
	"strings"
)

//FileJoin 结合path，兼容windows
func FileJoin(paths ...string) string {
	path := filepath.Join(paths...)
	path = strings.ReplaceAll(path, "\\", "/") // 将反斜杠替换为正斜杠
	return path
}
