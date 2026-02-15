package execution_service

import (
	"context"

	"net/http"

	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/request"
	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/response"
)

type ExecutionServiceRepository struct {
	http *http.Client
}

func NewExecutionServiceRepository() *ExecutionServiceRepository {
	return &ExecutionServiceRepository{
		http: http.DefaultClient,
	}
}

func (r *ExecutionServiceRepository)CreateExecution(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error){
	return nil, nil
}

func (r *ExecutionServiceRepository)FinishExecution(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error){
	return nil, nil
}