package gateway

import (
	"context"
	"errors"

	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/request"
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/logs"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/interfaces"
)

type BillingServiceGateway struct {
	repo            interfaces.IBillingServiceRepository
	serviceRepo     interfaces.IServiceGateway
	partsSupplyRepo interfaces.IPartsSupplyGateway
}

func NewBillingServiceGateway(repo interfaces.IBillingServiceRepository) *BillingServiceGateway {
	return &BillingServiceGateway{
		repo: repo,
	}
}

func (g *BillingServiceGateway) CreateEstimate(ctx context.Context, serviceOrder *entities.ServiceOrder, additionalRepair *entities.AdditionalRepair) (*entities.Estimate,
	error) {
	logger := logs.Logger()
	if serviceOrder == nil && additionalRepair == nil {
		return nil, errors.New("no service order and additional repair provided")
	}

	request := &request.EstimateRequest{
		Services:      mapServicesDomainToRequest(serviceOrder.Services),
		PartsSupplies: mapPartsSupplyDomainToRequest(serviceOrder.PartsSupplies),
	}

	if serviceOrder != nil {
		request.ServiceOrderID = serviceOrder.ID
	}
	if additionalRepair != nil {
		request.AdditionalRepairID = additionalRepair.ID
	}

	response, err := g.repo.CreateEstimate(ctx, request)
	if err != nil {
		logger.Error().Err(err).Msg("error creating estimate")
		return nil, err
	}

	estimateResponse := &entities.Estimate{
		ID:                 response.ID,
		ServiceOrderID:     response.ServiceOrderID,
		AdditionalRepairID: response.AdditionalRepairID,
		Value:              response.Value,
		Status:             response.Status,
	}

	return estimateResponse, nil
}

func (g *BillingServiceGateway) ApproveEstimate(ctx context.Context, serviceOrder *entities.ServiceOrder, additionalRepair *entities.AdditionalRepair) (*entities.Estimate, error) {
	logger := logs.Logger()
	if serviceOrder == nil && additionalRepair == nil {
		return nil, errors.New("no service order and additional repair provided")
	}

	estimate := &request.EstimateRequest{
		ID: serviceOrder.Estimate.ID,
	}

	if serviceOrder != nil {
		estimate.ServiceOrderID = serviceOrder.ID
	}
	if additionalRepair != nil {
		estimate.AdditionalRepairID = additionalRepair.ID
	}

	response, err := g.repo.ApproveEstimate(ctx, estimate)
	if err != nil {
		logger.Error().Err(err).Msg("error approving estimate")
		return nil, err
	}

	if response == nil {
		logger.Error().Msg("error approving estimate: the value from estimate approval is nil")
		return nil, errors.New("error approving estimate: the value from estimate approval is nil")
	}

	return &entities.Estimate{
		ID:     response.ID,
		Value:  response.Value,
		Status: response.Status,
	}, nil
}

func (g *BillingServiceGateway) RejectEstimate(ctx context.Context, serviceOrder *entities.ServiceOrder, additionalRepair *entities.AdditionalRepair) (*entities.Estimate, error) {
	logger := logs.Logger()
	if serviceOrder == nil && additionalRepair == nil {
		return nil, errors.New("no service order and additional repair provided")
	}

	estimate := &request.EstimateRequest{
		ID: serviceOrder.Estimate.ID,
	}

	if serviceOrder != nil {
		estimate.ServiceOrderID = serviceOrder.ID
	}
	if additionalRepair != nil {
		estimate.AdditionalRepairID = additionalRepair.ID
	}

	response, err := g.repo.RejectEstimate(ctx, estimate)
	if err != nil {
		logger.Error().Err(err).Msg("error rejecting estimate")
		return nil, err
	}

	if response == nil {
		logger.Error().Msg("error rejecting estimate: the value from estimate approval is nil")
		return nil, errors.New("error rejecting estimate: the value from estimate approval is nil")
	}

	return &entities.Estimate{
		ID:     response.ID,
		Value:  response.Value,
		Status: response.Status,
	}, nil
}

