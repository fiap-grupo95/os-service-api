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

func (g *BillingServiceGateway) CreateEstimate(ctx context.Context, serviceOrder *entities.ServiceOrder) (*float64,
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

	return response.Value, nil
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
			ID:          s.ID,
			Name:        s.Name,
			Description: s.Description,
			Price:       s.Price,
		}
		if s.QuantityTotal > 0 {
			ps.Quantity = s.QuantityTotal
		} else {
			ps.Quantity = s.QuantityReserve
		}

		partsSuppliesRequest = append(partsSuppliesRequest, ps)
	}
	return partsSuppliesRequest
}
