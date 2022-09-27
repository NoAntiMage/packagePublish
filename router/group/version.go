package group

import (
	"PackageServer/constant"
	"PackageServer/controller"

	"github.com/gin-gonic/gin"
)

func NewVersionGroup() *VersionGroup {
	return &VersionGroup{
		groupPath:   constant.UrlVersion,
		RpcLoginApi: controller.NewRpcLoginApi(),
	}
}

type VersionGroup struct {
	groupPath   string
	RpcLoginApi controller.RpcLoginApi
}

func (vg *VersionGroup) InitRouter(rg *gin.RouterGroup) *gin.RouterGroup {
	versionRouter := rg.Group(vg.groupPath)

	//	versionRouter.GET("/healthcheck", controller.HealthCheck)
	versionRouter.POST("/digestToken", vg.RpcLoginApi.RetrieveRpcToken)
	versionRouter.GET("/rpcLogin", vg.RpcLoginApi.RpcLogin) // internal rpc for DEBUG
	return versionRouter
}
