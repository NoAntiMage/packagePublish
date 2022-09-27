package service

import "PackageServer/repo"

type AreaServiceReSvc interface {
	HasRelation(areaId uint, serviceId uint) error
}

func NewAreaServiceReSvc() AreaServiceReSvc {
	return &areaServiceReSvc{
		AreaServiceRelationRepo: repo.NewAreaServiceRelationRepo(),
	}
}

type areaServiceReSvc struct {
	AreaServiceRelationRepo repo.AreaServiceRelationRepo
}

func (s areaServiceReSvc) HasRelation(areaId uint, serviceId uint) error {
	return s.AreaServiceRelationRepo.Get(areaId, serviceId)
}
