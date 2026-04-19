package service

import (
	"fmt"

	"github.com/blazeisclone/vehicle-dms-inventory/inventory/vehicle"
)

type Service struct {
	repo vehicle.VehicleRepo
}

func NewVehicleSvc(repo vehicle.VehicleRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(name, description string) (*vehicle.Vehicle, error) {
	v := &vehicle.Vehicle{Name: name, Description: description}
	if err := s.repo.Create(v); err != nil {
		return nil, fmt.Errorf("create vehicle: %w", err)
	}
	return v, nil
}

func (s *Service) GetAll() ([]vehicle.Vehicle, error) {
	return s.repo.GetAll()
}

func (s *Service) FindByID(id int) (*vehicle.Vehicle, error) {
	return s.repo.FindByID(id)
}

func (s *Service) Update(id int, name, description string) (*vehicle.Vehicle, error) {
	v := &vehicle.Vehicle{Name: name, Description: description}
	if err := s.repo.UpdateByID(id, v); err != nil {
		return nil, fmt.Errorf("update vehicle: %w", err)
	}
	return v, nil
}

func (s *Service) Delete(id int) error {
	return s.repo.DeleteByID(id)
}
