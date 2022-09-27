package group

import (
	"PackageServer/controller"

	"github.com/gin-gonic/gin"
)

func NewServiceOnlineGroup() *ServiceOnlineGroup {
	return &ServiceOnlineGroup{
		groupPath:        "serviceOnline",
		ServiceOnlineApi: controller.NewServiceOnlineApi(),
	}
}

type ServiceOnlineGroup struct {
	groupPath        string
	ServiceOnlineApi controller.ServiceOnlineApi
}

func (sg *ServiceOnlineGroup) InitRouter(rg *gin.RouterGroup) *gin.RouterGroup {
	serviceOnlineRouter := rg.Group(sg.groupPath)

	serviceOnlineRouter.POST("", sg.ServiceOnlineApi.Create)
	serviceOnlineRouter.GET("", sg.ServiceOnlineApi.List)
	serviceOnlineRouter.DELETE("/:serviceName", sg.ServiceOnlineApi.DeleteByName)
	serviceOnlineRouter.GET("/:serviceName", sg.ServiceOnlineApi.GetByName)
	serviceOnlineRouter.PUT("/:serviceName", sg.ServiceOnlineApi.UpdateByName)
	serviceOnlineRouter.GET("/:serviceName/areas", sg.ServiceOnlineApi.ListAreas)
	return serviceOnlineRouter
}
