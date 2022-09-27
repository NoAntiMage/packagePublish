package group

import (
	"PackageServer/controller"

	"github.com/gin-gonic/gin"
)

type PackageReceiveGroup struct {
	groupPath         string
	PackageReceiveApi controller.PackageReceiveApi
}

func NewPackageReceiveGroup() *PackageReceiveGroup {
	return &PackageReceiveGroup{
		groupPath:         "package",
		PackageReceiveApi: controller.NewPackageReceiveApi(),
	}
}

func (pg *PackageReceiveGroup) InitRouter(rg *gin.RouterGroup) *gin.RouterGroup {
	packageReceiveRouter := rg.Group(pg.groupPath)

	packageReceiveRouter.POST("/info", pg.PackageReceiveApi.PackInfo)
	packageReceiveRouter.POST("chunkUpload", pg.PackageReceiveApi.ChunkUpload)
	packageReceiveRouter.GET("/check", pg.PackageReceiveApi.PackCheck)

	return packageReceiveRouter
}
