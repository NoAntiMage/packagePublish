package rpc

import (
	"PackageServer/cache"
	"PackageServer/config"
	"PackageServer/constant"
)

type AuthRpc struct {
	UrlPath string
}

func NewAuthRpc() AuthRpc {
	return AuthRpc{
		UrlPath: constant.UrlVersion + constant.UrlAuth,
	}
}

func (b *AuthRpc) GetTokenByEnv(area string) (string, error) {
	var token string
	if config.ServerConf.RouterMode != "release" {
		return token, nil
	}

	jwtCache := cache.NewJsonWebTokenCache()
	to := string(constant.To)
	token, err := jwtCache.GetJwt(area, to)
	if err != nil {
		return "", err
	}

	return token, nil
}
