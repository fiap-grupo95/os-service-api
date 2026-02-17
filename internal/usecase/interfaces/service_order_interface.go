package interfaces

import (
	"context"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	dto "github.com/fiap-grupo95/os-service-api/internal/infrastructure/database/model"
)

type IServiceOrderGateway interface {
	Create(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error)
	GetByID(ctx context.Context, id string, isFullData bool) (*entities.ServiceOrder, error)
	Update(ctx context.Context, serviceOrder *entities.ServiceOrder) error
	List(ctx context.Context) ([]*entities.ServiceOrder, error)
	GetPartsSupplyServiceOrder(ctx context.Context, partsSupplyID string, serviceOrderID string) (*entities.ServiceOrderPartsSupply, error)
}

type IServiceOrderRepository interface {
	Create(ctx context.Context, serviceOrderDto *dto.ServiceOrderModel) (*dto.ServiceOrderModel, error)
	GetByID(ctx context.Context, id string) (*dto.ServiceOrderModel, error)
	Update(ctx context.Context, serviceOrderDto *dto.ServiceOrderModel) error
	List(ctx context.Context) ([]*dto.ServiceOrderModel, error)
	GetPartsSupplyServiceOrder(ctx context.Context, partsSupplyID string, serviceOrderID string) (*dto.PartsSupplyServiceOrder, error)
}
