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

type ServiceRepository struct {
	http *http.Client
}

var _ interfaces.IServiceRepository = (*ServiceRepository)(nil)

func NewServiceRepository() *ServiceRepository {
	return &ServiceRepository{http: http.DefaultClient}
}

func (s *ServiceRepository) GetByID(ctx context.Context, id string) (*response.ServiceResponse, error) {
	logger := logs.Logger()
	path := fmt.Sprintf(SERVICE_ID_ENDPOINT, id)
	resp, err := http.Get(path)
	if err != nil {
		logger.Error().Err(err).Msg("failed to find vehicle by plate")
		return nil, err
	}
	defer resp.Body.Close()

	var service response.ServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&service); err != nil {
		logger.Error().Err(err).Msg("failed to decode vehicle")
		return nil, err
	}
	return &service, nil
}