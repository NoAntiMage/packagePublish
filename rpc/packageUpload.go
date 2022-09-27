package rpc

import (
	"PackageServer/cache"
	"PackageServer/constant"
	"PackageServer/dto"
	"PackageServer/request"
)

type PackagePublishRpc interface {
	PostPackageInfo(remoteUrl string, info *dto.PackageInfoPost) error
	PostChunkUpload(remoteUrl string, chunkLocation string) error
	GetPackCheck(remoteUrl string, PackName string) error
}

func NewPackagePublishRpc(area string) (PackagePublishRpc, error) {
	obj := &packagePublishRpc{
		authRpc: NewAuthRpc(),
	}
	token, err := obj.authRpc.GetTokenByEnv(area)
	if err != nil {
		return nil, err
	}
	obj.token = token
	return obj, nil
}

type packagePublishRpc struct {
	authRpc  AuthRpc
	token    string
	jwtCache cache.JsonWebTokenCache
}

func (p *packagePublishRpc) PostPackageInfo(remoteUrl string, info *dto.PackageInfoPost) error {
	var reqInfoDto = dto.RequestInfo{
		TargetUrl: remoteUrl + p.authRpc.UrlPath + constant.UrlPackInfo,
		Token:     p.token,
	}
	return request.PostPackInfo(reqInfoDto, info)
}

/*
NEXT VERSION:
concurrent, breaker, retry  will be written here.
*/
func (p *packagePublishRpc) PostChunkUpload(remoteUrl string, chunkLocation string) error {
	var reqInfoDto = dto.RequestInfo{
		TargetUrl: remoteUrl + p.authRpc.UrlPath + constant.UrlChunkUpload,
		Token:     p.token,
	}
	return request.PostChunkUpload(reqInfoDto, chunkLocation)
}

func (p *packagePublishRpc) GetPackCheck(remoteUrl string, PackName string) error {
	var reqInfoDto = dto.RequestInfo{
		TargetUrl: remoteUrl + p.authRpc.UrlPath + constant.UrlPackCheck,
		Token:     p.token,
	}
	return request.GetPackCheck(reqInfoDto, PackName)
}
