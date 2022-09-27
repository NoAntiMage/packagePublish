package service

import (
	"PackageServer/constant"
	"PackageServer/db"
	"PackageServer/dto"
	"PackageServer/model"
	"PackageServer/repo"
	"PackageServer/rpc"
	"fmt"
	"net"

	"github.com/pkg/errors"
)

type AreaInfoSvc interface {
	Get(name string) (*dto.AreaInfo, error)
	GetById(id uint) (*dto.AreaInfo, error)
	Search(name string) ([]dto.AreaInfo, error)
	Update(name string, updation dto.AreaInfoUpdate) (*dto.AreaInfo, error)
	Exist(name string) (int64, error)
	List() ([]dto.AreaInfo, error)
	Add(item dto.AreaInfo) error
	Page(nums int, size int) (int64, []dto.AreaInfo, error)
	Delete(areaInfo string) error

	ListServices(areaName string) ([]dto.ServiceOnline, error)
	AddServices(areaName string, list dto.ServiceList) error
	DelServices(areaName string, list dto.ServiceList) error

	RemoteHealthCheck(areaName string) error
	GetUrl(areaName string) (string, error)
}

type areaInfoSvc struct {
	areaInfoRepo repo.AreaInfoRepo
}

func NewAreaInfoSvc() AreaInfoSvc {
	return &areaInfoSvc{
		areaInfoRepo: repo.NewAreaInfoRepo(),
	}
}

func (a *areaInfoSvc) Get(name string) (*dto.AreaInfo, error) {
	mo, err := a.areaInfoRepo.Get(name)
	if err != nil {
		return nil, err
	}
	d := dto.ToAreaInfoDto(mo)
	return &d, nil
}

func (a *areaInfoSvc) GetById(id uint) (*dto.AreaInfo, error) {
	mo, err := a.areaInfoRepo.GetById(id)
	if err != nil {
		return nil, err
	}
	d := dto.ToAreaInfoDto(mo)
	return &d, nil
}

func (a *areaInfoSvc) Search(name string) ([]dto.AreaInfo, error) {
	mos, err := a.areaInfoRepo.Search(name)
	if err != nil {
		return nil, err
	}
	var dtos []dto.AreaInfo
	for _, mo := range mos {
		d := dto.ToAreaInfoDto(mo)
		dtos = append(dtos, d)
	}
	return dtos, nil
}

func (a *areaInfoSvc) Update(name string, updation dto.AreaInfoUpdate) (*dto.AreaInfo, error) {
	mo, err := a.areaInfoRepo.Update(name, updation)
	if err != nil {
		return nil, err
	}
	d := dto.ToAreaInfoDto(*mo)
	return &d, nil
}

func (a *areaInfoSvc) Exist(name string) (int64, error) {
	return a.areaInfoRepo.Exist(name)
}

func (a *areaInfoSvc) List() ([]dto.AreaInfo, error) {
	var areaInfoDtos []dto.AreaInfo

	mos, err := a.areaInfoRepo.List()
	if err != nil {
		return nil, err
	}

	for _, mo := range mos {
		areaInfoDtos = append(areaInfoDtos, dto.ToAreaInfoDto(mo))
	}

	return areaInfoDtos, nil
}

func (a *areaInfoSvc) Add(d dto.AreaInfo) error {
	ipCheck := net.ParseIP(d.IpAddr)
	if ipCheck == nil {
		return errors.Wrap(constant.ErrIpInvalid, "service:areaInfoService:Add")
	}

	return a.areaInfoRepo.Add(d.AreaInfo)
}

func (a *areaInfoSvc) Page(nums int, size int) (int64, []dto.AreaInfo, error) {
	var areaInfoDtos []dto.AreaInfo
	total, mos, err := a.areaInfoRepo.Page(nums, size)

	for _, mo := range mos {
		areaInfoDtos = append(areaInfoDtos, dto.ToAreaInfoDto(mo))
	}

	return total, areaInfoDtos, err
}

func (a *areaInfoSvc) Delete(areaInfo string) error {
	return a.areaInfoRepo.Delete(areaInfo)
}

func (a *areaInfoSvc) ListServices(areaName string) ([]dto.ServiceOnline, error) {
	mos, err := a.areaInfoRepo.ListServicesOf(areaName)
	if err != nil {
		return nil, err
	}
	var dtos []dto.ServiceOnline
	for _, mo := range mos {
		dtos = append(dtos, toServiceOnlineDto(mo))
	}
	return dtos, nil
}

func (a *areaInfoSvc) AddServices(areaName string, list dto.ServiceList) error {
	var mos []model.ServiceOnline
	db.Db.Where("service_name IN ?", list.ServiceList).Find(&mos)
	if len(mos) == 0 {
		return errors.Wrap(constant.ErrParamNotFoundInDb, "service:areaInfo:AddServices")
	}
	return a.areaInfoRepo.AddServices(areaName, mos)
}

func (a *areaInfoSvc) DelServices(areaName string, list dto.ServiceList) error {
	var mos []model.ServiceOnline
	db.Db.Where("service_name IN ?", list.ServiceList).Find(&mos)
	if len(mos) == 0 {
		return errors.Wrap(constant.ErrParamNotFoundInDb, "service:areaInfo:DelServices")
	}
	return a.areaInfoRepo.DelServices(areaName, mos)
}

func (a *areaInfoSvc) RemoteHealthCheck(areaName string) error {
	do, err := a.Get(areaName)
	if err != nil {
		return err
	}

	return rpc.HealthCheck(*do)
}

func (a *areaInfoSvc) GetUrl(areaName string) (string, error) {
	do, err := a.Get(areaName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://%v:%v", do.IpAddr, do.Port), nil
}
