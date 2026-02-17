package gateway

import (
	"context"
	"errors"

	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/request"
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/logs"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/interfaces"
)

type ExecutionServiceGateway struct {
	executionRepository interfaces.IExecutionRepository
}

func NewExecutionServiceGateway(executionRepository interfaces.IExecutionRepository) *ExecutionServiceGateway {
	return &ExecutionServiceGateway{
		executionRepository: executionRepository,
	}
}

func (g *ExecutionServiceGateway) CreateExecution(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.Execution, error) {
	logger := logs.Logger()
	if serviceOrder == nil {
		return nil, errors.New("no service order provided")
	}

	execution := &request.ExecutionRequest{
		ServiceOrderID: serviceOrder.ID,
	}

	ExecutionResponse, err := g.executionRepository.CreateExecution(ctx, execution)
	if err != nil {
		logger.Error().Err(err).Msg("error creating execution")
		return nil, err
	}

	return &entities.Execution{
		ID:             ExecutionResponse.ID,
		ServiceOrderID: ExecutionResponse.ServiceOrderID,
		Status:         ExecutionResponse.Status,
		StartedAt:      ExecutionResponse.StartedAt,
		FinishedAt:     ExecutionResponse.FinishedAt,
	}, nil
}

func (g *ExecutionServiceGateway) FinishExecution(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.Execution, error) {
	logger := logs.Logger()
	if serviceOrder == nil {
		return nil, errors.New("no service order provided")
	}

	execution := &request.ExecutionRequest{
		ID:             serviceOrder.Execution.ID,
		ServiceOrderID: serviceOrder.ID,
	}

	ExecutionResponse, err := g.executionRepository.FinishExecution(ctx, execution)
	if err != nil {
		logger.Error().Err(err).Msg("error finishing execution")
		return nil, err
	}

	return &entities.Execution{
		ID:             ExecutionResponse.ID,
		ServiceOrderID: ExecutionResponse.ServiceOrderID,
		Status:         ExecutionResponse.Status,
		StartedAt:      ExecutionResponse.StartedAt,
		FinishedAt:     ExecutionResponse.FinishedAt,
	}, nil
}
