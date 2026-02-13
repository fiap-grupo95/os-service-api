package interfaces

import (
	"context"

	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/request"
	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
)

type IBillingServiceGateway interface {
	CreateEstimate(ctx context.Context, serviceOrder *entities.ServiceOrder) (*float64, error)
}

type IBillingServiceRepository interface {
	CreateEstimate(ctx context.Context, request *request.EstimateRequest) (*response.EstimateResponse, error)
}
