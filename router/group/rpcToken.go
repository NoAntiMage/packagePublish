package group

import (
	"PackageServer/controller"

	"github.com/gin-gonic/gin"
)

type RpcTokenGroup struct {
	groupPath   string
	rpcLoginApi controller.RpcLoginApi
}

func NewRpcTokenGroup() *RpcTokenGroup {
	return &RpcTokenGroup{
		groupPath:   "rpcToken",
		rpcLoginApi: controller.NewRpcLoginApi(),
	}
}

func (r *RpcTokenGroup) InitRouter(rg *gin.RouterGroup) *gin.RouterGroup {
	rpcTokenRouter := rg.Group(r.groupPath)

	rpcTokenRouter.POST("/refreshExpire", r.rpcLoginApi.RefreshTokenExipireTime)
	rpcTokenRouter.GET("/del", r.rpcLoginApi.TokenDelete)
	return rpcTokenRouter
}
