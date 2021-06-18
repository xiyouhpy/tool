package file

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const tmpDir = "/tmp/"

// Download 接口整理
type Download interface {
	// DownloadUrl 下载网络文件
	DownloadUrl(strURL string, dstFile string) (int64, error)
}

// isDownloadFile 判断下载的文件是否存在
func isDownloadFile(fileName string, fileSize int64) bool {
	info, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	if info.Size() != fileSize {
		os.Remove(fileName)
		return false
	}

	return true
}

// DownloadUrl 单个文件下载
func DownloadUrl(strURL string, dstFile string) (int64, error) {
	// 1、发送下载请求，获取下载对象
	client := new(http.Client)
	client.Timeout = time.Second * 600
	rsp, err := client.Get(strURL)
	if err != nil {
		return 0, err
	}
	defer rsp.Body.Close()

	// 2、判断是否已下载，未下载则下载到临时文件
	tmpFile := tmpDir + filepath.Base(dstFile) + ".downloading"
	fileSize, _ := strconv.ParseInt(rsp.Header.Get("Content-Length"), 10, 32)
	if !isDownloadFile(tmpFile, fileSize) {
		file, err := os.OpenFile(tmpFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			return 0, err
		}
		defer file.Close()

		// copy方法使用缓存写入，一次读取大致3M，能规避OOM
		if _, err := io.Copy(file, rsp.Body); err != nil {
			return 0, err
		}
	}

	// 3、移动临时文件到目标文件处
	os.Rename(tmpFile, dstFile)

	return fileSize, nil
}
