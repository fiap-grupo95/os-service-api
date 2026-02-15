package billing_service

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

type BillingServiceRepository struct {
	http *http.Client
}

func NewBillingServiceRepository() *BillingServiceRepository {
	return &BillingServiceRepository{
		http: http.DefaultClient,
	}
}

func (r *BillingServiceRepository) CreateEstimate(ctx context.Context, request *request.EstimateRequest) (*response.EstimateResponse, error) {
	logger := logs.Logger()
	if request == nil {
		return nil, errors.New("no service order provided")
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	body := io.NopCloser(bytes.NewReader(payload))
	resp, err := r.http.Post("", "application/json", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		logger.Error().Err(err).Msg("failed to create estimate")
		return nil, fmt.Errorf("failed to create estimate: status %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var response response.EstimateResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (r *BillingServiceRepository) RejectEstimate(ctx context.Context, request *request.EstimateRequest) (*response.EstimateResponse, error){
	return nil, nil
}

func (r *BillingServiceRepository) ApproveEstimate(ctx context.Context, request *request.EstimateRequest) (*response.EstimateResponse, error){
	return nil, nil
}

func (r *BillingServiceRepository) CancelEstimate(ctx context.Context, request *request.EstimateRequest) (*response.EstimateResponse, error){
	return nil, nil
}

func (r *BillingServiceRepository) CreatePayment(ctx context.Context, estimateID string) error {
	return nil
}

func (r *BillingServiceRepository) GetPaymentByEstimateID (ctx context.Context, estimateID string) (*response.PaymentResponse, error){
	return nil, nil
}