package operations

import (
	"context"
	"errors"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/logs"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/interfaces"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// ApproveEstimateStrategy - Estratégia para aprovar estimativa
type ApproveEstimateStrategy struct {
	serviceOrderRepo   interfaces.IServiceOrderGateway
	billingServiceRepo interfaces.IBillingServiceGateway
	partsSupplyRepo    interfaces.IPartsSupplyGateway
}

func (s *ApproveEstimateStrategy) GetTargetStatus() valueobject.ServiceOrderStatus {
	return valueobject.StatusAprovada
}

func (s *ApproveEstimateStrategy) Execute(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error) {
	logger := logs.LoggerWithContext(ctx)
	if txn := newrelic.FromContext(ctx); txn != nil {
		startSegment := txn.StartSegment("ApproveEstimateStrategy.Execute")
		defer startSegment.End()
	}

	// Write off parts supply - Dar baixa no estoque
	if err := s.partsSupplyRepo.WriteOff(ctx, serviceOrder.PartsSupplies); err != nil {
		logger.Error().Err(err).Any("parts_supply_id", serviceOrder.PartsSupplies).Msg("Error to write off parts supply")
		return nil, err
	}

	// Approve estimate - Aprovar orçamento
	estimate, err := s.billingServiceRepo.ApproveEstimate(ctx, serviceOrder)
	if err != nil {
		logger.Error().Err(err).Any("parts_supply_id", serviceOrder.PartsSupplies).Msg("Error approving estimate")
		return nil, err
	}
	if estimate == nil {
		logger.Error().Msg("Error approving estimate: the value from estimate approval is nil")
		return nil, errors.New("Error approving estimate: the value from estimate approval is nil")
	}
	serviceOrder.Estimate = estimate

	serviceOrder.Status = s.GetTargetStatus()
	return serviceOrder, nil
}
