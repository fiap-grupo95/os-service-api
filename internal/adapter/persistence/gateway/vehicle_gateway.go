package gateway

import (
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/interfaces"
)

type VehicleGateway struct {
	repo interfaces.IVehicleRepository
}

func NewVehicleGateway(repo interfaces.IVehicleRepository) *VehicleGateway {
	return &VehicleGateway{
		repo: repo,
	}
}


func (g *VehicleGateway) FindAll() ([]entities.Vehicle, error) {
	vehiclesResponse, err := g.repo.FindAll()
	if err != nil {
		return nil, err
	}
	vehicles := make([]entities.Vehicle, len(vehiclesResponse))
	for _, vehicle := range vehiclesResponse {
		v := entities.Vehicle{
			ID: vehicle.ID,
			CustomerID: vehicle.CustomerID,
			Plate: valueobject.Plate(vehicle.Plate),
			Model: vehicle.Model,
			Year: vehicle.Year,
			Brand: vehicle.Brand,
		}
		vehicles = append(vehicles, v)
	}
	return vehicles, nil
}

func (g *VehicleGateway) FindByID(id uint) (*entities.Vehicle, error) {
	return nil, nil
}
func (g *VehicleGateway) FindByPlate(plate valueobject.Plate) (*entities.Vehicle, error){
	return nil, nil
}
func (g *VehicleGateway) FindByCustomerID(customerID uint) ([]entities.Vehicle, error){
	return nil, nil
}
func (g *VehicleGateway) Create(vehicle entities.Vehicle) (*entities.Vehicle, error){
	return nil, nil
}
func (g *VehicleGateway) Update(vehicle entities.Vehicle) error{
	return nil
}