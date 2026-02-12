package interfaces

import (
	"context"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	dto "github.com/fiap-grupo95/os-service-api/internal/infrastructure/database/model"
)

type IServiceOrderGateway interface {
	Create(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error)
	GetByID(ctx context.Context, id uint, isFullData bool) (*entities.ServiceOrder, error)
	Update(ctx context.Context, serviceOrder *entities.ServiceOrder) error
	List(ctx context.Context) ([]*entities.ServiceOrder, error)
	UpdateEstimate(ctx context.Context, id uint, estimate float64) error
	GetPartsSupplyServiceOrder(ctx context.Context, partsSupplyID uint, serviceOrderID uint) (*entities.ServiceOrderPartsSupply, error)
}

type IServiceOrderRepository interface {
	Create(ctx context.Context, serviceOrderDto *dto.ServiceOrderModel) (*dto.ServiceOrderModel, error)
	GetByID(ctx context.Context, id uint) (*dto.ServiceOrderModel, error)
	Update(ctx context.Context, serviceOrderDto *dto.ServiceOrderModel) error
	List(ctx context.Context) ([]*dto.ServiceOrderModel, error)
	UpdateEstimate(ctx context.Context, id uint, estimate float64) error
	GetPartsSupplyServiceOrder(ctx context.Context, partsSupplyID uint, serviceOrderID uint) (*dto.PartsSupplyServiceOrder, error)
}

