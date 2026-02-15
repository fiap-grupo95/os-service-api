package interfaces

import (
	"context"

	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/request"
	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
)

type IExecutionGateway interface {
	CreateExecution(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.Execution, error)
	FinishExecution(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.Execution, error)
}

type IExecutionRepository interface {
	CreateExecution(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error)
	FinishExecution(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error)
}

