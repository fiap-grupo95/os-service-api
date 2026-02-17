package gateway

import (
	"context"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/interfaces"
)

type AdditionalRepairGateway struct {
	repo interfaces.IAdditionalRepairRepository
}

func NewAdditionalRepairGateway(repo interfaces.IAdditionalRepairRepository) *AdditionalRepairGateway {
	return &AdditionalRepairGateway{
		repo: repo,
	}
}

func (a *AdditionalRepairGateway) CreateAdditionalRepair(ctx context.Context, additionalRepair *entities.AdditionalRepair) (*entities.AdditionalRepair, error) {
	return a.repo.CreateAdditionalRepair(ctx, additionalRepair)
}

func (a *AdditionalRepairGateway) GetByID(ctx context.Context, id string) (*entities.AdditionalRepair, error) {
	return a.repo.GetByID(ctx, id)
}

func (a *AdditionalRepairGateway) GetByServiceOrderID(ctx context.Context, serviceOrderID string) ([]entities.AdditionalRepair, error) {
	return a.repo.GetByServiceOrderID(ctx, serviceOrderID)
}

func (a *AdditionalRepairGateway) UpdateAdditionalRepair(ctx context.Context, additionalRepair *entities.AdditionalRepair) (*entities.AdditionalRepair, error) {
	return a.repo.UpdateAdditionalRepair(ctx, additionalRepair)
}
