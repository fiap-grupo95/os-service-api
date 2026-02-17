package interfaces

import (
	"context"

	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"

)

type IServiceGateway interface {
	GetByID(ctx context.Context, id string) (*entities.Service, error)
}

type IServiceRepository interface {
	GetByID(ctx context.Context, id string) (*response.ServiceResponse, error)
}
