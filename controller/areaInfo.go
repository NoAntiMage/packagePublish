package controller

import (
	"PackageServer/constant"
	"PackageServer/dto"
	"PackageServer/logger"
	"PackageServer/service"
	"fmt"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AreaInfoApi interface {
	Create(c *gin.Context)
	GetByName(c *gin.Context)
	UpdateByName(c *gin.Context)
	List(c *gin.Context)
	DeleteByName(c *gin.Context)

	ListServices(c *gin.Context)
	AddServices(c *gin.Context)
	DelServices(c *gin.Context)
}

type areaInfoApi struct {
	BaseApi
	areaInfoSvc service.AreaInfoSvc
}

func NewAreaInfoApi() AreaInfoApi {
	return &areaInfoApi{
		areaInfoSvc: service.NewAreaInfoSvc(),
	}
}

func (a *areaInfoApi) Create(c *gin.Context) {
	var d dto.AreaInfo
	err := c.BindJSON(&d)
	if err != nil {
		logger.Log.Errorf("%+v", err)
		a.Error400(c, errors.Cause(err))
		return
	}

	logger.Log.Debugf("areaInfoDto: %v", d)

	validation := validator.New()
	err = validation.Struct(d)
	if err != nil {
		logger.Log.Errorf("%+v", err)
		a.Error400(c, errors.Cause(err))
		return
	}

	total, err := a.areaInfoSvc.Exist(d.AreaName)
	if err != nil {
		logger.Log.Errorf("%+v", err)
		a.Error502(c, errors.Cause(err))
		return
	}

	if total != 0 {
		logger.Log.Errorf("%+v", constant.ErrRecordExist)
		a.Error502(c, constant.ErrRecordExist)
		return
	}

	err = a.areaInfoSvc.Add(d)
	if err != nil {
		logger.Log.Errorf("%+v", err)
		a.Error502(c, errors.Cause(err))
		return
	}

	a.Success(c, "AreaInfoCreate", nil)
}

func (a *areaInfoApi) GetByName(c *gin.Context) {
	name := c.Param("areaName")
	d, err := a.areaInfoSvc.Get(name)
	if err != nil {
		logger.Log.Errorf("%+v", err)
		a.Error502(c, errors.Cause(err))
		return
	}

	vo := dto.ToAreaInfoVo(*d)
	a.Success(c, "AreaInfoGet", vo)
}

func (a *areaInfoApi) UpdateByName(c *gin.Context) {
	var d dto.AreaInfoUpdate
	name := c.Param("areaName")

	err := c.BindJSON(&d)
	if err != nil {
		logger.Log.Errorf("%+v", err)
		a.Error400(c, errors.Cause(err))
		return
	}

	validation := validator.New()
	if err := validation.Struct(&d); err != nil {
		logger.Log.Errorf("%+v", err)
		a.Error400(c, errors.Cause(err))
		return
	}

	areaInfoDto, err := a.areaInfoSvc.Update(name, d)
	if err != nil {
		logger.Log.Errorf("%+v", err)
		a.Error502(c, errors.Cause(err))
		return
	}

	vo := dto.ToAreaInfoVo(*areaInfoDto)
	a.Success(c, "AreaInfoUpdate", vo)
}

func (a *areaInfoApi) List(c *gin.Context) {
	var dtos []dto.AreaInfo
	var vos []dto.AreaInfoVo
	var err error

	substring := c.Query("like")
	if substring != "" {
		dtos, err = a.areaInfoSvc.Search(substring)
	} else {
		dtos, err = a.areaInfoSvc.List()
	}

	if err != nil {
		logger.Log.Errorf("%+v", err)
		a.Error502(c, errors.Cause(err))
		return
	}

	for _, d := range dtos {
		vos = append(vos, dto.ToAreaInfoVo(d))
	}

	a.Success(c, "AreaInfoList", vos)
}

func (a *areaInfoApi) DeleteByName(c *gin.Context) {
	name := c.Param("areaName")
	if err := a.areaInfoSvc.Delete(name); err != nil {
		logger.Log.Errorf("%+v", err)
		a.Error502(c, errors.Cause(err))
		return
	}
	a.Success(c, "AreaInfoDelete", nil)
}

func (a *areaInfoApi) ListServices(c *gin.Context) {
	areaName := c.Param("areaName")
	if areaName == "" {
		logger.Log.Errorf("%+v", constant.ErrParamInvalid)
		a.Error400(c, constant.ErrParamInvalid)
		return
	}
	dtos, err := a.areaInfoSvc.ListServices(areaName)
	if err != nil {
		logger.Log.Errorf("%+v", err)
		a.Error502(c, errors.Cause(err))
		return
	}
	a.Success(c, "ListServices", dtos)
}

func (a *areaInfoApi) AddServices(c *gin.Context) {
	areaName := c.Param("areaName")
	fmt.Printf("areaName: %v", areaName)
	if areaName == "" {
		logger.Log.Errorf("%+v", constant.ErrParamInvalid)
		a.Error400(c, constant.ErrParamInvalid)
		return
	}

	var serviceList dto.ServiceList
	if err := c.BindJSON(&serviceList); err != nil {
		logger.Log.Errorf("%+v", err)
		a.Error400(c, err)
		return
	}

	logger.Log.Debugf("area %v serviceList %v", areaName, serviceList.ServiceList)

	validation := validator.New()
	if err := validation.Struct(serviceList); err != nil {
		logger.Log.Errorf("%+v", err)
		a.Error400(c, constant.ErrParamInvalid)
		return
	}

	if err := a.areaInfoSvc.AddServices(areaName, serviceList); err != nil {
		logger.Log.Errorf("%+v", err)
		a.Error502(c, errors.Cause(err))
		return
	}
	a.Success(c, "AddServices", nil)
}

func (a *areaInfoApi) DelServices(c *gin.Context) {
	areaName := c.Param("areaName")
	if areaName == "" {
		logger.Log.Errorf("%+v", constant.ErrParamInvalid)
		a.Error400(c, constant.ErrParamInvalid)
		return
	}

	var serviceList dto.ServiceList
	if err := c.BindJSON(&serviceList); err != nil {
		logger.Log.Errorf("%+v", err)
		a.Error400(c, err)
		return
	}

	logger.Log.Debugf("area %v serviceList %v", areaName, serviceList.ServiceList)

	validation := validator.New()
	if err := validation.Struct(serviceList); err != nil {
		logger.Log.Errorf("%+v", err)
		a.Error400(c, constant.ErrParamInvalid)
		return
	}

	if err := a.areaInfoSvc.DelServices(areaName, serviceList); err != nil {
		logger.Log.Errorf("%+v", err)
		a.Error502(c, errors.Cause(err))
		return
	}
	a.Success(c, "DelServices", nil)
}
