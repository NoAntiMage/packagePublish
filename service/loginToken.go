package service

import (
	"PackageServer/constant"
	"PackageServer/dto"
	"PackageServer/logger"
	"PackageServer/util"
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrLoginTokenNotMatch     = errors.New("login token not match")
	ErrPasswordConfigNotFound = errors.New("password config not found")
)

type LoginTokenSvc interface {
	GenerateToken(nowTimeStamp, secret, password string) (value string)
	GenerateKey(user string) string
	SaveToken(key, value string) error
	CreateAndSaveToken(user string, nowTimeStamp string) (string, error)
	ValidateToken(dto.LoginToken) error
	//	updateToken()
	DeleteToken(user string) error
}

func NewLoginTokenService() LoginTokenSvc {
	return &loginTokenSvc{
		md5Util:   util.NewMd5Util(),
		redisUtil: util.NewRedisUtil(),
	}
}

type loginTokenSvc struct {
	md5Util   util.Md5Util
	redisUtil util.RedisUtil
}

func (l *loginTokenSvc) GenerateToken(nowTimeStamp, secret, password string) (value string) {
	value = l.md5Util.Md5Sum(nowTimeStamp + secret + password)
	logger.Log.Debugf("value is : %v", value)
	return value
}

func (l *loginTokenSvc) GenerateKey(user string) (key string) {
	key = fmt.Sprintf("%v%v", constant.LoginTokenKey, user)
	logger.Log.Debugf("key is : %v", key)
	return key
}

func (l *loginTokenSvc) SaveToken(key, value string) error {
	//TODO expire from config
	expireSeconds := 180
	if err := l.redisUtil.SetExpire(key, value, expireSeconds); err != nil {
		return err
	}
	return nil
}

func (l *loginTokenSvc) CreateAndSaveToken(user string, nowTimeStamp string) (string, error) {
	secret := constant.Secret
	password, ok := util.AppConfig.Get("password").(string)
	if ok == false {
		return "", errors.Wrap(ErrPasswordConfigNotFound, "service:LoginTokenService:CreateAndSaveToken:")
	}

	key := l.GenerateKey(user)
	value := l.GenerateToken(nowTimeStamp, secret, password)

	err := l.SaveToken(key, value)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (l *loginTokenSvc) ValidateToken(loginToken dto.LoginToken) error {
	key := l.GenerateKey(loginToken.User)
	value, err := l.redisUtil.Get(key)
	if err != nil {
		logger.Log.Infof("redis key %v is not found", key)
		return err
	}

	logger.Log.Debugf("value is : %v, loginToken is : %v", value, loginToken.LoginToken)
	if value == loginToken.LoginToken {
		logger.Log.Debugf("%v loginToken validation pass", loginToken.User)
		return nil
	} else {
		return errors.Wrap(ErrLoginTokenNotMatch, "service:loginToken:ValidateToken:")
	}
}

func (l *loginTokenSvc) DeleteToken(user string) error {
	return nil
}
