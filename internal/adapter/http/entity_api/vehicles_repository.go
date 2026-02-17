package entity_api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/logs"
	domainrepo "github.com/fiap-grupo95/os-service-api/internal/usecase/interfaces"

	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
)

type VehicleRepository struct {
	http *http.Client
}

func NewVehicleRepository() domainrepo.IVehicleRepository {
	return &VehicleRepository{http: http.DefaultClient}
}

func (r *VehicleRepository) FindAll() ([]response.VehicleResponse, error) {
	logger := logs.Logger()
	resp, err := http.Get(VEHICLES_ENDPOINT)
	if err != nil {
		logger.Error().Err(err).Msg("failed to find all vehicles")
		return nil, err
	}
	defer resp.Body.Close()

	var vehicles []response.VehicleResponse
	if err := json.NewDecoder(resp.Body).Decode(&vehicles); err != nil {
		logger.Error().Err(err).Msg("failed to decode vehicles")
		return nil, err
	}
	return vehicles, nil
}

func (r *VehicleRepository) FindByID(id string) (*response.VehicleResponse, error) {
	logger := logs.Logger()
	path := fmt.Sprintf(VEHICLES_ID_ENDPOINT, id)
	resp, err := http.Get(path)
	if err != nil {
		logger.Error().Err(err).Msg("failed to find vehicle by ID")
		return nil, err
	}
	defer resp.Body.Close()

	var vehicle response.VehicleResponse
	if err := json.NewDecoder(resp.Body).Decode(&vehicle); err != nil {
		logger.Error().Err(err).Msg("failed to decode vehicle")
		return nil, err
	}
	return &vehicle, nil
}

func (r *VehicleRepository) FindByPlate(plate valueobject.Plate) (*response.VehicleResponse, error) {
	logger := logs.Logger()
	path := fmt.Sprintf(VEHICLES_PLATE_ENDPOINT, plate)
	resp, err := http.Get(path)
	if err != nil {
		logger.Error().Err(err).Msg("failed to find vehicle by plate")
		return nil, err
	}
	defer resp.Body.Close()

	var vehicle response.VehicleResponse
	if err := json.NewDecoder(resp.Body).Decode(&vehicle); err != nil {
		logger.Error().Err(err).Msg("failed to decode vehicle")
		return nil, err
	}
	return &vehicle, nil
}

func (r *VehicleRepository) FindByCustomerID(customerID string) ([]response.VehicleResponse, error) {
	logger := logs.Logger()
	path := fmt.Sprintf(VEHICLES_CUSTOMER_ID_ENDPOINT, customerID)
	resp, err := http.Get(path)
	if err != nil {
		logger.Error().Err(err).Msg("failed to find vehicles by customer ID")
		return nil, err
	}
	defer resp.Body.Close()

	var vehicles []response.VehicleResponse
	if err := json.NewDecoder(resp.Body).Decode(&vehicles); err != nil {
		logger.Error().Err(err).Msg("failed to decode vehicles")
		return nil, err
	}
	return vehicles, nil
}

// Ensure VehicleRepository implements the domain contract.
var _ domainrepo.IVehicleRepository = (*VehicleRepository)(nil)
