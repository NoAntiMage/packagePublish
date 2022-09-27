package controller

import (
	"PackageServer/constant"
	"PackageServer/dto"
	"PackageServer/logger"
	"PackageServer/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type PublishPlanApi interface {
	PublishVersion(c *gin.Context)
}

type publishPlanApi struct {
	BaseApi
	areaInfoSvc      service.AreaInfoSvc
	serviceOnlineSvc service.ServiceOnlineSvc
	publishPlanSvc   service.PublishPlanSvc
	areaServiceReSvc service.AreaServiceReSvc
}

func NewPublishPlanApi() PublishPlanApi {
	return &publishPlanApi{
		areaInfoSvc:      service.NewAreaInfoSvc(),
		serviceOnlineSvc: service.NewServiceOnlineSvc(),
		publishPlanSvc:   service.NewPublishPlanSvc(),
		areaServiceReSvc: service.NewAreaServiceReSvc(),
	}
}

func (p *publishPlanApi) PublishVersion(c *gin.Context) {
	var vo dto.PublishPlanAdd
	var do dto.PublishPlanLog

	if err := c.BindJSON(&vo); err != nil {
		logger.Log.Errorf("%+v ", err)
		p.Error400(c, err)
		return
	}
	validator := validator.New()
	if err := validator.Struct(vo); err != nil {
		logger.Log.Errorf("%+v ", err)
		p.Error400(c, constant.ErrParamInvalid)
		return
	}

	areaInfoDto, err := p.areaInfoSvc.Get(vo.AreaName)
	if err != nil {
		logger.Log.Errorf("%+v ", err)
		p.Error502(c, errors.Cause(err))
		return
	}
	do.AreaInfoId = areaInfoDto.ID

	serviceDto, err := p.serviceOnlineSvc.Get(vo.ServiceName)
	if err != nil {
		logger.Log.Errorf("%+v ", err)
		p.Error502(c, errors.Cause(err))
		return
	}
	do.ServiceOnlineId = serviceDto.ID

	if err := p.areaServiceReSvc.HasRelation(do.AreaInfoId, do.ServiceOnlineId); err != nil {
		logger.Log.Errorf("%+v ", err)
		p.Error502(c, errors.Cause(err))
		return
	}

	do.Version = vo.Version
	logger.Log.Debugf("get publishPlan: %v", do)
	if err := p.publishPlanSvc.PublishVersion(do); err != nil {
		logger.Log.Errorf("%+v ", err)
		p.Error502(c, errors.Cause(err))
		return
	}
	p.Success(c, "ok", nil)
}

//TODO below 2022.08.23
func (p *publishPlanApi) PublishJobStatus(c *gin.Context) {}

func (p *publishPlanApi) PublishJobStop(c *gin.Context) {}

func (p *publishPlanApi) PublishJobContinue(c *gin.Context) {}

func (p *publishPlanApi) PublishJobRetry(c *gin.Context) {}
