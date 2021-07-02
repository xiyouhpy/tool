package store

import (
	"errors"
	"io/ioutil"

	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// RedisInterface 接口整理
type RedisInterface interface {
	// NewRedis 获取 redis 对象
	NewRedis(redisCluster string, configPath string) (*RedisCli, error)
	// Set redis set 方法
	Set(key string, value string) bool
	// SetEX redis setEX 方法
	SetEX(key string, value string, seconds int) bool
	// SetNX redis setNX 方法 <该方法只有在key不存在的时候才会设置成功>
	SetNX(key string, value string) bool
	// Get redis get 方法
	Get(key string) (string, error)
	// Del redis del 方法
	Del(key string) bool
	// Expire redis expire 方法
	Expire(key string, seconds int) bool
	// Exists redis exists 方法
	Exists(key string) bool
}

// redisCli redis 对象全局变量
var redisCli *RedisCli

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

// getConf 获取 redis 配置信息
func getConf(configPath string) map[string]redisConfig {
	yamlInfo, err := ioutil.ReadFile(configPath)
	if err != nil {
		logrus.Warnf("ioutil.ReadFile err, path:%s, err:%s", configPath, err.Error())
		return nil
	}

	config := make(map[string]redisConfig)
	err = yaml.Unmarshal(yamlInfo, config)
	if err != nil {
		logrus.Warnf("yaml.Unmarshal err, path:%s, err:%s", configPath, err.Error())
		return nil
	}

	return config
}

// getRedis 初始化 redis
func getRedis(redisCluster string, configPath string) (*RedisCli, error) {
	// redis 配置获取
	config := getConf(configPath)
	if config[redisCluster].host == "" || config[redisCluster].port == "" {
		logrus.Warnf("getConf err, conf:%v", config)
		return nil, errors.New("get host/port empty")
	}
	redisClusterConf := redisConfig{
		host:   config[redisCluster].host,
		port:   config[redisCluster].port,
		passwd: config[redisCluster].passwd,
	}

	// redis 服务连接
	client, err := redis.Dial("tcp", redisClusterConf.host+":"+redisClusterConf.port)
	if err != nil {
		logrus.Warnf("redis.Dial err, err:%s", err.Error())
		return nil, err
	}

	// redis 密码鉴权
	if redisClusterConf.passwd != "" {
		if _, err = client.Do("auth", redisClusterConf.passwd); err != nil {
			logrus.Warnf("redis.Do auth err, err:%s", err.Error())
			return nil, err
		}
		logrus.Infof("auth ok!, %s:%s", redisClusterConf.host, redisClusterConf.port)
	}
	logrus.Infof("connect to redis, %s:%s", redisClusterConf.host, redisClusterConf.port)

	// 赋值 redis 对象全局变量
	redisCli = &RedisCli{client: client}

	return redisCli, nil
}

// NewRedis 获取 redis 对象
func NewRedis(redisCluster string, configPath string) (*RedisCli, error) {
	if redisCli != nil {
		return redisCli, nil
	}

	return getRedis(redisCluster, configPath)
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

	_, err := conn.client.Do("SET", key, value, "EX", string(seconds))
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
	if ret != int64(1) || err != nil {
		logrus.Warnf("redis.Do SETNX err, ret:%d, err:%s", ret, err.Error())
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
