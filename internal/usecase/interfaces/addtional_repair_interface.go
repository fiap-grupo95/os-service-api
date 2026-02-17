package interfaces

import (
	"context"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
)

type IAdditionalRepairGateway interface {
	CreateAdditionalRepair(ctx context.Context, additionalRepair *entities.AdditionalRepair) (*entities.AdditionalRepair, error)
	GetByID(ctx context.Context, id string) (*entities.AdditionalRepair, error)
	AddPartSupplyAndService(ctx context.Context, additionalRepairID string, services []entities.Service, partsSupplies []entities.PartsSupply, newEstimate float64) error
	ReplacePartSupplyAndService(ctx context.Context, additionalRepairID string, services []entities.Service, partsSupplies []entities.PartsSupply, newEstimate float64) error
	GetByServiceOrder(ctx context.Context, serviceOrderId string) ([]entities.AdditionalRepair, error)
	UpdateAdditionalRepair(ctx context.Context, additionalRepair *entities.AdditionalRepair) (*entities.AdditionalRepair, error)
	GetPartsSupplyQuantity(ctx context.Context, partsSupplyID string, additionalRepairID string) (int, error)
}

type IAdditionalRepairRepository interface {
	CreateAdditionalRepair(ctx context.Context, additionalRepair *entities.AdditionalRepair) (*entities.AdditionalRepair, error)
	GetByID(ctx context.Context, id string) (*entities.AdditionalRepair, error)
	AddPartSupplyAndService(ctx context.Context, additionalRepairID string, services []entities.Service, partsSupplies []entities.PartsSupply, newEstimate float64) error
	ReplacePartSupplyAndService(ctx context.Context, additionalRepairID string, services []entities.Service, partsSupplies []entities.PartsSupply, newEstimate float64) error
	GetByServiceOrder(ctx context.Context, serviceOrderId string) ([]entities.AdditionalRepair, error)
	UpdateAdditionalRepair(ctx context.Context, additionalRepair entities.AdditionalRepair) (entities.AdditionalRepair, error)
	GetPartsSupplyQuantity(ctx context.Context, partsSupplyID string, additionalRepairID string) (int, error)
}
