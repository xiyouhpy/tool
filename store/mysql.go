package store

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

// MysqlCli mysql 对象结构
type MysqlCli struct {
	client *sql.DB
}

// mysqlConfig mysql 配置文件结构
type mysqlConfig struct {
	host   string
	port   string
	user   string
	passwd string
	dbname string
}

// 全局变量定义
var (
	// mysqlCli mysql 对象
	mysqlCli *MysqlCli
)

// getMysql 初始化 mysql，使用 utf-8 编码
func getMysql(mysqlConf mysqlConfig) (*MysqlCli, error) {
	// mysql 配置获取
	if mysqlConf.host == "" || mysqlConf.port == "" || mysqlConf.user == "" || mysqlConf.passwd == "" {
		return nil, errors.New("ip/port is empty")
	}

	// mysql 服务连接
	dbServer := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8",
		mysqlConf.user, mysqlConf.passwd, mysqlConf.host, mysqlConf.port, mysqlConf.dbname)
	client, err := sql.Open("mysql", dbServer)
	if err != nil {
		logrus.Warnf("mysql Open err, err:%s", err.Error())
		return nil, err
	}

	if err = client.Ping(); err != nil {
		logrus.Warnf("ping mysql err, err:%s", err.Error())
		return nil, err
	}

	// 赋值 redis 对象全局变量
	mysqlCli = &MysqlCli{client: client}

	return mysqlCli, nil
}

// NewMysql 获取 mysql 对象
func NewMysql(host string, port string, user string, passwd string, dbname string) (*MysqlCli, error) {
	if host == "" || port == "" || user == "" || passwd == "" || dbname == "" {
		logrus.Warnf("NewMysql params err")
		return nil, errors.New("params err")
	}

	mysqlConf := mysqlConfig{
		host:   host,
		port:   port,
		user:   user,
		passwd: passwd,
		dbname: dbname,
	}
	return getMysql(mysqlConf)
}
