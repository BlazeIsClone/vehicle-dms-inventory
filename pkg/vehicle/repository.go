package vehicle

type VehicleRepository interface {
	Create(vehicle *Vehicle) error
	GetAll() (*[]Vehicle, error)
	FindByID(id int) (*Vehicle, error)
	UpdateByID(id int, vehicle *Vehicle) error
	DeleteByID(id int) error
}
