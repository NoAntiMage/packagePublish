package model

import "gorm.io/gorm"

type AreaInfo struct {
	gorm.Model
	AreaName       string           `json:"AreaName" gorm:"unique" validate:"required"`
	IpAddr         string           `json:"IpAddr"`
	Port           string           `json:"Port"`
	UrlPath        string           `json:"UrlPath"`
	WorkDir        string           `json:"WorkDir"`
	UpdateDir      string           `json:"UpdateDir"`
	Password       string           `json:"Password"`
	ServiceOnlines []*ServiceOnline `gorm:"many2many:area_service_relation;"`
}
