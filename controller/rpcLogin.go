package controller

import (
	"PackageServer/constant"
	"PackageServer/dto"
	"PackageServer/logger"
	"PackageServer/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type RpcLoginApi interface {
	RpcLogin(c *gin.Context)
	RetrieveRpcToken(c *gin.Context)
	RefreshTokenExipireTime(c *gin.Context)
	TokenDelete(c *gin.Context)
}

func NewRpcLoginApi() RpcLoginApi {
	return &rpcLoginApi{
		rpcLoginSvc:   service.NewRpcLoginService(),
		loginTokenSvc: service.NewLoginTokenService(),
	}
}

type rpcLoginApi struct {
	BaseApi
	rpcLoginSvc   service.RpcLoginSvc
	loginTokenSvc service.LoginTokenSvc
}

// --- client func ---

func (r *rpcLoginApi) RpcLogin(c *gin.Context) {
	area := c.Query("area")
	if area == "" {
		logger.Log.Errorf("%+v", constant.ErrParamIsNotComplete.Error())
		r.Error400(c, constant.ErrParamIsNotComplete)
		return
	}

	loginJwt, err := r.rpcLoginSvc.Login(area)
	if err != nil {
		logger.Log.Errorf("%+v", err)
		r.Error502(c, errors.Cause(err))
		return
	}

	logger.Log.Infof("login area %v successfully. jwt: %v", area, loginJwt)
	r.Success(c, "ok",
		map[string]string{
			"jwt": loginJwt,
		})
	return
}

// --- server func ---

func (r *rpcLoginApi) RetrieveRpcToken(c *gin.Context) {
	var loginTokenDto dto.LoginToken
	err := c.BindJSON(&loginTokenDto)
	//assert param is complete
	if err != nil {
		logger.Log.Errorf("%+v", constant.ErrBindJson)
		r.Error400(c, constant.ErrBindJson)
		return
	} else if loginTokenDto.LoginToken == "" || loginTokenDto.User == "" {
		logger.Log.Errorf("%v ERR: %v ", loginTokenDto, constant.ErrParamIsNotComplete.Error())
		r.Error400(c, constant.ErrParamIsNotComplete)
		return
	}

	logger.Log.Debugf("loginToken: %v", loginTokenDto)

	rpcToken, err := r.rpcLoginSvc.RetrieveRpcToken(loginTokenDto)
	if err != nil {
		logger.Log.Errorf("%+v", err)
		r.Error502(c, errors.Cause(err))
		return
	}

	logger.Log.Infof("jwt is created successfully! returning jwt")
	r.Success(c, "ok",
		map[string]string{
			"jwt": rpcToken,
		},
	)
	return
}

func (r *rpcLoginApi) RefreshTokenExipireTime(c *gin.Context) {
	var RpcTokenUpdateDto dto.RpcTokenUpdate
	c.BindJSON(&RpcTokenUpdateDto)

	validator := validator.New()
	err := validator.Struct(&RpcTokenUpdateDto)
	if err != nil {
		logger.Log.Errorf("%+v", err)
		r.Error400(c, constant.ErrParamIsNotComplete)
		return
	}

	if err := r.rpcLoginSvc.RpcTokenExpireRefresh(RpcTokenUpdateDto); err != nil {
		logger.Log.Errorf("%+v", err)
		r.Error502(c, errors.Cause(err))
		return
	}

	r.Success(c, "ok", nil)
}

// UNDONE
func (r *rpcLoginApi) TokenDelete(c *gin.Context) {
	area := c.Query("area")
	if area == "" {
		logger.Log.Errorf("%+v", constant.ErrParamIsNotComplete)
		r.Error400(c, constant.ErrParamIsNotComplete)
		return
	}
	err := r.rpcLoginSvc.RpcTokenDelete(area)
	if err != nil {
		logger.Log.Errorf("%+v", err)
		r.Error502(c, errors.Cause(err))
		return
	}
	logger.Log.Infof("delete rpcToken from: %v", area)
	r.Success(c, "ok", nil)
}
