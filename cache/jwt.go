package cache

import (
	"PackageServer/constant"
	"PackageServer/logger"
	"PackageServer/util"
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrJwtNotMatch = errors.New("ErrJwtNotMatch")
)

type JsonWebTokenCache interface {
	GenerateKey(area string, direction string) (key string)
	GenerateJwt(area string) (string, error)
	CacheJwt(area string, jwt string, direction string) error
	DeleteJwtCache(area string, direction string) error
	GetJwt(area string, direction string) (value string, err error)
	ValidateJwt(area string, jwt string, direction string) error
	UpdateExpireTimeOfJwt(area string, direction string, ttl int) error
}

func NewJsonWebTokenCache() JsonWebTokenCache {
	return &jsonWebTokenCache{
		Jwt:       util.NewJwt(),
		RedisUtil: util.NewRedisUtil(),
	}
}

type jsonWebTokenCache struct {
	Claim     util.CustomClaims
	Jwt       util.Jwt
	RedisUtil util.RedisUtil
}

func (j *jsonWebTokenCache) GenerateKey(area string, direction string) (key string) {
	key = fmt.Sprintf("%v%v:%v", constant.JwtTokenKey, direction, area)
	logger.Log.Debugf("jwt key: %v", key)
	return key
}

func (j *jsonWebTokenCache) GenerateJwt(area string) (jwt string, err error) {
	j.Claim.User = area
	return j.Jwt.CreateTokenWithExpire(j.Claim)
}

func (j *jsonWebTokenCache) CacheJwt(area string, jwt string, direction string) error {
	key := j.GenerateKey(area, direction)
	//TODO read from config
	JwtExpireSeconds := 30 * 60
	return j.RedisUtil.SetExpire(key, jwt, JwtExpireSeconds)
}

func (j *jsonWebTokenCache) GetJwt(area string, direction string) (value string, err error) {
	key := j.GenerateKey(area, direction)
	value, err = j.RedisUtil.Get(key)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (j *jsonWebTokenCache) ValidateJwt(area string, jwt string, direction string) error {
	key := j.GenerateKey(area, direction)
	value, err := j.RedisUtil.Get(key)
	if err != nil {
		return err
	}

	if value != jwt {
		return errors.Wrapf(ErrJwtNotMatch, "jsonWebTokenCache:ValidateJwt")
	}
	return nil
}

func (j *jsonWebTokenCache) UpdateExpireTimeOfJwt(area string, direction string, ttl int) error {
	key := j.GenerateKey(area, direction)
	return j.RedisUtil.Expire(key, ttl)
}

func (j *jsonWebTokenCache) DeleteJwtCache(area string, direction string) error {
	key := j.GenerateKey(area, direction)
	return j.RedisUtil.Del(key)
}
