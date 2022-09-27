package dto

import (
	"time"
)

type PublishPlanAdd struct {
	AreaName    string `json:"AreaName" validate:"required"`
	ServiceName string `json:"ServiceName" validate:"required"`
	Version     string `json:"Version" validate:"required"`
	FromDate    time.Time
}

type PublishPlanLog struct {
	AreaInfoId      uint
	ServiceOnlineId uint
	Version         string
	FromDate        time.Time `gorm:"type: datetime"`
}

type PlanStageUpdate struct {
	AreaInfoId      uint
	ServiceOnlineId uint
	Version         string
	Stage           string
}

func PublishPlanLog2Updation(planDto PublishPlanLog) *PlanStageUpdate {
	updation := &PlanStageUpdate{
		AreaInfoId:      planDto.AreaInfoId,
		ServiceOnlineId: planDto.ServiceOnlineId,
		Version:         planDto.Version,
	}
	return updation
}