func (g *BillingServiceGateway) CancelEstimate(ctx context.Context, serviceOrder *entities.ServiceOrder, additionalRepair *entities.AdditionalRepair) (*entities.Estimate, error) {
	logger := logs.Logger()
	if serviceOrder == nil && additionalRepair == nil {
		return nil, errors.New("no service order and additional repair provided")
	}

	estimate := &request.EstimateRequest{
		ID: serviceOrder.Estimate.ID,
	}

	if serviceOrder != nil {
		estimate.ServiceOrderID = serviceOrder.ID
	}

	if additionalRepair != nil {
		estimate.AdditionalRepairID = additionalRepair.ID
	}

	response, err := g.repo.CancelEstimate(ctx, estimate)
	if err != nil {
		logger.Error().Err(err).Msg("error cancelling estimate")
		return nil, err
	}

	if response == nil {
		logger.Error().Msg("error cancelling estimate: the value from estimate approval is nil")
		return nil, errors.New("error cancelling estimate: the value from estimate approval is nil")
	}

	return &entities.Estimate{
		ID:     response.ID,
		Value:  response.Value,
		Status: response.Status,
	}, nil
}

func (g *BillingServiceGateway) CreatePayment(ctx context.Context, estimateID string) (*entities.Payment, error) {
	logger := logs.Logger()
	if estimateID == "" {
		return nil, errors.New("no estimate ID provided")
	}

	payment, err := g.repo.CreatePayment(ctx, estimateID)
	if err != nil {
		logger.Error().Err(err).Msg("error creating payment")
		return nil, err
	}

	if payment == nil {
		logger.Error().Msg("error creating payment: the payment is nil")
		return nil, errors.New("error creating payment: the payment is nil")
	}

	return &entities.Payment{
		ID:          payment.ID,
		EstimateID:  payment.EstimateID,
		PaymentDate: payment.PaymentDate,
		Amount:      payment.Amount,
	}, nil
}

func (g *BillingServiceGateway) GetPaymentByEstimateID(ctx context.Context, estimateID string) (*entities.Payment, error) {
	logger := logs.Logger()
	if estimateID == "" {
		return nil, errors.New("no estimate ID provided")
	}

	payment, err := g.repo.GetPaymentByEstimateID(ctx, estimateID)
	if err != nil {
		logger.Error().Err(err).Msg("error getting payment by estimate ID")
		return nil, err
	}

	if payment == nil {
		logger.Error().Msg("error creating payment: the payment is nil")
		return nil, errors.New("error creating payment: the payment is nil")
	}

	return &entities.Payment{
		ID:          payment.ID,
		EstimateID:  payment.EstimateID,
		PaymentDate: payment.PaymentDate,
		Amount:      payment.Amount,
	}, nil
}

func mapServicesDomainToRequest(services []entities.Service) []request.ServiceRequest {
	servicesRequest := make([]request.ServiceRequest, 0, len(services))
	for _, s := range services {
		service := request.ServiceRequest{
			ID:          s.ID,
			Name:        s.Name,
			Description: s.Description,
			Price:       s.Price,
		}
		servicesRequest = append(servicesRequest, service)
	}
	return servicesRequest
}

func mapPartsSupplyDomainToRequest(partsSupplies []entities.PartsSupply) []request.PartsSupplyRequest {
	partsSuppliesRequest := make([]request.PartsSupplyRequest, 0, len(partsSupplies))
	for _, s := range partsSupplies {
		ps := request.PartsSupplyRequest{
			ID:       s.ID,
			Price:    s.Price,
			Quantity: s.Quantity,
		}
		partsSuppliesRequest = append(partsSuppliesRequest, ps)
	}
	return partsSuppliesRequest
}
