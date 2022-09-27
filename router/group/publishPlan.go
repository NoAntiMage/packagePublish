package group

import (
	"PackageServer/controller"

	"github.com/gin-gonic/gin"
)

type PublishPlanGroup struct {
	groupPath      string
	PublishPlanApi controller.PublishPlanApi
}

func NewPublishPlanGroup() *PublishPlanGroup {
	return &PublishPlanGroup{
		groupPath:      "publishPlan",
		PublishPlanApi: controller.NewPublishPlanApi(),
	}
}

func (pg *PublishPlanGroup) InitRouter(rg *gin.RouterGroup) *gin.RouterGroup {
	publishPlanRouter := rg.Group(pg.groupPath)

	publishPlanRouter.GET("", controller.Todo)
	publishPlanRouter.POST("", pg.PublishPlanApi.PublishVersion)
	publishPlanRouter.GET("/:planName", controller.Todo)

	return publishPlanRouter
}
