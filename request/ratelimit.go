package request

import (
	"context"
	"errors"
	"time"

	"golang.org/x/time/rate"
)

// Request 接口整理
type Request interface {
	// NewLimiter 初始化令牌信息；bucketSize 表示每秒产生令牌的个数，也即是每秒接收多少个请求
	NewLimiter(bucketSize int) (*RateLimit, error)
	// IsPass 令牌桶限流请求是否放行，参数表示获取令牌等待超时时间，单位：ms
	IsPass(timeOut time.Duration) bool
}

// RateLimit 对象结构体
type RateLimit struct {
	*rate.Limiter
}

// NewLimiter 初始化令牌信息；bucketSize 表示每秒产生令牌的个数，也即是每秒接收多少个请求
func NewLimiter(bucketSize int) (*RateLimit, error) {
	if bucketSize <= 0 {
		return nil, errors.New("param error")
	}

	// 计算多少毫秒释放一个请求
	millisecondNum := int(1000 / bucketSize)
	if millisecondNum < 1 {
		millisecondNum = 1
	}

	// 获取令牌对象，第一个参数为令牌生产时间间隔，第二个参数为令牌桶最大容量
	limiter := rate.NewLimiter(rate.Every(time.Duration(millisecondNum)*time.Millisecond), bucketSize)

	// 使用过程中可通过下面两个逻辑进行修改对应的两个参数
	// limiter.SetBurst(bucketSize)
	// limiter.SetLimit(10)

	return &RateLimit{limiter}, nil
}

// IsPass 令牌桶限流请求是否放行，参数表示获取令牌等待超时时间，单位：ms
func (rate *RateLimit) IsPass(timeOut time.Duration) bool {
	// 1、设置该执行逻辑的 context 超时时间
	ctx, cancel := context.WithTimeout(context.Background(), timeOut*time.Millisecond)
	defer cancel()

	// 2、等待监听是否有可用的令牌
	if err := rate.Wait(ctx); err != nil {
		return false
	}

	return true
}
