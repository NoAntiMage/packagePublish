package model

import (
	"time"

	"gorm.io/gorm"
)

type PublishPlanLog struct {
	gorm.Model
	AreaInfoId      uint
	AreaInfo        AreaInfo
	ServiceOnlineId uint
	ServiceOnline   ServiceOnline
	Version         string
	FromDate        time.Time `gorm:"type: datetime"`
	ToDate          time.Time `gorm:"type: datetime"`
	PublishStage    string
	IsComplete      bool
}
