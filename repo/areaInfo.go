package repo

import (
	"PackageServer/constant"
	"PackageServer/db"
	"PackageServer/dto"
	"PackageServer/model"
	"net"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type AreaInfoRepo interface {
	Get(name string) (model.AreaInfo, error)
	GetById(id uint) (model.AreaInfo, error)
	Exist(name string) (int64, error)
	List() ([]model.AreaInfo, error)
	Search(name string) ([]model.AreaInfo, error)
	Add(item model.AreaInfo) error
	Update(name string, updation dto.AreaInfoUpdate) (*model.AreaInfo, error)
	Page(nums int, size int) (int64, []model.AreaInfo, error)
	Delete(areaName string) error

	GetWithServices(areaName string) (model.AreaInfo, error)
	ListServicesOf(name string) ([]model.ServiceOnline, error)
	AddServices(name string, serviceOnlines []model.ServiceOnline) error
	DelServices(name string, serviceOnlines []model.ServiceOnline) error
}

func NewAreaInfoRepo() AreaInfoRepo {
	return &areaInfoRepo{}
}

type areaInfoRepo struct{}

func (a *areaInfoRepo) Get(areaName string) (model.AreaInfo, error) {
	var areaInfo model.AreaInfo
	if err := db.Db.Where("area_name = ?", areaName).First(&areaInfo).Error; err != nil {
		return areaInfo, errors.Wrap(err, "repo:areaInfo")
	}
	return areaInfo, nil
}

func (a *areaInfoRepo) GetById(id uint) (model.AreaInfo, error) {
	var areaInfo model.AreaInfo
	if err := db.Db.First(&areaInfo, id).Error; err != nil {
		return areaInfo, errors.Wrap(err, "repo:areaInfo")
	}
	return areaInfo, nil
}

func (a *areaInfoRepo) Exist(areaName string) (int64, error) {
	var total int64
	if err := db.Db.Model(&model.AreaInfo{}).Where("area_name = ?", areaName).Count(&total).Error; err != nil {
		return 0, errors.Wrap(err, "repo:areaInfo")
	}
	return total, nil

}

func (a *areaInfoRepo) List() ([]model.AreaInfo, error) {
	var areaInfos []model.AreaInfo
	if err := db.Db.Find(&areaInfos).Error; err != nil {
		return areaInfos, errors.Wrap(err, "repo:areaInfo")
	}
	return areaInfos, nil
}

func (a *areaInfoRepo) Search(name string) ([]model.AreaInfo, error) {
	var areaInfo []model.AreaInfo
	if err := db.Db.Where("area_name like ?", "%"+name+"%").Find(&areaInfo).Error; err != nil {
		return areaInfo, errors.Wrap(err, "repo:areaInfo")
	}
	return areaInfo, nil
}

func (a *areaInfoRepo) Add(item model.AreaInfo) error {
	if err := db.Db.Create(&item).Error; err != nil {
		return errors.Wrap(err, "repo:areaInfo")
	}
	return nil
}

func (a *areaInfoRepo) Update(name string, updation dto.AreaInfoUpdate) (*model.AreaInfo, error) {
	var mo model.AreaInfo
	if err := db.Db.Where("area_name = ? ", name).First(&mo).Error; err != nil {
		return nil, err
	}
	mo.IpAddr = updation.IpAddr
	mo.Port = updation.Port
	mo.UrlPath = updation.UrlPath

	ipCheck := net.ParseIP(updation.IpAddr)
	if ipCheck == nil {
		return nil, errors.Wrap(constant.ErrIpInvalid, "repo:areaInfoRepo:Update")
	}

	if err := db.Db.Save(&mo).Error; err != nil {
		return nil, err
	}
	return &mo, nil
}

func (a *areaInfoRepo) Page(num int, size int) (int64, []model.AreaInfo, error) {
	var total int64
	var areaInfos []model.AreaInfo
	err := db.Db.Model(&model.AreaInfo{}).Count(&total).Order("area_name").Offset((num - 1) * size).Limit(size).Find(&areaInfos).Error
	return total, areaInfos, errors.Wrap(err, "repo:areaInfo")
}

func (a *areaInfoRepo) Delete(areaName string) error {
	areaInfo, err := a.Get(areaName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrap(constant.ErrRecordNotfound, "repo:serviceOnline")
		}
		return errors.Wrap(err, "repo:areaInfo")
	}
	err = db.Db.Delete(&areaInfo).Error
	return errors.Wrap(err, "repo:areaInfo")
}

func (a *areaInfoRepo) GetWithServices(areaName string) (model.AreaInfo, error) {
	var areaInfo model.AreaInfo
	if err := db.Db.Debug().Preload("ServiceOnlines").Where("area_name = ?", areaName).First(&areaInfo).Error; err != nil {
		return areaInfo, errors.Wrap(err, "repo:areaInfo")
	}
	return areaInfo, nil
}

func (a *areaInfoRepo) ListServicesOf(name string) ([]model.ServiceOnline, error) {
	areaInfo, err := a.Get(name)
	if err != nil {
		return nil, err
	}

	var mos []model.ServiceOnline
	if err := db.Db.Debug().Model(&areaInfo).Association("ServiceOnlines").Find(&mos); err != nil {
		return mos, errors.Wrap(err, "repo:areaInfo")
	}
	return mos, nil
}

func (a *areaInfoRepo) AddServices(name string, serviceOnlines []model.ServiceOnline) error {
	areaInfo, err := a.Get(name)
	if err != nil {
		return err
	}
	if err = db.Db.Debug().Model(&areaInfo).Association("ServiceOnlines").Append(serviceOnlines); err != nil {
		return errors.Wrap(err, "repo:areaInfo")
	}
	return nil
}

func (a *areaInfoRepo) DelServices(name string, serviceOnlines []model.ServiceOnline) error {
	areaInfo, err := a.Get(name)
	if err != nil {
		return err
	}
	if err = db.Db.Debug().Model(&areaInfo).Association("ServiceOnlines").Delete(serviceOnlines); err != nil {
		return errors.Wrap(err, "repo:areaInfo")
	}
	return nil
}
