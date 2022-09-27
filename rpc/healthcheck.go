package rpc

import (
	"PackageServer/dto"
	"PackageServer/request"
)

func HealthCheck(d dto.AreaInfo) error {
	return request.GetPing(d)
}
