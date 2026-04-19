package vehicle

import (
	"fmt"
)

type Service struct {
	repo VehicleRepo
}

func NewVehicleSvc(repo VehicleRepo) *Service {
	return &Service{repo: repo}
}

func (svc *Service) Create(cmd CreateVehicleCommand) (*Vehicle, error) {
	vehicle := &Vehicle{Name: cmd.Name, Description: cmd.Description}

	if err := svc.repo.Create(vehicle); err != nil {
		return nil, fmt.Errorf("create vehicle: %w", err)
	}

	return vehicle, nil
}

func (svc *Service) GetAll() ([]Vehicle, error) {
	return svc.repo.GetAll()
}

func (svc *Service) FindByID(id int) (*Vehicle, error) {
	return svc.repo.FindByID(id)
}

func (svc *Service) Update(id int, cmd UpdateVehicleCommand) (*Vehicle, error) {
	vehicle := &Vehicle{Name: cmd.Name, Description: cmd.Description}

	if err := svc.repo.UpdateByID(id, vehicle); err != nil {
		return nil, fmt.Errorf("update vehicle: %w", err)
	}

	return vehicle, nil
}

func (svc *Service) Delete(id int) error {
	return svc.repo.DeleteByID(id)
}
