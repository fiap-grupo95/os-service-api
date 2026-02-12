package entity_api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/logs"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/interfaces"
)

type PartsSupplyRepository struct {
	http *http.Client
}

var _ interfaces.IPartsSupplyRepository = (*PartsSupplyRepository)(nil)

func NewPartsSupplyRepository() *PartsSupplyRepository {
	return &PartsSupplyRepository{http: http.DefaultClient}
}

func (s *PartsSupplyRepository) GetByID(ctx context.Context, id uint) (*response.PartsSupplyResponse, error) {
	logger := logs.Logger()
	path := fmt.Sprintf(PARTS_SUPPLY_ID_ENDPOINT, id)
	resp, err := http.Get(path)
	if err != nil {
		logger.Error().Err(err).Msg("failed to find vehicle by plate")
		return nil, err
	}
	defer resp.Body.Close()

	var partsSupply response.PartsSupplyResponse
	if err := json.NewDecoder(resp.Body).Decode(&partsSupply); err != nil {
		logger.Error().Err(err).Msg("failed to decode vehicle")
		return nil, err
	}
	return &partsSupply, nil
}

func (s *PartsSupplyRepository) GetByServiceOrderID(ctx context.Context, serviceOrderID uint) ([]response.PartsSupplyResponse, error) {
	logger := logs.Logger()
	path := fmt.Sprintf(PARTS_SUPPLY_SERVICE_ORDER_ID_ENDPOINT, serviceOrderID)
	resp, err := http.Get(path)
	if err != nil {
		logger.Error().Err(err).Msg("failed to find vehicle by plate")
		return nil, err
	}
	defer resp.Body.Close()

	var partsSupply []response.PartsSupplyResponse
	if err := json.NewDecoder(resp.Body).Decode(&partsSupply); err != nil {
		logger.Error().Err(err).Msg("failed to decode vehicle")
		return nil, err
	}
	return partsSupply, nil
}
