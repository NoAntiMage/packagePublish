package service

import (
	"PackageServer/dto"
	"PackageServer/repo"
)

type PlanStageSvc interface {
	Update(updation *dto.PlanStageUpdate) error
}

func NewPlanStageSvc() PlanStageSvc {
	return &planStageSvc{
		publishPlanRepo: repo.NewPublishPlanRepo(),
	}
}

type planStageSvc struct {
	publishPlanRepo repo.PublishPlanRepo
}

func (p *planStageSvc) Update(updation *dto.PlanStageUpdate) error {
	d := dto.PublishPlanLog{
		AreaInfoId:      updation.AreaInfoId,
		ServiceOnlineId: updation.ServiceOnlineId,
		Version:         updation.Version,
	}
	mo, err := p.publishPlanRepo.Get(d)
	if err != nil {
		return err
	}
	mo.PublishStage = updation.Stage
	return p.publishPlanRepo.Save(mo)
}
