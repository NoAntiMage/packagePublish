package service

import (
	"PackageServer/cache"
	"PackageServer/config"
	"PackageServer/constant"
	"PackageServer/dto"
	"PackageServer/logger"
	"PackageServer/rpc"
	"fmt"
)

type RpcLoginSvc interface {
	Login(area string) (token string, err error)
	Logout() error
	RetrieveRpcToken(loginTokenDto dto.LoginToken) (rpcToken string, err error)
	RemoteAreaTokenCheck(area string) error
	RpcTokenExpireRefresh(rpcTokenUpdateDto dto.RpcTokenUpdate) error
	RpcTokenDelete(area string) error
}

func NewRpcLoginService() RpcLoginSvc {
	return &rpcLoginSvc{
		LoginTokenSvc: NewLoginTokenService(),
		areaInfoSvc:   NewAreaInfoSvc(),
		jwtCache:      cache.NewJsonWebTokenCache(),
	}
}

type rpcLoginSvc struct {
	LoginTokenSvc LoginTokenSvc
	areaInfoSvc   AreaInfoSvc
	jwtCache      cache.JsonWebTokenCache
}

// --- client func ---

/*Description:
1. client and server hold the Secret
2. client request TimeStamp with loginUser
3. client generate loginToken = func(TimeStamp + Secret + Password)
4. client post loginToken to server
5. server validate the loginToken and return the JWT as RpcToken.

actually we implement digestAuth and session mode.
jwt is regarded as sessionSecret.
*/
func (l *rpcLoginSvc) Login(area string) (jwt string, err error) {
	var loginTokenDto dto.LoginToken
	loginTokenDto.User = config.ServerConf.Name
	loginTokenDto.Area = area

	areaInfoDto, err := l.areaInfoSvc.Get(area)
	if err != nil {
		return "", err
	}

	ts, err := rpc.RetrieveTimeStamp(loginTokenDto, *areaInfoDto)
	logger.Log.Debugf("service:rpcLoginSvc: get timeStamp: %v ", ts)
	if err != nil {
		return "", err
	}

	secret := constant.Secret
	password := areaInfoDto.Password

	loginTokenDto.LoginToken = l.LoginTokenSvc.GenerateToken(ts.Timestamp, secret, password)
	logger.Log.Debugf("service:rpcLoginSvc: caculate digestToken: %v", loginTokenDto.LoginToken)

	jwt, err = rpc.RetrieveJwt(loginTokenDto, *areaInfoDto)
	if err != nil {
		return jwt, err
	}

	direction := string(constant.To)
	err = l.jwtCache.CacheJwt(area, jwt, direction)
	if err != nil {
		return jwt, err
	}
	return jwt, nil
}

func (l *rpcLoginSvc) Logout() error {
	return nil
}

/* TODO HERE 2
Token to remote should be checked
1. is localStorage token exist in remote
2. if not exist, client require a new token from remote
3. save token to storage(cache for now)
*/
func (l *rpcLoginSvc) RemoteAreaTokenCheck(area string) error {
	areaInfoDto, err := l.areaInfoSvc.Get(area)
	if err != nil {
		return err
	}
	ttl := constant.RpcDefaultExpireTime
	remoteUrl := fmt.Sprintf("http://%v:%v", areaInfoDto.IpAddr, areaInfoDto.Port)
	if err := rpc.UpdateRpcTokenExpireTime(area, remoteUrl, ttl); err != nil {
		logger.Log.Debugf("rpcLoginSvc:RemoteAreaTokenCheck:UpdateRpcTokenExpireTime Fail: ERR: %v", err)
		_, err := l.Login(area)
		return err
	}
	return nil
}

// --- server func ---

func (l *rpcLoginSvc) RetrieveRpcToken(loginTokenDto dto.LoginToken) (rpcToken string, err error) {
	err = l.LoginTokenSvc.ValidateToken(loginTokenDto)
	if err != nil {
		return "", err
	}

	jwt, err := l.jwtCache.GenerateJwt(loginTokenDto.User)
	if err != nil {
		return "", err
	}

	direction := string(constant.From)
	err = l.jwtCache.CacheJwt(loginTokenDto.User, jwt, direction)
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func (l *rpcLoginSvc) RpcTokenExpireRefresh(rpcTokenUpdateDto dto.RpcTokenUpdate) error {
	if rpcTokenUpdateDto.ExpireTime == 0 {
		rpcTokenUpdateDto.ExpireTime = constant.RpcDefaultExpireTime
	}
	from := string(constant.From)
	return l.jwtCache.UpdateExpireTimeOfJwt(rpcTokenUpdateDto.User, from, rpcTokenUpdateDto.ExpireTime)
}

func (l *rpcLoginSvc) RpcTokenDelete(user string) error {
	from := string(constant.From)
	return l.jwtCache.DeleteJwtCache(user, from)
}
