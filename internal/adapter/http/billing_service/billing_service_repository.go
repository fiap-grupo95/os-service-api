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
	resp, err := r.http.Post(ESTIMATE_CREATE_ENDPOINT, "application/json", body)
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

func (r *BillingServiceRepository) RejectEstimate(ctx context.Context, request *request.EstimateRequest) (*response.EstimateResponse, error) {
	logger := logs.Logger()
	if request == nil || request.ID == "" {
		return nil, errors.New("no estimate ID provided")
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(ESTIMATE_REJECT_ENDPOINT, request.ID)
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
		logger.Error().Err(err).Msg("failed to reject estimate")
		return nil, fmt.Errorf("failed to reject estimate: status %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var response response.EstimateResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (r *BillingServiceRepository) ApproveEstimate(ctx context.Context, request *request.EstimateRequest) (*response.EstimateResponse, error) {
	logger := logs.Logger()
	if request == nil || request.ID == "" {
		return nil, errors.New("no estimate ID provided")
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(ESTIMATE_APPROVE_ENDPOINT, request.ID)
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
		logger.Error().Err(err).Msg("failed to approve estimate")
		return nil, fmt.Errorf("failed to approve estimate: status %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var response response.EstimateResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (r *BillingServiceRepository) CancelEstimate(ctx context.Context, request *request.EstimateRequest) (*response.EstimateResponse, error) {
	logger := logs.Logger()
	if request == nil || request.ID == "" {
		return nil, errors.New("no estimate ID provided")
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(ESTIMATE_CANCEL_ENDPOINT, request.ID)
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
		logger.Error().Err(err).Msg("failed to cancel estimate")
		return nil, fmt.Errorf("failed to cancel estimate: status %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var response response.EstimateResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (r *BillingServiceRepository) CreatePayment(ctx context.Context, estimateID string) (*response.PaymentResponse, error) {
	logger := logs.Logger()
	if estimateID == "" {
		return nil, errors.New("no estimate ID provided")
	}

	payload := map[string]string{"estimate_id": estimateID}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, PAYMENT_CREATE_ENDPOINT, bytes.NewReader(payloadBytes))
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
		logger.Error().Err(err).Msg("failed to create payment")
		return nil, fmt.Errorf("failed to create payment: status %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var response response.PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (r *BillingServiceRepository) GetPaymentByEstimateID(ctx context.Context, estimateID string) (*response.PaymentResponse, error) {
	logger := logs.Logger()
	if estimateID == "" {
		return nil, errors.New("no estimate ID provided")
	}

	url := fmt.Sprintf(PAYMENT_BY_ESTIMATE_ID_ENDPOINT, estimateID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		logger.Error().Err(err).Msg("failed to get payment by estimate ID")
		return nil, fmt.Errorf("failed to get payment by estimate ID: status %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var response response.PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}
