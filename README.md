# golang 工具库

## file 目录
### file.go —— 文件操作相关
```go
// func IsFileExists 判断文件是否存在 demo
strFileName := "hpy.go"
if IsFileExists(strFileName) {
	fmt.Printf("file[%s] is exist", strFileName)
} else {
	fmt.Printf("file[%s] is not exist", strFileName)
}
```

```go
// func IsDirExists 判断目录是否存在 demo
strDirName := "/Users/hanpeiyan/go"
if IsDirExists(strDirName) {
    fmt.Printf("dir[%s] is exist", strDirName)
} else {
    fmt.Printf("dir[%s] is not exist", strDirName)
}
```

### download.go —— 数据下载相关
```go
// func DownloadUrl 根据 url 下载相关内容文件 demo
strUrl := "http://nginx.org/download/nginx-1.18.0.tar.gz"
dstFile := "/Users/hanpeiyan/data/nginx.tar.gz"
fileSize, err := DownloadUrl(strUrl, dstFile)
if err != nil {
	return err
}
fmt.Printf("download file size:%d", fileSize)
```

## request 目录
### ratelimit.go —— 请求QPS控制相关
该库只是对 golang.org/x/time/rate 库的相关方法做了简单封装（原库基于令牌桶实现）
```go
// func NewLimiter 初始化请求，参数表示每秒钟限制的请求个数，也就是QPS；下面示例的QPS限制为 10
rate, err := request.NewLimiter(10)
if err != nil {
    return err
}

// func IsPass 参数 50 表示命中 QPS 控制后等待 50ms 超时时间
if rate.IsPass(50) {
    fmt.Println("qps pass")
} else {
    fmt.Println("qps stop")
}
```