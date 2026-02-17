package gateway

import (
	"context"
	"errors"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/logs"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/interfaces"
	"github.com/newrelic/go-agent/v3/newrelic"
)

const (
	ErrInvalidStatus = "invalid status to create service order"
	InvalidID        = "invalid id"
)

type ServiceOrderGateway struct {
	repo            interfaces.IServiceOrderRepository
	vehicleRepo     interfaces.IVehicleGateway
	customerRepo    interfaces.ICustomerGateway
	partsSupplyRepo interfaces.IPartsSupplyGateway
	serviceRepo     interfaces.IServiceGateway
}

func NewServiceOrderGateway(repository interfaces.IServiceOrderRepository, vehicleRepo interfaces.IVehicleGateway, customerRepo interfaces.ICustomerGateway, partsSupplyRepo interfaces.IPartsSupplyGateway, serviceRepo interfaces.IServiceGateway) *ServiceOrderGateway {
	return &ServiceOrderGateway{
		repo:            repository,
		vehicleRepo:     vehicleRepo,
		customerRepo:    customerRepo,
		partsSupplyRepo: partsSupplyRepo,
		serviceRepo:     serviceRepo,
	}
}

func (s *ServiceOrderGateway) Create(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error) {
	logger := logs.Logger()
	if !serviceOrder.Status.IsRecebida() {
		return nil, errors.New(ErrInvalidStatus)
	}

	createdServiceOrder, err := s.repo.Create(ctx, serviceOrder)
	if err != nil {
		logger.Error().Msg(err.Error())
		return nil, err
	}
	return createdServiceOrder, nil
}

func (s *ServiceOrderGateway) GetByID(ctx context.Context, id string, isFullData bool) (*entities.ServiceOrder, error) {
	logger := logs.Logger()
	if txn := newrelic.FromContext(ctx); txn != nil {
		logger = logs.LoggerWithContext(ctx)
	}

	serviceOrderRecord, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Error().Msg(err.Error())
		return nil, err
	}

	if serviceOrderRecord == nil {
		return nil, nil
	}

	if !isFullData {
		return serviceOrderRecord, nil
	}

	serviceOrder := serviceOrderRecord

	vehicle, err := s.vehicleRepo.FindByID(serviceOrderRecord.VehicleID)
	if err != nil {
		logger.Error().Msg(err.Error())
		return nil, err
	}

	customer, err := s.customerRepo.GetByID(serviceOrderRecord.CustomerID)
	if err != nil {
		logger.Error().Msg(err.Error())
		return nil, err
	}

	partsSupplies, err := s.getPartsSupplies(ctx, serviceOrderRecord.PartsSupplies)
	if err != nil {
		return nil, err
	}

	services, err := s.getService(ctx, serviceOrderRecord.Services)
	if err != nil {
		return nil, err
	}

	if len(serviceOrderRecord.AdditionalRepairs) > 0 {
		for i, ar := range serviceOrderRecord.AdditionalRepairs {
			partsSupplies, err := s.getPartsSupplies(ctx, ar.PartsSupplies)
			if err != nil {
				return nil, err
			}
			services, err := s.getService(ctx, ar.Services)
			if err != nil {
				return nil, err
			}

			serviceOrder.AdditionalRepairs[i].PartsSupplies = partsSupplies
			serviceOrder.AdditionalRepairs[i].Services = services
		}
	}

	serviceOrder.Vehicle = vehicle
	serviceOrder.Customer = customer
	serviceOrder.PartsSupplies = partsSupplies
	serviceOrder.Services = services

	return serviceOrder, nil
}

func (s *ServiceOrderGateway) getPartsSupplies(ctx context.Context, partsSupplies []entities.PartsSupply) (partsSuppliesList []entities.PartsSupply, err error) {
	logger := logs.Logger()

	for _, ps := range partsSupplies {
		partsSupply, err := s.partsSupplyRepo.GetByID(ctx, ps.ID)
		if err != nil {
			logger.Error().Msg(err.Error())
			return nil, err
		}
		partsSuppliesList = append(partsSuppliesList, *partsSupply)
	}

	return partsSuppliesList, nil
}

func (s *ServiceOrderGateway) getService(ctx context.Context, services []entities.Service) (servicesList []entities.Service, err error) {
	logger := logs.Logger()

	for _, svc := range services {
		service, err := s.serviceRepo.GetByID(ctx, svc.ID)
		if err != nil {
			logger.Error().Msg(err.Error())
			return nil, err
		}
		servicesList = append(servicesList, *service)
	}
	return servicesList, nil
}

func (s *ServiceOrderGateway) Update(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error) {
	logger := logs.Logger()
	if serviceOrder == nil || serviceOrder.ID == "" {
		return nil, errors.New(InvalidID)
	}
	so, err := s.repo.Update(ctx, serviceOrder)
	if err != nil {
		logger.Error().Msg(err.Error())
		return nil, err
	}
	return so, nil
}

func (s *ServiceOrderGateway) List(ctx context.Context) ([]*entities.ServiceOrder, error) {
	logger := logs.Logger()
	var serviceOrders []*entities.ServiceOrder

	dtoList, err := s.repo.List(ctx)
	if err != nil {
		logger.Error().Msg(err.Error())
		return nil, err
	}

	for _, so := range dtoList {
		serviceOrders = append(serviceOrders, so)
	}

	return serviceOrders, nil
}
