package entity_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/request"
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

func (r *VehicleRepository) FindByID(id uint) (*response.VehicleResponse, error) {
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

func (r *VehicleRepository) FindByCustomerID(customerID uint) ([]response.VehicleResponse, error) {
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

func (r *VehicleRepository) Create(vehicle request.VehicleCreateRequest) (*response.VehicleResponse, error) {
	logger := logs.Logger()
	path := fmt.Sprintf(VEHICLES_ENDPOINT)

	payload, err := json.Marshal(vehicle)
	if err != nil {
		logger.Error().Err(err).Msg("failed to marshal vehicle")
		return nil, err
	}

	body := io.NopCloser(bytes.NewReader(payload))
	resp, err := r.http.Post(path, "application/json", body)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create vehicle")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		logger.Error().Err(err).Msg("failed to create vehicle")
		return nil, fmt.Errorf("failed to create vehicle: status %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var createdVehicle response.VehicleResponse
	if err := json.NewDecoder(resp.Body).Decode(&createdVehicle); err != nil {
		logger.Error().Err(err).Msg("failed to decode vehicle")
		return nil, err
	}

	return &createdVehicle, nil
}

func (r *VehicleRepository) Update(vehicle request.VehicleUpdateRequest) error {
	logger := logs.Logger()
	// We need a valid ID to update the vehicle
	if vehicle.CustomerID == nil {
		logger.Error().Msg("customer ID is required for update")
		return fmt.Errorf("customer ID is required for update")
	}

	path := fmt.Sprintf(VEHICLES_ID_ENDPOINT, *vehicle.CustomerID)

	payload, err := json.Marshal(vehicle)
	if err != nil {
		logger.Error().Err(err).Msg("failed to marshal vehicle")
		return err	
	}

	req, err := http.NewRequest(http.MethodPatch, path, bytes.NewReader(payload))
	if err != nil {
		logger.Error().Err(err).Msg("failed to create request")
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := r.http.Do(req)
	if err != nil {
		logger.Error().Err(err).Msg("failed to send request")
		return err
	}	
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		logger.Error().Err(err).Msg("failed to update vehicle")
		return fmt.Errorf("failed to update vehicle: status %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// Ensure VehicleRepository implements the domain contract.
var _ domainrepo.IVehicleRepository = (*VehicleRepository)(nil)
