package interfaces

import (
	"context"

	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
)

type IPartsSupplyGateway interface {
	GetByID(ctx context.Context, id uint) (*entities.PartsSupply, error)
	GetByServiceOrderID(ctx context.Context, serviceOrderID uint) ([]entities.PartsSupply, error)
}

type IPartsSupplyRepository interface {
	GetByID(ctx context.Context, id uint) (*response.PartsSupplyResponse, error)
	GetByServiceOrderID(ctx context.Context, serviceOrderID uint) ([]response.PartsSupplyResponse, error)
}
