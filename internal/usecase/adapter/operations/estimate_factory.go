package operations

import (
	"context"
	"errors"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
)

const(
	ESTIMATE_APPROVE = "estimate_approve"
	ESTIMATE_REJECT  = "estimate_reject"
	ESTIMATE_CANCEL  = "estimate_cancel"
	
	ErrInvalidEstimateOperation = "invalid estimate operation"
)

// EstimateOperationStrategy define a interface para operações de estimativa
type EstimateOperationStrategy interface {
    Execute(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error)
    GetTargetStatus() valueobject.ServiceOrderStatus
}

// EstimateStrategyFactory cria a estratégia apropriada baseada na operação
type EstimateStrategyFactory struct {
    strategies map[string]EstimateOperationStrategy
}

func NewEstimateStrategyFactory() *EstimateStrategyFactory {
    return &EstimateStrategyFactory{
        strategies: map[string]EstimateOperationStrategy{
            ESTIMATE_APPROVE: &ApproveEstimateStrategy{},
            ESTIMATE_REJECT:  &RejectEstimateStrategy{},
            ESTIMATE_CANCEL:  &CancelEstimateStrategy{},
        },
    }
}

func (f *EstimateStrategyFactory) GetStrategy(operation string) (EstimateOperationStrategy, error) {
    strategy, exists := f.strategies[operation]
    if !exists {
        return nil, errors.New(ErrInvalidEstimateOperation)
    }
    return strategy, nil
}