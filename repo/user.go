package repo

import "PackageServer/model"

type UserRepo interface {
	Get(name string) (*model.User, error)
	Exist(name string) error
	List() ([]model.User, error)
	Save(item model.User) error
	Page(nums int, size int) ([]model.User, error)
	Delete() error
}

func NewUserRepo() UserRepo {
	return &userRepo{}
}

type userRepo struct{}

func (u *userRepo) Get(name string) (*model.User, error) {
	return nil, nil
}

func (u *userRepo) Exist(name string) error {
	return nil
}

func (u *userRepo) List() ([]model.User, error) {
	return nil, nil
}

func (u *userRepo) Save(item model.User) error {
	return nil
}

func (u *userRepo) Page(nums int, size int) ([]model.User, error) {
	return nil, nil
}

func (u *userRepo) Delete() error {
	return nil
}
