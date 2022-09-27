package rpc

import (
	"PackageServer/constant"
	"PackageServer/dto"
	"PackageServer/logger"
	"PackageServer/request"

	"github.com/pkg/errors"
)

func RetrieveTimeStamp(loginTokenDto dto.LoginToken, areaInfoDto dto.AreaInfo) (*dto.TimeStamp, error) {
	ts, err := request.GetTimeStamp(loginTokenDto, areaInfoDto)
	if err != nil {
		return nil, err

	}
	logger.Log.Debugf("ts: %v", ts)
	return ts, nil
}

func RetrieveJwt(loginTokenDto dto.LoginToken, areaInfoDto dto.AreaInfo) (jwt string, err error) {
	if loginTokenDto.User == "" || loginTokenDto.LoginToken == "" {
		return "", errors.Wrap(constant.ErrParamIsNotComplete, "rpc:RetrieveJwt:")
	}
	jwt, err = request.PostDigestToken(loginTokenDto, areaInfoDto)
	if err != nil {
		return jwt, err
	}
	return jwt, nil
}

func UpdateRpcTokenExpireTime(area string, remoteUrl string, ttl int) error {
	authRpc := NewAuthRpc()
	token, err := authRpc.GetTokenByEnv(area)
	if err != nil {
		return err
	}
	var reqInfoDto = dto.RequestInfo{
		TargetUrl: remoteUrl + authRpc.UrlPath + constant.UrlRpcTokenRefreshExpire,
		Token:     token,
	}
	return request.PostRpcTokenExpireTime(reqInfoDto, ttl)
}
