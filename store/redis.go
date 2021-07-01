package store

import (
	"errors"

	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
)

// RedisClient redis 对象结构
type RedisClient struct {
	client redis.Conn
}

// RedisInterface 接口整理
type RedisInterface interface {
	// NewRedis 获取 redis 对象
	NewRedis() (*RedisClient, error)
	// Set redis set 方法
	Set(key string, value string) error
	// Get redis get 方法
	Get(key string) (string, error)
}

// redisClient redis 对象全局变量
var redisClient *RedisClient

// getRedis 初始化 redis
func getRedis() (*RedisClient, error) {
	// 获取 redis 配置信息
	host := "127.0.0.1"
	port := "6379"
	passwd := ""

	// 连接 redis
	client, err := redis.Dial("tcp", host+":"+port)
	if err != nil {
		logrus.Warnf("redis.Dial err, err:%s", err.Error())
		return nil, err
	}

	// 验证 redis
	if passwd != "" {
		if _, err = client.Do("auth", passwd); err != nil {
			logrus.Warnf("redis.Do auth err, err:%s", err.Error())
			return nil, err
		}
	}

	// 赋值 redis 对象全局变量
	redisClient = &RedisClient{client: client}

	return redisClient, nil
}

// NewRedis 获取 redis 对象
func NewRedis() (*RedisClient, error) {
	if redisClient != nil {
		return redisClient, nil
	}

	return getRedis()
}

// Set redis set 方法
func (redis *RedisClient) Set(key string, value string) error {
	if key == "" || value == "" {
		return errors.New("params error")
	}

	_, err := redis.client.Do("Set", key, value)
	if err != nil {
		logrus.Warnf("redis.Do Set err, err:%s", err.Error())
		return err
	}

	return nil
}

// Get redis get 方法
func (redis *RedisClient) Get(key string) (string, error) {
	if key == "" {
		return "", errors.New("params error")
	}

	value, err := redis.client.Do("Get", key)
	if err != nil {
		logrus.Warnf("redis.Do Get err, err:%s", err.Error())
		return "", err
	}

	return string(value.([]byte)), nil
}
