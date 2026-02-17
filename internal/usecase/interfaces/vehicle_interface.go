package interfaces

import (
	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
)

// VehicleGateway defines the contracts required by the use case layer
type IVehicleGateway interface {
	FindByID(id string) (*entities.Vehicle, error)
}

// VehicleRepository defines the contracts required by the adapter layer
type IVehicleRepository interface {
	FindByID(id string) (*response.VehicleResponse, error)
}
