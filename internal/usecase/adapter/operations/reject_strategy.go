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

// RejectEstimateStrategy - Estratégia para rejeitar estimativa
type RejectEstimateStrategy struct{
	partsSupplyRepo interfaces.IPartsSupplyGateway
	billingServiceRepo interfaces.IBillingServiceGateway
}
 
func (s *RejectEstimateStrategy) GetTargetStatus() valueobject.ServiceOrderStatus {
    return valueobject.StatusRejeitada
}
 
func (s *RejectEstimateStrategy) Execute(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error) {
    logger := logs.LoggerWithContext(ctx)
	if txn := newrelic.FromContext(ctx); txn != nil {
		startSegment := txn.StartSegment("RejectEstimateStrategy.Execute")
		defer startSegment.End()
	}
    
	if err := s.partsSupplyRepo.Release(ctx, serviceOrder.PartsSupplies); err != nil {
		logger.Error().Err(err).Any("parts_supply_id", serviceOrder.PartsSupplies).Msg("Error unreserving parts supply")
		return nil, err
	}

	estimate, err := s.billingServiceRepo.RejectEstimate(ctx, serviceOrder)
	if err != nil {
		logger.Error().Err(err).Any("parts_supply_id", serviceOrder.PartsSupplies).Msg("Error rejecting estimate")
		return nil, err
	}
	if estimate == nil {
		logger.Error().Msg("Error rejecting estimate: the value from estimate approval is nil")
		return nil, errors.New("error approving estimate: the value from estimate approval is nil")
	}
	serviceOrder.Estimate = estimate

    
    serviceOrder.Status = s.GetTargetStatus()
    return serviceOrder, nil
}