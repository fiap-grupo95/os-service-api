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

func (g *BillingServiceGateway) CreateEstimate(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.Estimate,
	error) {
	logger := logs.Logger()
	if serviceOrder == nil {
		return nil, errors.New("no service order provided")
	}

	request := &request.EstimateRequest{
		ServiceOrderID: serviceOrder.ID,
		Services:       mapServicesDomainToRequest(serviceOrder.Services),
		PartsSupplies:  mapPartsSupplyDomainToRequest(serviceOrder.PartsSupplies),
	}

	response, err := g.repo.CreateEstimate(ctx, request)
	if err != nil {
		logger.Error().Err(err).Msg("error creating estimate")
	}

	return &entities.Estimate{
		ID:     response.ID,
		Value:  response.Value,
		Status: response.Status,
	}, nil
}

func (g *BillingServiceGateway) ApproveEstimate(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.Estimate, error) {
	logger := logs.Logger()
	if serviceOrder == nil {
		return nil, errors.New("no service order provided")
	}

	estimate := &request.EstimateRequest{
		ID:             serviceOrder.Estimate.ID,
		ServiceOrderID: serviceOrder.ID,
		Services:       mapServicesDomainToRequest(serviceOrder.Services),
		PartsSupplies:  mapPartsSupplyDomainToRequest(serviceOrder.PartsSupplies),
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

func (g *BillingServiceGateway) RejectEstimate(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.Estimate, error) {
	logger := logs.Logger()
	if serviceOrder == nil {
		return nil, errors.New("no service order provided")
	}

	estimate := &request.EstimateRequest{
		ID:             serviceOrder.Estimate.ID,
		ServiceOrderID: serviceOrder.ID,
		Services:       mapServicesDomainToRequest(serviceOrder.Services),
		PartsSupplies:  mapPartsSupplyDomainToRequest(serviceOrder.PartsSupplies),
	}

	response, err := g.repo.RejectEstimate(ctx, estimate)
	if err != nil {
		logger.Error().Err(err).Msg("error rejecting estimate")
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

func (g *BillingServiceGateway) CancelEstimate(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.Estimate, error) {
	logger := logs.Logger()
	if serviceOrder == nil {
		return nil, errors.New("no service order provided")
	}

	estimate := &request.EstimateRequest{
		ID:             serviceOrder.Estimate.ID,
		ServiceOrderID: serviceOrder.ID,
		Services:       mapServicesDomainToRequest(serviceOrder.Services),
		PartsSupplies:  mapPartsSupplyDomainToRequest(serviceOrder.PartsSupplies),
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

func (g *BillingServiceGateway) getServicesByIDs(ctx context.Context, services []entities.Service) ([]entities.Service, error) {
	logger := logs.Logger()
	if len(services) == 0 {
		return nil, errors.New("no services provided")
	}

	var serviceDb []entities.Service
	for _, s := range services {
		item, err := g.serviceRepo.GetByID(ctx, s.ID)
		if err != nil {
			logger.Error().Err(err).Any("service_id", s.ID).Msg("error getting service by ID")
			return nil, err
		}
		serviceDb = append(serviceDb, *item)
	}

	// Assuming we have a service repository to get the services by IDs
	return serviceDb, nil
}

func (g *BillingServiceGateway) getPartsSupplyByIDs(ctx context.Context, partsSupplies []entities.PartsSupply) ([]entities.PartsSupply, error) {
	logger := logs.Logger()
	if len(partsSupplies) == 0 {
		return nil, errors.New("no services provided")
	}

	var psDb []entities.PartsSupply
	for _, ps := range partsSupplies {
		item, err := g.partsSupplyRepo.GetByID(ctx, ps.ID)
		if err != nil {
			logger.Error().Err(err).Any("parts_supply_id", ps.ID).Msg("error getting Parts Supplies by ID")
			return nil, err
		}
		psDb = append(psDb, *item)
	}

	return psDb, nil
}

func mapServicesDomainToRequest(services []entities.Service) []request.ServiceRequest {
	servicesRequest := make([]request.ServiceRequest, len(services))
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
	partsSuppliesRequest := make([]request.PartsSupplyRequest, len(partsSupplies))
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
