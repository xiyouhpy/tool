package file

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var tmpDir = "/tmp/"

// IsFileExists 判断文件是否存在
// 		return: 存在返回: true，不存在返回: false
func IsFileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

// DownloadFile 单个文件下载
func DownloadFile(dstFileName string, strURL string) (int64, error) {
	// 1、创建临时文件
	_, fileName := filepath.Split(dstFileName)
	tmpFileName := tmpDir + fileName + ".downloading"
	file, err := os.OpenFile(tmpFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// 2、发送下载请求
	rsp, err := http.Get(strURL)
	if err != nil {
		return 0, err
	}
	defer rsp.Body.Close()

	// 3、下载数据写入文件（copy方法使用缓存写入，一次读取大致3M，能规避OOM）
	length, err := io.Copy(file, rsp.Body)
	if err != nil {
		return length, err
	}

	// 4、移动临时文件到目标文件处
	os.Rename(tmpFileName, dstFileName)

	return length, err
}
