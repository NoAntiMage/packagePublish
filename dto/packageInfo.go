package dto

import "strings"

type PackageInfo struct {
	FileName string `json:"Name" validate:"required"`
	Path     string `json:"Path" validate:"required"`
	Md5Sum   string `json:"Md5" validate:"required"`
}

type PackagePushJob struct {
	Id        uint   `validate:"required"`
	FileName  string `validate:"required"`
	Path      string `validate:"required"`
	Md5Sum    string `validate:"required"`
	Area      string `validate:"required"`
	RemoteUrl string `validate:"required"`
}

type PackageInfoPost struct {
	FileName string `json:"FileName" validate:"required"`
	Md5Sum   string `json:"Md5" validate:"required"`
	ChunkNum int    `json:"ChunkNum" validate:"required"`
}

type PackageInfoEnv struct {
	FileName    string `validate:"required"`
	ServiceName string `validate:"required"`
	AreaName    string `validate:"required"`
	Version     string `validate:"required"`
	Md5Sum      string `validate:"required"`
	ChunkNum    int    `validate:"required"`
}

type ChunkInfo struct {
	ChunkName string `validate:"required"`
	FileName  string `validate:"required"`
}

type ChunkLack struct {
	FileName string `json:"FileName" validate:"required"`
	LackList []int  `json:"LackList"  validate:"required"`
}

func ChunkLocToChunkInfo(chunkLocation string) *ChunkInfo {
	var do ChunkInfo
	l1 := strings.Split(chunkLocation, "/")
	do.ChunkName = l1[len(l1)-1]

	l2 := strings.Split(do.ChunkName, ".")
	do.FileName = strings.Join(l2[0:2], ".")
	return &do
}
