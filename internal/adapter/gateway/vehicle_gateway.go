package gateway

import (
	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/response"
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
		v := mapVehicleResponseToDomain(&vehicle)
		vehicles = append(vehicles, *v)
	}
	return vehicles, nil
}

func (g *VehicleGateway) FindByID(id uint) (*entities.Vehicle, error) {
	vehicleResponse, err := g.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	vehicle := mapVehicleResponseToDomain(vehicleResponse)
	return vehicle, nil
}
func (g *VehicleGateway) FindByPlate(plate valueobject.Plate) (*entities.Vehicle, error) {
	vehicleResponse, err := g.repo.FindByPlate(plate)
	if err != nil {
		return nil, err
	}
	vehicle := mapVehicleResponseToDomain(vehicleResponse)
	return vehicle, nil
}
func (g *VehicleGateway) FindByCustomerID(customerID uint) ([]entities.Vehicle, error) {
	vehicleResponse, err := g.repo.FindByCustomerID(customerID)
	if err != nil {
		return nil, err
	}
	vehicles := make([]entities.Vehicle, len(vehicleResponse))
	for _, vehicle := range vehicleResponse {
		v := mapVehicleResponseToDomain(&vehicle)
		vehicles = append(vehicles, *v)
	}
	return vehicles, nil
}
// func (g *VehicleGateway) Create(vehicle entities.Vehicle) (*entities.Vehicle, error) {
// 	return nil, nil
// }
// func (g *VehicleGateway) Update(vehicle entities.Vehicle) error {
// 	return nil
// }

func mapVehicleResponseToDomain(resp *response.VehicleResponse) *entities.Vehicle{
	if resp == nil {
		return nil
	}
	return &entities.Vehicle{
		ID:         resp.ID,
		CustomerID: resp.CustomerID,
		Plate:      valueobject.Plate(resp.Plate),
		Model:      resp.Model,
		Year:       resp.Year,
		Brand:      resp.Brand,
	}
}