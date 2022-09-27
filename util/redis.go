package util

import (
	"PackageServer/db"
	"PackageServer/logger"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

type RedisUtil interface {
	Ping() (string, error)
	Set(key, value string) error
	SetExpire(key, value string, second int) error
	Get(key string) (string, error)
	Del(key string) error
	Expire(key string, second int) error
}

var (
	ErrCacheMiss = errors.New("ErrCacheMiss")
)

func NewRedisUtil() RedisUtil {
	return &redisUtil{
		Conn: *db.RedisConn,
	}
}

type redisUtil struct {
	Conn redis.Conn
}

func (r *redisUtil) Ping() (string, error) {
	reply, err := redis.String(r.Conn.Do("ping"))
	if err != nil {
		return "", errors.Wrap(err, "util:redisUtil:Ping")
	}
	return reply, nil
}

func (r *redisUtil) Set(key, value string) error {
	reply, err := redis.String(r.Conn.Do("SET", key, value))
	logger.Log.Debugf("redisKeySet %v %v", key, reply)
	if err != nil {
		return errors.Wrap(err, "util:redisUtil:Set")
	}
	return nil
}

func (r *redisUtil) SetExpire(key, value string, second int) error {
	reply, err := redis.String(r.Conn.Do("SETEX", key, second, value))
	logger.Log.Debugf("redisKeySet %v %v", key, reply)

	if err != nil {
		return errors.Wrap(err, "util:redisUtil:SetExpire")
	}
	return nil
}

func (r *redisUtil) Get(key string) (string, error) {
	reply, err := redis.String(r.Conn.Do("GET", key))
	logger.Log.Debugf("redisKeyGet %v %v", key, reply)
	if err != nil {
		if errors.Is(err, redis.ErrNil) {
			return "", errors.Wrap(ErrCacheMiss, "service:redisServer:Get")
		}
		return reply, errors.Wrap(err, "service:redisServer:Get")
	}
	return reply, nil
}

func (r *redisUtil) Del(key string) error {
	reply, err := redis.Int64(r.Conn.Do("DEL", key))
	logger.Log.Debugf("redisKeyDel %v %v", key, reply)
	if err != nil {
		return errors.Wrap(err, "util:redisUtil:Get")
	}
	return nil
}

func (r *redisUtil) Expire(key string, second int) error {
	reply, err := redis.Int64(r.Conn.Do("EXPIRE", key, second))
	logger.Log.Debugf("redisKeyExpire %v %v, expireTime: %v", key, reply, second)
	if err != nil {
		return errors.Wrapf(err, "util:redisUtil:Expire")
	}
	return nil
}
