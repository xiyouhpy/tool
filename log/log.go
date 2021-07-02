package log

import (
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var (
	// logPath 日志保存目录文件
	logPath = "./log/access.log"

	// rotationTime 设置日志切割时间间隔，单位：小时（每隔1小时切割）
	rotationTime = 1

	// defaultRotationCount 设置日志保留个数（保留3天）
	rotationCount = 72
)

// ConfInfo 日志配置结构
type ConfInfo struct {
	LogPath       string
	RotationTime  int
	RotationCount int
}

// NewLogger 日志打印设置
func NewLogger(c ConfInfo) {
	if c.LogPath != "" {
		logPath = c.LogPath
	}
	if c.RotationTime > 0 {
		rotationTime = c.RotationTime
	}
	if c.RotationCount > 0 {
		rotationCount = c.RotationCount
	}

	writer, _ := rotatelogs.New(
		logPath+".%Y%m%d%H",
		rotatelogs.WithLinkName(logPath),
		rotatelogs.WithRotationCount(uint(rotationCount)),
		rotatelogs.WithRotationTime(time.Duration(rotationTime)*time.Hour),
	)
	logrus.SetOutput(writer)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
}
