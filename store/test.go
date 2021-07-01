package main

import (
	"errors"
	"fmt"

	"github.com/garyburd/redigo/redis"
)

type RedisClient struct {
	client redis.Conn
}

var redisClient RedisClient

// 初始化 redis
func initRedis() (*RedisClient, error) {
	// 获取 redis 配置信息
	host := "127.0.0.1"
	port := "6379"
	passwd := ""

	// 连接 redis
	client, err1 := redis.Dial("tcp", host+":"+port)
	if err1 != nil {
		return nil, err1
	}

	// 验证 redis
	_, err2 := client.Do("auth", passwd)
	if err2 != nil {
		return nil, err2
	}

	redisClient = RedisClient{
		client: client,
	}
	return &redisClient, nil
}

func (redis *RedisClient) Set(key string, value string) error {
	if key == "" || value == "" {
		return errors.New("params error")
	}

	_, err := redis.client.Do("Set", key, value)
	if err != nil {
		return err
	}

	return nil
}

func (redis *RedisClient) Get(key string) (string, error) {
	if key == "" {
		return "", errors.New("params error")
	}

	value, err := redis.client.Do("Get", key)
	if err != nil {
		return "", err
	}

	return string(value.([]byte)), nil
}

func main() {
	initRedis()

	err := redisClient.Set("test", "abc")
	if err != nil {
		return
	}

	val, err := redisClient.Get("test")
	if err != nil {
		return
	}

	fmt.Println(val)
}
