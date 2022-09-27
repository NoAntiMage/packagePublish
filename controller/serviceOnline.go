package controller

import (
	"PackageServer/constant"
	"PackageServer/dto"
	"PackageServer/logger"
	"PackageServer/service"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ServiceOnlineApi interface {
	Create(c *gin.Context)
	GetByName(c *gin.Context)
	UpdateByName(c *gin.Context)
	List(c *gin.Context)
	DeleteByName(c *gin.Context)

	ListAreas(c *gin.Context)
}

type serviceOnlineApi struct {
	BaseApi
	serviceOnlineSvc service.ServiceOnlineSvc
}

func NewServiceOnlineApi() ServiceOnlineApi {
	return &serviceOnlineApi{
		serviceOnlineSvc: service.NewServiceOnlineSvc(),
	}
}

func (s *serviceOnlineApi) Create(c *gin.Context) {
	var d dto.ServiceOnline
	err := c.BindJSON(&d)
	if err != nil {
		logger.Log.Errorf("%+v", err)
		s.Error400(c, errors.Cause(err))
		return
	}

	logger.Log.Debugf("serviceOnlineDto: %v", d)

	validation := validator.New()
	err = validation.Struct(d)
	if err != nil {
		logger.Log.Errorf("%+v", err)
		s.Error400(c, errors.Cause(err))
		return
	}

	total, err := s.serviceOnlineSvc.Exist(d.ServiceName)
	if err != nil {
		logger.Log.Errorf("%+v", err)
		s.Error502(c, errors.Cause(err))
		return
	}

	if total != 0 {
		logger.Log.Errorf("%+v", constant.ErrRecordExist)
		s.Error502(c, constant.ErrRecordExist)
		return
	}

	err = s.serviceOnlineSvc.Add(d)
	if err != nil {
		logger.Log.Errorf("%+v", err)
		s.Error502(c, errors.Cause(err))
		return
	}
	s.Success(c, "ServiceOnlineCreate", nil)
}

func (s *serviceOnlineApi) GetByName(c *gin.Context) {
	var d *dto.ServiceOnline

	name := c.Param("serviceName")
	d, err := s.serviceOnlineSvc.Get(name)
	if err != nil {
		logger.Log.Errorf("%+v", err)
		s.Error502(c, errors.Cause(err))
		return
	}

	s.Success(c, "ServiceOnlineGet", d)
}

//UNDONE
func (s *serviceOnlineApi) UpdateByName(c *gin.Context) {
	var d dto.ServiceOnline
	name := c.Param("serviceName")

	err := c.BindJSON(&d)
	if err != nil {
		logger.Log.Errorf("%+v", err)
		s.Error400(c, errors.Cause(err))
		return
	}

	validation := validator.New()
	if err := validation.Struct(&d); err != nil {
		logger.Log.Errorf("%+v", err)
		s.Error400(c, errors.Cause(err))
		return
	}

	serviceOnlineDto, err := s.serviceOnlineSvc.Update(name, d)
	if err != nil {
		logger.Log.Errorf("%+v", err)
		s.Error502(c, errors.Cause(err))

		return
	}
	s.Success(c, "ServiceOnlineUpdate", serviceOnlineDto)
}

func (s *serviceOnlineApi) List(c *gin.Context) {
	var dtos []dto.ServiceOnline
	var err error

	substring := c.Query("like")
	if substring != "" {
		dtos, err = s.serviceOnlineSvc.Search(substring)
	} else {
		dtos, err = s.serviceOnlineSvc.List()
	}

	if err != nil {
		logger.Log.Errorf("%+v", err)
		s.Error502(c, errors.Cause(err))
		return
	}

	s.Success(c, "ServiceOnlineList", dtos)
}

func (s *serviceOnlineApi) DeleteByName(c *gin.Context) {
	name := c.Param("serviceName")
	if err := s.serviceOnlineSvc.Delete(name); err != nil {
		logger.Log.Errorf("%+v", err)
		if errors.Is(err, constant.ErrRecordNotfound) {
			s.Error404(c, errors.Cause(err))
		} else {
			s.Error502(c, errors.Cause(err))
		}
		return
	}
	s.Success(c, "ServiceOnlineDelete", nil)
}

func (s *serviceOnlineApi) ListAreas(c *gin.Context) {
	serviceName := c.Param("serviceName")
	if serviceName == "" {
		logger.Log.Errorf("%+v", constant.ErrParamIsNotComplete)
		s.Error400(c, constant.ErrParamIsNotComplete)
		return
	}
	dtos, err := s.serviceOnlineSvc.ListAreas(serviceName)
	if err != nil {
		logger.Log.Errorf("%+v", err)
		s.Error502(c, errors.Cause(err))
		return
	}
	s.Success(c, "serviceServices", dtos)
}
