package operations

import (
	"context"
	"errors"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/logs"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/interfaces"
)

// CancelEstimateStrategy - Estratégia para cancelar estimativa
type CancelEstimateStrategy struct{
	partsSupplyRepo interfaces.IPartsSupplyGateway
	billingServiceRepo interfaces.IBillingServiceGateway
}
 
func (s *CancelEstimateStrategy) GetTargetStatus() valueobject.ServiceOrderStatus {
    return valueobject.StatusEmDiagnostico
}
 
func (s *CancelEstimateStrategy) Execute(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error) {
    logger := logs.LoggerWithContext(ctx)
    
    if err := s.partsSupplyRepo.Release(ctx, serviceOrder.PartsSupplies); err != nil {
        logger.Error().Err(err).Any("parts_supply_id", serviceOrder.PartsSupplies).Msg("Error unreserving parts supply")
        return nil, err
    }

	// TODO: Implementar a lógica de cancelamento de orçamento
	estimate, err := s.billingServiceRepo.CancelEstimate(ctx, serviceOrder)
    if err != nil {
        logger.Error().Err(err).Any("parts_supply_id", serviceOrder.PartsSupplies).Msg("Error canceling estimate")
        return nil, err
    }
    if estimate == nil {
        logger.Error().Msg("Error canceling estimate: the value from estimate approval is nil")
        return nil, errors.New("error canceling estimate: the value from estimate approval is nil")
    }
    serviceOrder.Estimate = estimate    
	
    serviceOrder.Status = s.GetTargetStatus()
    return serviceOrder, nil
}