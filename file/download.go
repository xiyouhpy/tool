package file

import (
	"io"
	"net/http"
	"os"
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
	// 1、创建下载文件
	file, err := os.OpenFile(dstFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// 2、发送网络请求下载文件
	rsp, err := http.Get(strURL)
	if err != nil {
		return 0, err
	}
	defer rsp.Body.Close()

	// 3、将网络请求拉取的数据写入到文件中（copy方法使用缓存写入，一次读取大致3M，能规避OOM）
	length, err := io.Copy(file, rsp.Body)
	if err != nil {
		return length, err
	}

	// 4、获取当前执行路径，拼接保存的目标文件绝对路径
	os.Rename(tmpDir+dstFileName+".download", dstFileName)

	return length, err
}
