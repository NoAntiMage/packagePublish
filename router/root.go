package router

import (
	"PackageServer/config"
	"PackageServer/logger"
	"PackageServer/middleware"
	group "PackageServer/router/group"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Engine              *gin.Engine
	CommonGroup         *group.CommonGroup
	VerionGroup         *group.VersionGroup
	AuthGroup           *group.AuthGroup
	RpcTokenGroup       *group.RpcTokenGroup
	AreaInfoGroup       *group.AreaInfoGroup
	ServiceOnlineGroup  *group.ServiceOnlineGroup
	PublishPlanGroup    *group.PublishPlanGroup
	PackageReceiveGroup *group.PackageReceiveGroup
}

func NewRouter() *Router {
	return &Router{
		Engine:              gin.New(),
		CommonGroup:         group.NewCommonGroup(),
		VerionGroup:         group.NewVersionGroup(),
		AuthGroup:           group.NewAuthGroup(),
		RpcTokenGroup:       group.NewRpcTokenGroup(),
		AreaInfoGroup:       group.NewAreaInfoGroup(),
		ServiceOnlineGroup:  group.NewServiceOnlineGroup(),
		PublishPlanGroup:    group.NewPublishPlanGroup(),
		PackageReceiveGroup: group.NewPackageReceiveGroup(),
	}
}

func (r *Router) RouterInit() {
	r.Engine.Use(gin.RecoveryWithWriter(logger.FileAndStdoutWriter))
	r.Engine.Use(middleware.LoggerToFileAndStdout())

	groupRoot := r.Engine.Group("")

	r.CommonGroup.InitRouter(groupRoot)
	versionRouter := r.VerionGroup.InitRouter(groupRoot)
	authRouter := r.AuthGroup.InitRouter(versionRouter)
	r.RpcTokenGroup.InitRouter(authRouter)

	if config.ServerConf.Name == "manager" {
		r.AreaInfoGroup.InitRouter(authRouter)
		r.ServiceOnlineGroup.InitRouter(authRouter)
		r.PublishPlanGroup.InitRouter(authRouter)

	} else if config.ServerConf.Name == "worker" {
		r.PackageReceiveGroup.InitRouter(authRouter)
	}
}
