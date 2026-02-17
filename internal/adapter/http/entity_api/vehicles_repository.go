package entity_api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/logs"
	domainrepo "github.com/fiap-grupo95/os-service-api/internal/usecase/interfaces"
)

type VehicleRepository struct {
	http *http.Client
}

func NewVehicleRepository() domainrepo.IVehicleRepository {
	return &VehicleRepository{http: http.DefaultClient}
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

// Ensure VehicleRepository implements the domain contract.
var _ domainrepo.IVehicleRepository = (*VehicleRepository)(nil)
