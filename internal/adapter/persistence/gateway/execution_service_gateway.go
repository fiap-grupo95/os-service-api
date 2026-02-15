package gateway

import (
	"context"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/interfaces"
)

type ExecutionServiceGateway struct {
	executionServiceRepository interfaces.IExecutionRepository
}

func NewExecutionServiceGateway(executionServiceRepository interfaces.IExecutionRepository) *ExecutionServiceGateway {
	return &ExecutionServiceGateway{
		executionServiceRepository: executionServiceRepository,
	}
}

func (g *ExecutionServiceGateway)CreateExecution(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.Execution, error){
	return nil, nil
}

func (g *ExecutionServiceGateway)FinishExecution(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.Execution, error){
	return nil, nil
}