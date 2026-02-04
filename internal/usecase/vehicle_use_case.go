package usecase

import (
	"errors"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
	domainrepo "github.com/fiap-grupo95/os-service-api/internal/usecase/interfaces"
)

var (
	ErrVehicleNotFound      = errors.New("vehicle not found")
	ErrInvalidPlateFormat   = errors.New("invalid plate format")
	ErrVehicleAlreadyExists = errors.New("vehicle already exists")
	ErrInvalidID            = errors.New("invalid vehicle ID")
)

var (
	MessageVehicleUpdatedSuccessfully = "Vehicle updated successfully"
	MessageVehicleNotFound            = "Vehicle not found"
	MessageInvalidPlateFormat         = "Invalid plate format"
	MessageErrorUpdatingVehicle       = "Error updating the vehicle"
	MessageErrorSearch                = "Error searching existing vehicles"
)

type VehicleServiceInterface interface {
	GetAllVehicles() ([]entities.Vehicle, error)
	GetVehicleByID(id uint) (*entities.Vehicle, error)
	GetVehicleByPlate(plate string) (*entities.Vehicle, error)
	GetVehiclesByCustomerID(customerID uint) ([]entities.Vehicle, error)
}

type VehicleService struct {
	repo domainrepo.IVehicleGateway
}

func NewVehicleService(repo domainrepo.IVehicleGateway) VehicleServiceInterface {
	return &VehicleService{repo: repo}
}

func (s *VehicleService) GetAllVehicles() ([]entities.Vehicle, error) {
	vehicles, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	return vehicles, nil
}
func (s *VehicleService) GetVehicleByID(id uint) (*entities.Vehicle, error) {
	vehicle, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if vehicle == nil || vehicle.ID == 0 {
		return nil, ErrVehicleNotFound
	}
	return vehicle, nil
}
func (s *VehicleService) GetVehicleByPlate(plate string) (*entities.Vehicle, error) {
	voPlate := valueobject.ParsePlate(plate)
	if !voPlate.IsValidFormat() {
		return nil, ErrInvalidPlateFormat
	}
	vehicle, err := s.repo.FindByPlate(voPlate)
	if err != nil {
		return nil, err
	}
	if vehicle == nil || vehicle.ID == 0 {
		return nil, ErrVehicleNotFound
	}
	return vehicle, nil
}
func (s *VehicleService) GetVehiclesByCustomerID(customerID uint) ([]entities.Vehicle, error) {
	result, err := s.repo.FindByCustomerID(customerID)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return []entities.Vehicle{}, nil
	}
	return result, nil
}