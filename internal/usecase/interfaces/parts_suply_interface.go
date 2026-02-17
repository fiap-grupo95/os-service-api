package interfaces

import (
	"context"

	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/request"
	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
)

type IPartsSupplyGateway interface {
	GetByID(ctx context.Context, id string) (*entities.PartsSupply, error)
	GetByServiceOrderID(ctx context.Context, serviceOrderID string) ([]entities.PartsSupply, error)
	Reserve(ctx context.Context, partsSupply []entities.PartsSupply) error
	Release(ctx context.Context, partsSupply []entities.PartsSupply) error
	WriteOff(ctx context.Context, partsSupply []entities.PartsSupply) error
	AuthorizeReserve(ctx context.Context, partsSupply []entities.PartsSupply) error
}

type IPartsSupplyRepository interface {
	GetByID(ctx context.Context, id string) (*response.PartsSupplyResponse, error)
	GetByServiceOrderID(ctx context.Context, serviceOrderID string) ([]response.PartsSupplyResponse, error)
	Reserve(ctx context.Context, partsSupply []request.PartsSupplyRequest) error
	Release(ctx context.Context, partsSupply []request.PartsSupplyRequest) error
	WriteOff(ctx context.Context, partsSupply []request.PartsSupplyRequest) error
	AuthorizeReserve(ctx context.Context, partsSupply []request.PartsSupplyRequest) error
}
