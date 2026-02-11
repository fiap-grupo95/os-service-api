package interfaces

import (
	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
)

// VehicleGateway defines the contracts required by the use case layer
type IVehicleGateway interface {
	FindAll() ([]entities.Vehicle, error)
	FindByID(id uint) (*entities.Vehicle, error)
	FindByPlate(plate valueobject.Plate) (*entities.Vehicle, error)
	FindByCustomerID(customerID uint) ([]entities.Vehicle, error)
}

// VehicleRepository defines the contracts required by the adapter layer
type IVehicleRepository interface {
	FindAll() ([]response.VehicleResponse, error)
	FindByID(id uint) (*response.VehicleResponse, error)
	FindByPlate(plate valueobject.Plate) (*response.VehicleResponse, error)
	FindByCustomerID(customerID uint) ([]response.VehicleResponse, error)
}
