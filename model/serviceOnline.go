package model

import "gorm.io/gorm"

type ServiceOnline struct {
	gorm.Model
	ServiceName string      `gorm:"unique" json:"ServiceName" validate:"required"`
	ArchiveType string      `gorm:"type:varchar(6)" json:"ArchiveType" validate:"required"`
	AreaInfos   []*AreaInfo `gorm:"many2many:area_service_relation;"`
}
