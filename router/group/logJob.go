package group

import (
	"PackageServer/controller"

	"github.com/gin-gonic/gin"
)

type LogJob struct {
	groupPath string
}

func NewLogJob() *LogJob {
	return &LogJob{
		groupPath: "logJob",
	}
}

func (lg *LogJob) InitRouter(rg *gin.RouterGroup) *gin.RouterGroup {
	logJobRouter := rg.Group(lg.groupPath)
	logJobRouter.GET("job", controller.Todo)
	logJobRouter.POST("job", controller.Todo)
	logJobRouter.GET("job/:jobName", controller.Todo)
	logJobRouter.GET("job/:jobName/receive", controller.Todo) // for callBack, not web
	logJobRouter.GET("job/:jobName/retry", controller.Todo)
	logJobRouter.GET("job/:jobName/:fileName/download", controller.Todo)

	return logJobRouter
}
