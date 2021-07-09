package store

import (
	"errors"

	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
)

// RedisInterface 接口整理
type RedisInterface interface {
	// NewRedis 获取 redis 对象
	NewRedis(host string, port string, passwd string) (*RedisCli, error)

	// Del redis del 方法
	Del(key string) bool
	// Expire redis expire 方法
	Expire(key string, seconds int) bool
	// Exists redis exists 方法
	Exists(key string) bool
	// TTL redis exists 方法
	TTL(key string) bool

	// Get redis get 方法
	Get(key string) (string, error)
	// Set redis set 方法
	Set(key string, value string) bool
	// SetEX redis setEX 方法
	SetEX(key string, value string, seconds int) bool
	// SetNX redis setNX 方法 <该方法只有在key不存在的时候才会设置成功>
	SetNX(key string, value string) bool
}

// RedisCli redis 对象结构
type RedisCli struct {
	client redis.Conn
}

// redisConfig redis 配置文件结构
type redisConfig struct {
	host   string
	port   string
	passwd string
}

// 全局变量定义
var (
	// redisCli redis 对象
	redisCli *RedisCli
)

// getRedis 初始化 redis
func getRedis(redisConf redisConfig) (*RedisCli, error) {
	// redis 配置获取
	if redisConf.host == "" || redisConf.port == "" {
		return nil, errors.New("ip/port is empty")
	}

	// redis 服务连接
	client, err := redis.Dial("tcp", redisConf.host+":"+redisConf.port)
	if err != nil {
		logrus.Warnf("redis.Dial err, err:%s", err.Error())
		return nil, err
	}

	// redis 密码鉴权
	if redisConf.passwd != "" {
		if _, err = client.Do("auth", redisConf.passwd); err != nil {
			logrus.Warnf("redis.Do auth err, err:%s", err.Error())
			client.Close()
			return nil, err
		}
		logrus.Infof("auth ok!, %s:%s", redisConf.host, redisConf.port)
	}
	logrus.Infof("connect to redis, %s:%s", redisConf.host, redisConf.port)

	// 赋值 redis 对象全局变量
	redisCli = &RedisCli{client: client}

	return redisCli, nil
}

// NewRedis 获取 redis 对象
func NewRedis(host string, port string, passwd string) (*RedisCli, error) {
	if host == "" || port == "" {
		logrus.Warnf("NewRedis params err")
		return nil, errors.New("params err")
	}

	redisConf := redisConfig{
		host:   host,
		port:   port,
		passwd: passwd,
	}
	return getRedis(redisConf)
}

// Set redis set 方法
func (conn *RedisCli) Set(key string, value string) bool {
	if key == "" || value == "" {
		logrus.Warnf("params error, key:%s, value:%s", key, value)
		return false
	}

	_, err := conn.client.Do("SET", key, value)
	if err != nil {
		logrus.Warnf("redis.Do SET err, err:%s", err.Error())
		return false
	}

	return true
}

// SetEX redis setEX 方法
func (conn *RedisCli) SetEX(key string, value string, seconds int) bool {
	if key == "" || value == "" || seconds <= 0 {
		logrus.Warnf("params error, key:%s, value:%s, second:%d", key, value, seconds)
		return false
	}

	_, err := conn.client.Do("SETEX", key, seconds, value)
	if err != nil {
		logrus.Warnf("redis.Do SETEX err, err:%s", err.Error())
		return false
	}

	return true
}

// SetNX redis setNX 方法 <该方法只有在key不存在的时候才会设置成功>
func (conn *RedisCli) SetNX(key string, value string) bool {
	if key == "" {
		logrus.Warnf("params error, key:%s", key)
		return false
	}

	ret, err := conn.client.Do("SETNX", key, value)
	if err != nil {
		logrus.Warnf("redis.Do SETNX err, ret:%d, err:%s", ret, err.Error())
		return false
	}
	if ret != int64(1) {
		logrus.Infof("redis.Do SETNX succ, ret:%d", ret)
		return false
	}

	return true
}

// Get redis get 方法
func (conn *RedisCli) Get(key string) (string, error) {
	if key == "" {
		logrus.Warnf("params error, key:%s", key)
		return "", errors.New("params error, key:" + key)
	}

	value, err := redis.String(conn.client.Do("GET", key))
	if err != nil {
		logrus.Warnf("redis.Do GET err, err:%s", err.Error())
		return "", err
	}

	return value, nil
}

// Del redis del 方法
func (conn *RedisCli) Del(key string) bool {
	if key == "" {
		logrus.Warnf("params error, key:%s", key)
		return false
	}

	_, err := conn.client.Do("DEL", key)
	if err != nil {
		logrus.Warnf("redis.Do DEL err, err:%s", err.Error())
		return false
	}

	return true
}

// Expire redis expire 方法
func (conn *RedisCli) Expire(key string, seconds int) bool {
	if key == "" || seconds <= 0 {
		logrus.Warnf("params error, key:%s, second:%d", key, seconds)
		return false
	}

	ret, err := conn.client.Do("EXPIRE", key, seconds)
	if ret != int64(1) || err != nil {
		logrus.Warnf("redis.Do EXPIRE err, ret:%d, err:%s", ret, err.Error())
		return false
	}

	return true
}

// Exists redis exists 方法
func (conn *RedisCli) Exists(key string) bool {
	if key == "" {
		logrus.Warnf("params error, key:%s", key)
		return false
	}

	isExists, err := redis.Bool(conn.client.Do("EXISTS", key))
	if err != nil {
		logrus.Warnf("redis.Do EXISTS err, err:%s", err.Error())
		return false
	}

	return isExists
}

// TTL redis ttl 方法
func (conn *RedisCli) TTL(key string) (int64, error) {
	if key == "" {
		logrus.Warnf("params error, key:%s", key)
		return -1, errors.New("params error, key:" + key)
	}

	intTime, err := redis.Int64(conn.client.Do("TTL", key))
	if err != nil {
		logrus.Warnf("redis.Do TTL err, err:%s", err.Error())
		return -1, err
	}

	return intTime, nil
}
