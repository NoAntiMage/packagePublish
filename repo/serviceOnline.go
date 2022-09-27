package repo

import (
	"PackageServer/constant"
	"PackageServer/db"
	"PackageServer/model"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type ServiceOnlineRepo interface {
	Get(name string) (model.ServiceOnline, error)
	GetById(id uint) (model.ServiceOnline, error)
	Exist(name string) (int64, error)
	List() ([]model.ServiceOnline, error)
	Search(name string) ([]model.ServiceOnline, error)
	Add(item model.ServiceOnline) error
	Page(nums int, size int) (int64, []model.ServiceOnline, error)
	Delete(name string) error
	GetWithAreas(name string) (model.ServiceOnline, error)
	ListAreasOf(name string) ([]model.AreaInfo, error)
}

func NewServiceOnlineRepo() ServiceOnlineRepo {
	return &serviceOnlineRepo{}
}

type serviceOnlineRepo struct{}

func (s *serviceOnlineRepo) Get(name string) (model.ServiceOnline, error) {
	var serviceOnline model.ServiceOnline
	if err := db.Db.Where("service_name = ?", name).First(&serviceOnline).Error; err != nil {
		return serviceOnline, errors.Wrap(err, "repo:serviceOnline")
	}
	return serviceOnline, nil
}

func (s *serviceOnlineRepo) GetById(id uint) (model.ServiceOnline, error) {
	var serviceOnline model.ServiceOnline
	if err := db.Db.First(&serviceOnline, id).Error; err != nil {
		return serviceOnline, errors.Wrap(err, "repo:serviceOnline")
	}
	return serviceOnline, nil
}

func (s *serviceOnlineRepo) Exist(name string) (int64, error) {
	var total int64
	if err := db.Db.Model(&model.ServiceOnline{}).Where("service_name = ?", name).Count(&total).Error; err != nil {
		return 0, errors.Wrap(err, "repo:serviceOnline:")
	}
	return total, nil
}

func (s *serviceOnlineRepo) List() ([]model.ServiceOnline, error) {
	var serviceOnlines []model.ServiceOnline
	if err := db.Db.Find(&serviceOnlines).Error; err != nil {
		return serviceOnlines, errors.Wrap(err, "repo:serviceOnline")
	}
	return serviceOnlines, nil
}

func (s *serviceOnlineRepo) Search(name string) ([]model.ServiceOnline, error) {
	var serviceOnlines []model.ServiceOnline
	if err := db.Db.Where("service_name like ?", "%"+name+"%").Find(&serviceOnlines).Error; err != nil {
		return serviceOnlines, errors.Wrap(err, "repo: serviceOnline")
	}
	return serviceOnlines, nil
}

func (s *serviceOnlineRepo) Add(item model.ServiceOnline) error {
	if err := db.Db.Create(&item).Error; err != nil {
		return errors.Wrap(err, "repo:serviceOnline")
	}
	return nil
}

func (s *serviceOnlineRepo) Page(num int, size int) (int64, []model.ServiceOnline, error) {
	var total int64
	var serviceOnlines []model.ServiceOnline
	if err := db.Db.Model(&model.ServiceOnline{}).Count(&total).Offset((num - 1) * size).Limit(size).Find(&serviceOnlines).Error; err != nil {
		return 0, nil, errors.Wrap(err, "repo:serviceOnline")
	}
	return total, serviceOnlines, nil
}

func (s *serviceOnlineRepo) Delete(name string) error {
	serviceOnline, err := s.Get(name)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrap(constant.ErrRecordNotfound, "repo:serviceOnline")
		}
		return err
	}
	err = db.Db.Delete(&serviceOnline).Error
	return errors.Wrap(err, "repo:serviceOnline")
}

func (s *serviceOnlineRepo) GetWithAreas(name string) (model.ServiceOnline, error) {
	var serviceOnline model.ServiceOnline
	if err := db.Db.Preload("AreaInfos").Where("service_name = ?", name).First(&serviceOnline).Error; err != nil {
		return serviceOnline, errors.Wrap(err, "repo:serviceOnline")
	}
	return serviceOnline, nil
}

func (s *serviceOnlineRepo) ListAreasOf(name string) ([]model.AreaInfo, error) {
	service, err := s.Get(name)
	if err != nil {
		return nil, err
	}

	var mos []model.AreaInfo
	if err := db.Db.Debug().Model(&service).Association("AreaInfos").Find(&mos); err != nil {
		return mos, errors.Wrap(err, "repo:serviceOnline")
	}
	return mos, nil
}

// NOT ALLOWED in ServiceSvc
// modify in AreaSvc
func (s *serviceOnlineRepo) AddAreas(name string, areaInfos []model.AreaInfo) error {
	return nil
}

func (s *serviceOnlineRepo) DelAreas() error {
	return nil
}
