package group

import (
	"PackageServer/config"
	"PackageServer/middleware"

	"github.com/gin-gonic/gin"
)

func NewAuthGroup() *AuthGroup {
	return &AuthGroup{
		groupPath: "auth",
	}
}

type AuthGroup struct {
	groupPath string
}

func (ag *AuthGroup) InitRouter(rg *gin.RouterGroup) *gin.RouterGroup {
	authRouter := rg.Group(ag.groupPath)

	if config.ServerConf.RouterMode == "release" && config.ServerConf.Name == "worker" {
		authRouter.Use(middleware.JwtAuth())
	}

	return authRouter
}
