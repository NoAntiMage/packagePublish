package repo

import (
	"PackageServer/db"
	"PackageServer/model"

	"github.com/pkg/errors"
)

type AreaServiceRelationRepo interface {
	Get(areaId uint, serviceId uint) error
	List() ([]model.AreaServiceRelation, error)
}

func NewAreaServiceRelationRepo() AreaServiceRelationRepo {
	return &areaServiceRelationRepo{}
}

type areaServiceRelationRepo struct{}

func (r areaServiceRelationRepo) Get(areaId uint, serviceId uint) error {
	var mo model.AreaServiceRelation
	if err := db.Db.Where("area_info_id = ? AND service_online_id = ?", areaId, serviceId).First(&mo).Error; err != nil {
		return errors.Wrapf(err, "repo:AreaServiceReRepo")
	}
	return nil
}

func (r areaServiceRelationRepo) List() ([]model.AreaServiceRelation, error) {
	var mos []model.AreaServiceRelation
	if err := db.Db.Find(&mos).Error; err != nil {
		return nil, errors.Wrapf(err, "repo:AreaServiceReRepo")
	}
	return mos, nil
}
