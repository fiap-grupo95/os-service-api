package gateway

import (
	"context"

	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/interfaces"
)

type ServiceGateway struct {
	repo interfaces.IServiceRepository
}

func NewServiceGateway(repo interfaces.IServiceRepository) *ServiceGateway {
	return &ServiceGateway{repo: repo}
}

func (s *ServiceGateway) GetByID(ctx context.Context, id uint) (*entities.Service, error){
	response, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapServiceResponseToDomain(ctx, response), nil
}

func mapServiceResponseToDomain(ctx context.Context, response *response.ServiceResponse) *entities.Service{
	if response == nil {
		return nil
	}
	return &entities.Service{
		ID: response.ID,
		Name: response.Name,
		Description: response.Description,
		Price: response.Price,
	}
}