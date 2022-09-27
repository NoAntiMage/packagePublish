package service

import (
	"PackageServer/dto"
	"PackageServer/model"
	"PackageServer/repo"
)

type ServiceOnlineSvc interface {
	Get(name string) (*dto.ServiceOnline, error)
	Update(name string, updation dto.ServiceOnline) (*dto.ServiceOnline, error)
	Exist(name string) (int64, error)
	List() ([]dto.ServiceOnline, error)
	Search(name string) ([]dto.ServiceOnline, error)
	Add(item dto.ServiceOnline) error
	Page(nums int, size int) (int64, []dto.ServiceOnline, error)
	Delete(name string) error

	ListAreas(name string) ([]dto.AreaInfo, error)
}

type serviceOnlineSvc struct {
	ServiceOnlineRepo repo.ServiceOnlineRepo
}

func NewServiceOnlineSvc() ServiceOnlineSvc {
	return &serviceOnlineSvc{
		ServiceOnlineRepo: repo.NewServiceOnlineRepo(),
	}
}

func (s *serviceOnlineSvc) Get(name string) (*dto.ServiceOnline, error) {
	var mo model.ServiceOnline
	mo, err := s.ServiceOnlineRepo.Get(name)
	if err != nil {
		return nil, err
	}
	d := toServiceOnlineDto(mo)
	return &d, nil
}

func (s *serviceOnlineSvc) Update(name string, updation dto.ServiceOnline) (*dto.ServiceOnline, error) {
	return nil, nil
}

func (s *serviceOnlineSvc) Exist(name string) (int64, error) {
	return s.ServiceOnlineRepo.Exist(name)
}

func (s *serviceOnlineSvc) List() ([]dto.ServiceOnline, error) {
	var dtos []dto.ServiceOnline
	mos, err := s.ServiceOnlineRepo.List()
	if err != nil {
		return nil, err
	}

	for _, mo := range mos {
		dtos = append(dtos, toServiceOnlineDto(mo))
	}
	return dtos, nil
}

func (s *serviceOnlineSvc) Search(name string) ([]dto.ServiceOnline, error) {
	mos, err := s.ServiceOnlineRepo.Search(name)
	if err != nil {
		return nil, err
	}
	var dtos []dto.ServiceOnline
	for _, mo := range mos {
		d := toServiceOnlineDto(mo)
		dtos = append(dtos, d)
	}
	return dtos, nil
}

func (s *serviceOnlineSvc) Add(do dto.ServiceOnline) error {
	return s.ServiceOnlineRepo.Add(do.ServiceOnline)
}

func (s *serviceOnlineSvc) Page(nums int, size int) (int64, []dto.ServiceOnline, error) {
	var dtos []dto.ServiceOnline
	total, mos, err := s.ServiceOnlineRepo.Page(nums, size)

	for _, mo := range mos {
		dtos = append(dtos, toServiceOnlineDto(mo))
	}

	return total, dtos, err
}

func (s *serviceOnlineSvc) Delete(name string) error {
	return s.ServiceOnlineRepo.Delete(name)
}

func (s *serviceOnlineSvc) ListAreas(name string) ([]dto.AreaInfo, error) {
	mos, err := s.ServiceOnlineRepo.ListAreasOf(name)
	if err != nil {
		return nil, err
	}

	var dtos []dto.AreaInfo
	for _, mo := range mos {
		dtos = append(dtos, dto.ToAreaInfoDto(mo))
	}
	return dtos, nil
}

func toServiceOnlineDto(mo model.ServiceOnline) dto.ServiceOnline {
	serviceOnlineDto := dto.ServiceOnline{ServiceOnline: mo}
	return serviceOnlineDto
}
