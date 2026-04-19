package service

type CreateVehicleCommand struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateVehicleCommand struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
