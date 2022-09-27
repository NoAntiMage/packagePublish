package group

import (
	"PackageServer/controller"

	"github.com/gin-gonic/gin"
)

type BaseGroup interface {
	InitRouter(rg *gin.RouterGroup) *gin.RouterGroup
}

func NewCommonGroup() *CommonGroup {
	return &CommonGroup{
		CommonApi: controller.NewCommonApi(),
	}
}

type CommonGroup struct {
	CommonApi controller.CommonApi
}

func (cg *CommonGroup) InitRouter(rg *gin.RouterGroup) *gin.RouterGroup {
	rg.GET("/ping", cg.CommonApi.Ping)
	rg.GET("/timestamp", cg.CommonApi.TimeStamp)
	return rg
}
