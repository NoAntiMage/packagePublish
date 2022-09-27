package controller

type AreaServiceRelationApi interface {
}

type areaServiceRelationApi struct {
	BaseApi
}

func NewAreaServiceRelation() AreaServiceRelationApi {
	return &areaServiceRelationApi{}
}
