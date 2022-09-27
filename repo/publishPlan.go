package repo

import (
	"PackageServer/constant"
	"PackageServer/db"
	"PackageServer/dto"
	"PackageServer/model"

	"github.com/pkg/errors"
)

type PublishPlanRepo interface {
	Get(dto dto.PublishPlanLog) (*model.PublishPlanLog, error)
	GetById(id uint) (*model.PublishPlanLog, error)
	Count(dto dto.PublishPlanLog) (int64, error)
	List() ([]model.PublishPlanLog, error)
	ListByAreaId(areaId uint) ([]model.PublishPlanLog, error)
	Add(dto dto.PublishPlanLog) (id uint, err error)
	Save(mo *model.PublishPlanLog) error
	Page(nums int, size int) (int64, []model.PublishPlanLog, error)
}

func NewPublishPlanRepo() PublishPlanRepo {
	return &publishPlanRepo{}
}

type publishPlanRepo struct{}

func (p *publishPlanRepo) Get(dto dto.PublishPlanLog) (*model.PublishPlanLog, error) {
	var mo model.PublishPlanLog
	if err := db.Db.Where("area_info_id = ? AND service_online_id = ? AND version = ?", dto.AreaInfoId, dto.ServiceOnlineId, dto.Version).Last(&mo).Error; err != nil {
		return nil, errors.Wrap(err, "repo:PublishPlanRepo")
	}
	return &mo, nil
}

func (p *publishPlanRepo) GetById(id uint) (*model.PublishPlanLog, error) {
	var mo model.PublishPlanLog
	if err := db.Db.First(&mo, id).Error; err != nil {
		return nil, errors.Wrap(err, "repo:PublishPlanRepo")
	}
	return &mo, nil
}

func (p *publishPlanRepo) Count(dto dto.PublishPlanLog) (int64, error) {
	var total int64
	if err := db.Db.Model(&model.PublishPlanLog{}).Where("area_info_id = ? AND service_online_id = ?", dto.AreaInfoId, dto.ServiceOnlineId).Count(&total).Error; err != nil {
		return total, errors.Wrap(err, "repo:PublishPlanRepo")
	}
	return total, nil
}

func (p *publishPlanRepo) List() ([]model.PublishPlanLog, error) {
	var mos []model.PublishPlanLog
	if err := db.Db.Find(&mos).Error; err != nil {
		return nil, errors.Wrap(err, "repo:PublishPlan")
	}
	return mos, nil
}

func (p *publishPlanRepo) ListByAreaId(areaId uint) ([]model.PublishPlanLog, error) {
	var mos []model.PublishPlanLog
	if err := db.Db.Where("area_info_id = ?", areaId).Find(&mos).Error; err != nil {
		return nil, errors.Wrap(err, "repo:PublishPlan")
	}
	return mos, nil
}

func (p *publishPlanRepo) Add(dto dto.PublishPlanLog) (id uint, err error) {
	var mo model.PublishPlanLog
	mo = model.PublishPlanLog{
		AreaInfoId:      dto.AreaInfoId,
		ServiceOnlineId: dto.ServiceOnlineId,
		Version:         dto.Version,
		PublishStage:    constant.PlanCreated,
		IsComplete:      false,
	}
	if err := db.Db.Create(&mo).Error; err != nil {
		return 0, errors.Wrap(err, "repo:PublishPlan")
	}
	return mo.ID, nil
}

func (p *publishPlanRepo) Save(mo *model.PublishPlanLog) error {
	if err := db.Db.Save(mo).Error; err != nil {
		return errors.Wrap(err, "repo:PublishPlan")
	}
	return nil
}

func (p *publishPlanRepo) Page(nums int, size int) (int64, []model.PublishPlanLog, error) {
	return 0, nil, nil
}
