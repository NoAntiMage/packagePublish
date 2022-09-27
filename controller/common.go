package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"PackageServer/dto"
	"PackageServer/logger"
	"PackageServer/repo"
	"PackageServer/service"
)

func Todo(c *gin.Context) {
	c.String(http.StatusOK, "todo")
}

type CommonApi interface {
	Ping(c *gin.Context)
	TimeStamp(c *gin.Context)
}

func NewCommonApi() CommonApi {
	return &commonApi{
		userRepo:      repo.NewUserRepo(),
		loginTokenSvc: service.NewLoginTokenService(),
	}
}

type commonApi struct {
	BaseApi
	userRepo      repo.UserRepo
	loginTokenSvc service.LoginTokenSvc
}

func (co *commonApi) Ping(c *gin.Context) {
	co.Success(c, "pong", nil)
}

func (co *commonApi) TimeStamp(c *gin.Context) {
	nowTimeStamp := strconv.Itoa(int(time.Now().Unix()))

	tsDto := dto.TimeStamp{
		Timestamp: nowTimeStamp,
	}

	user := c.Query("user")
	if user != "" {
		err := co.userRepo.Exist(user)
		if err != nil {
			logger.Log.Errorf("%+v", err)
			co.Success(c, err.Error(), tsDto)
			return
		}

		logger.Log.Infof("login user %v request timestamp", user)
		logger.Log.Debugf("user %v loginTokenSvc start", user)
		_, err = co.loginTokenSvc.CreateAndSaveToken(user, nowTimeStamp)
		if err != nil {
			logger.Log.Errorf("%+v", err)
		}
	}

	co.Success(c, "ok", tsDto)
}
