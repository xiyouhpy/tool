package file

import "os"

// File 接口整理
type File interface {
	// IsFileExists 判断文件是否存在
	IsFileExists(fileName string) bool
	// IsDirExists 判断目录是否存在
	IsDirExists(dirName string) bool
}

// IsFileExists 判断文件是否存在
func IsFileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

// IsDirExists 判断目录是否存在
func IsDirExists(dirName string) bool {
	d, err := os.Stat(dirName)
	if err != nil {
		return false
	}

	return d.IsDir()
}
