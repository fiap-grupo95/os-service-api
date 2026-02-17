package interfaces

import (
	"context"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
)

type IAdditionalRepairGateway interface {
	CreateAdditionalRepair(ctx context.Context, additionalRepair *entities.AdditionalRepair) (*entities.AdditionalRepair, error)
	GetByID(ctx context.Context, id string) (*entities.AdditionalRepair, error)
	GetByServiceOrderID(ctx context.Context, serviceOrderID string) ([]entities.AdditionalRepair, error)
	UpdateAdditionalRepair(ctx context.Context, additionalRepair *entities.AdditionalRepair) (*entities.AdditionalRepair, error)
}

type IAdditionalRepairRepository interface {
	CreateAdditionalRepair(ctx context.Context, additionalRepair *entities.AdditionalRepair) (*entities.AdditionalRepair, error)
	GetByID(ctx context.Context, id string) (*entities.AdditionalRepair, error)
	GetByServiceOrderID(ctx context.Context, serviceOrderID string) ([]entities.AdditionalRepair, error)
	UpdateAdditionalRepair(ctx context.Context, additionalRepair *entities.AdditionalRepair) (*entities.AdditionalRepair, error)
}
