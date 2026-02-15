package execution_service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/request"
	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/logs"
)

type ExecutionServiceRepository struct {
	http *http.Client
}

func NewExecutionServiceRepository() *ExecutionServiceRepository {
	return &ExecutionServiceRepository{
		http: http.DefaultClient,
	}
}

func (r *ExecutionServiceRepository) CreateExecution(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
	logger := logs.Logger()
	if execution == nil {
		return nil, errors.New("no execution request provided")
	}

	payload, err := json.Marshal(execution)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, EXECUTION_SERVICE_ENDPOINT, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		logger.Error().Err(err).Msg("failed to create execution")
		return nil, fmt.Errorf("failed to create execution: status %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var response response.ExecutionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (r *ExecutionServiceRepository) FinishExecution(ctx context.Context, execution *request.ExecutionRequest) (*response.ExecutionResponse, error) {
	logger := logs.Logger()
	if execution == nil {
		return nil, errors.New("no execution request provided")
	}
	if execution.ID == "" {
		return nil, errors.New("no service order ID provided")
	}

	payload, err := json.Marshal(execution)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(EXECUTION_SERVICE_FINISH_ENDPOINT, execution.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		logger.Error().Err(err).Msg("failed to finish execution")
		return nil, fmt.Errorf("failed to finish execution: status %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var response response.ExecutionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}
