package lock

import (
	"errors"

	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
)

// tryLock 分布式抢占锁对象
type tryLock struct {
	key     string
	value   string
	conn    redis.Conn
	timeout int
}

// prefixKey redis 分布式锁 key 前缀
const prefixKey = "try_lock_"

// defaultTimeout redis 分布式锁默认超时时间，单位:秒
const defaultTimeout = 5

// NewTryLock 基于 redis 实现的分布式抢占锁
func NewTryLock(host string, port string, key string, value string, timeout int) (*tryLock, error) {
	if host == "" || port == "" || key == "" || value == "" {
		logrus.Warnf("NewRedis params err")
		return nil, errors.New("params err")
	}

	// 获取 redis 对象
	redisCli, err := redis.Dial("tcp", host+":"+port)
	if err != nil {
		logrus.Warnf("redis.Dial err, err:%s", err.Error())
		return nil, err
	}

	if timeout <= 0 {
		timeout = defaultTimeout
	}
	lock := &tryLock{
		key:     prefixKey + key,
		value:   value,
		conn:    redisCli,
		timeout: timeout,
	}

	return lock, nil
}

// TryLock 尝试获取 redis 锁
func (lock *tryLock) TryLock() bool {
	_, err := redis.String(lock.conn.Do("SET", lock.key, lock.value, "EX", lock.timeout, "NX"))
	if err != nil {
		return false
	}

	return true
}

// UnLock 释放 redis 锁
func (lock *tryLock) UnLock() {
	if _, err := lock.conn.Do("DEL", lock.key); err != nil {
		logrus.Warnf("DEL err, err:%s", err.Error())
	}

	return
}
