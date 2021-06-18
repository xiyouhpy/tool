package file

import "os"

// File 接口整理
type File interface {
	// IsFileExists 判断文件是否存在
	IsFileExists(fileName string) bool
}

// IsFileExists 判断文件是否存在
func IsFileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}
