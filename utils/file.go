package utils

import (
	"path/filepath"
)

// FileExistsFunc 检查文件是否存在
func FileExists(path string) bool {
	f, err := filepath.Glob(path)
	if err == nil && len(f) > 0 {
		return true
	}
	return false
}
