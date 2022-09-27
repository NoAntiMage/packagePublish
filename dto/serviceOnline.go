package dto

import "PackageServer/model"

type ServiceOnline struct {
	model.ServiceOnline
}

type ServiceOnlineUpdation struct {
}

type ServiceList struct {
	ServiceList []string `json:"ServiceList" validate:"required"`
}
