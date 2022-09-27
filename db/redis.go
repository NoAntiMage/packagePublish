package db

import (
	"PackageServer/config"
	"PackageServer/logger"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

var RedisConn *redis.Conn
var RedisPool *redis.Pool

func NewRedisConn() *redis.Conn {
	return RedisConn
}

func InitCache() {
	RedisConn, _ = initRedisConn()
	RedisPool = initRedisPool()
}

func initRedisConn() (*redis.Conn, error) {
	addr := fmt.Sprintf("%v:%v", config.RedisConf.Ip, config.RedisConf.Port)
	rConn, err := redis.Dial("tcp", addr)
	if err != nil {
		logger.Log.Error(err)
		rConn.Close()
		return nil, err
	}
	_, err = rConn.Do("AUTH", config.RedisConf.Password)
	if err != nil {
		logger.Log.Error(err)
		rConn.Close()
		return nil, err
	}
	return &rConn, nil
}

func initRedisPool() *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			addr := fmt.Sprintf("%v:%v", config.RedisConf.Ip, config.RedisConf.Port)
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			_, err = c.Do("AUTH", config.RedisConf.Password)
			if err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}
}
