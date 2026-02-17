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

func (g *VehicleGateway) FindByID(id string) (*entities.Vehicle, error) {
	vehicleResponse, err := g.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	vehicle := mapVehicleResponseToDomain(vehicleResponse)
	return vehicle, nil
}


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