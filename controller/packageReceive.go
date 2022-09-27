package controller

import (
	"PackageServer/constant"
	"PackageServer/dto"
	"PackageServer/service"

	"PackageServer/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type PackageReceiveApi interface {
	PackInfo(c *gin.Context)
	ChunkUpload(c *gin.Context)
	PackCheck(c *gin.Context)
}

func NewPackageReceiveApi() PackageReceiveApi {
	return &packageReceiveApi{
		PackageMgtSvc: service.NewPackageManagementSvc(),
	}
}

type packageReceiveApi struct {
	BaseApi
	PackageMgtSvc service.PackageManagementSvc
}

func (p *packageReceiveApi) PackInfo(c *gin.Context) {
	var d dto.PackageInfoPost
	c.BindJSON(&d)

	validation := validator.New()
	err := validation.Struct(d)
	if err != nil {
		logger.Log.Errorf("%+v", err)
		p.Error400(c, errors.Cause(err))
		return
	}

	if err = p.PackageMgtSvc.PackInfoReceive(d); err != nil {
		logger.Log.Errorf("%+v", err)
		p.Error502(c, errors.Cause(err))
		return
	}
	p.Success(c, "ok", nil)
}

func (p *packageReceiveApi) ChunkUpload(c *gin.Context) {
	var chunkInfo dto.ChunkInfo
	chunkInfo.FileName = c.Request.Header.Get("packageName")
	file, err := c.FormFile("chunk")
	if err != nil {
		logger.Log.Errorf("%+v", err)
		p.Error400(c, err)
		return
	}
	chunkInfo.ChunkName = file.Filename
	logger.Log.Infof("chunk name %v", chunkInfo.ChunkName)

	validator := validator.New()
	if err := validator.Struct(chunkInfo); err != nil {
		p.Error400(c, errors.Cause(err))
		return
	}

	if err = p.PackageMgtSvc.ChunkUpload(c, file, chunkInfo); err != nil {
		logger.Log.Errorf("%+v", err)
		p.Error502(c, errors.Cause(err))
		return
	}
	p.Success(c, "ok", nil)
}

/* description:
the api check if all the chunks of package have been uploaded.
If not, return the chunksId which should be upload again.
If all uploaded, merge chunks to package and verify md5 of package.
*/
func (p *packageReceiveApi) PackCheck(c *gin.Context) {
	packName := c.Query("packageName")
	if packName == "" {
		p.Error400(c, constant.ErrParamInvalid)
		return
	}
	chunkLackDto, err := p.PackageMgtSvc.PackCheck(packName)
	if errors.Is(err, service.ErrLackOfChunks) && len(chunkLackDto.LackList) != 0 {
		logger.Log.Infof("%v", err)
		p.Success(c, errors.Cause(err).Error(), *chunkLackDto)
		return
	} else if errors.Is(err, service.ErrPackMd5NotMatch) {
		logger.Log.Infof("%v", err)
		p.Success(c, errors.Cause(err).Error(), nil)
		return

	} else if err != nil {
		logger.Log.Infof("%v", err)
		p.Error502(c, errors.Cause(err))
		return
	}

	p.Success(c, "ok", nil)
}
