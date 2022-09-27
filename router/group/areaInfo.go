package group

import (
	"PackageServer/controller"

	"github.com/gin-gonic/gin"
)

func NewAreaInfoGroup() *AreaInfoGroup {
	return &AreaInfoGroup{
		groupPath:   "areaInfo",
		AreaInfoApi: controller.NewAreaInfoApi(),
	}
}

type AreaInfoGroup struct {
	groupPath   string
	AreaInfoApi controller.AreaInfoApi
}

func (ag AreaInfoGroup) InitRouter(rg *gin.RouterGroup) *gin.RouterGroup {
	areaInfoGroup := rg.Group(ag.groupPath)

	areaInfoGroup.POST("", ag.AreaInfoApi.Create)
	areaInfoGroup.GET("", ag.AreaInfoApi.List)
	areaInfoGroup.DELETE("/:areaName", ag.AreaInfoApi.DeleteByName)
	areaInfoGroup.GET("/:areaName", ag.AreaInfoApi.GetByName)
	areaInfoGroup.PUT("/:areaName", ag.AreaInfoApi.UpdateByName)

	areaInfoGroup.GET("/:areaName/services", ag.AreaInfoApi.ListServices)
	areaInfoGroup.POST("/:areaName/services/add", ag.AreaInfoApi.AddServices)
	areaInfoGroup.POST("/:areaName/services/delete", ag.AreaInfoApi.DelServices)

	areaInfoGroup.GET("/:areaName/service/:serviceName/publishPlan", controller.Todo)
	areaInfoGroup.GET("/:areaName/service/:serviceName/logJob", controller.Todo)

	return areaInfoGroup
}
