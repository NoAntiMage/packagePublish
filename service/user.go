package service

import "PackageServer/model"

type UserSvc interface {
	Get(name string) error
	Exist(name string) error
	List() ([]model.User, error)
}
