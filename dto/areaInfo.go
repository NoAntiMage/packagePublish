package dto

import "PackageServer/model"

type AreaInfo struct {
	model.AreaInfo
}

type AreaInfoVo struct {
	AreaName  string `json:"AreaName"`
	IpAddr    string `json:"IpAddr"`
	Port      string `json:"Port"`
	UrlPath   string `json:"UrlPath"`
	WorkDir   string `json:"WorkDir"`
	UpdateDir string `json:"UpdateDir"`
}

type AreaInfoUpdate struct {
	IpAddr  string `json:"IpAddr" validate:"required"`
	Port    string `json:"Port" validate:"required"`
	UrlPath string `json:"UrlPath" validate:"required"`
}

func ToAreaInfoDto(areaInfo model.AreaInfo) AreaInfo {
	areaInfoDto := AreaInfo{AreaInfo: areaInfo}
	return areaInfoDto
}

func ToAreaInfoVo(d AreaInfo) AreaInfoVo {
	areaInfoVo := AreaInfoVo{
		AreaName:  d.AreaName,
		IpAddr:    d.IpAddr,
		Port:      d.Port,
		UrlPath:   d.UrlPath,
		WorkDir:   d.WorkDir,
		UpdateDir: d.UpdateDir,
	}
	return areaInfoVo
}
