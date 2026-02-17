package gateway

import (
	"context"
	"errors"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	dto "github.com/fiap-grupo95/os-service-api/internal/infrastructure/database/model"
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

	serviceOrderDto := &dto.ServiceOrderModel{
		CustomerID: serviceOrder.CustomerID,
		VehicleID:  serviceOrder.VehicleID,
		ServiceOrderStatus: dto.ServiceOrderStatus{
			Description: serviceOrder.Status.String(),
		},
	}

	createdServiceOrder, err := s.repo.Create(ctx, serviceOrderDto)
	if err != nil {
		logger.Error().Msg(err.Error())
		return nil, err
	}
	return createdServiceOrder.ToDomain(), nil
}

func (s *ServiceOrderGateway) GetByID(ctx context.Context, id string, isFullData bool) (*entities.ServiceOrder, error) {
	logger := logs.Logger()
	if txn := newrelic.FromContext(ctx); txn != nil {
		logger = logs.LoggerWithContext(ctx)
	}

	serviceOrderModel, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Error().Msg(err.Error())
		return nil, err
	}

	if serviceOrderModel == nil {
		return nil, nil
	}

	if !isFullData {
		return serviceOrderModel.ToDomain(), nil
	}

	serviceOrder := serviceOrderModel.ToDomain()

	vehicle, err := s.vehicleRepo.FindByID(serviceOrderModel.VehicleID)
	if err != nil {
		logger.Error().Msg(err.Error())
		return nil, err
	}

	customer, err := s.customerRepo.GetByID(serviceOrderModel.CustomerID)
	if err != nil {
		logger.Error().Msg(err.Error())
		return nil, err
	}

	partsSupplies, err := s.getPartsSupplies(ctx, serviceOrderModel.PartsSupplies)
	if err != nil {
		return nil, err
	}

	services, err := s.getService(ctx, serviceOrderModel.Services)
	if err != nil {
		return nil, err
	}

	if len(serviceOrderModel.AdditionalRepairs) > 0 {
		for i, ar := range serviceOrderModel.AdditionalRepairs {
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

func (s *ServiceOrderGateway) getPartsSupplies(ctx context.Context, partsSupplies []dto.PartsSupplyModel) (partsSuppliesList []entities.PartsSupply, err error) {
	logger := logs.Logger()

	// TODO: Revisar modelo de dados de PartsSupplies do Model para armazenas apenas IDs
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

func (s *ServiceOrderGateway) getService(ctx context.Context, services []dto.ServiceModel) (servicesList []entities.Service, err error) {
	logger := logs.Logger()

	// TODO: Revisar modelo de dados de Services do Model para armazenas apenas IDs
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

func (s *ServiceOrderGateway) Update(ctx context.Context, serviceOrder *entities.ServiceOrder) error {
	// TODO: Implement me
	return nil
}

func (s *ServiceOrderGateway) List(ctx context.Context) ([]*entities.ServiceOrder, error) {
	logger := logs.Logger()
	var serviceOrders []*entities.ServiceOrder

	dtoList, err := s.repo.List(ctx)
	if err != nil {
		logger.Error().Msg(err.Error())
		return nil, err
	}

	for _, s := range dtoList {
		serviceOrders = append(serviceOrders, s.ToDomain())
	}

	return serviceOrders, nil
}

func (s *ServiceOrderGateway) UpdateEstimate(ctx context.Context, id string, estimate float64) error {
	// TODO: Implement me
	return nil
}

func (s *ServiceOrderGateway) GetPartsSupplyServiceOrder(ctx context.Context, partsSupplyID string, serviceOrderID string) (*entities.ServiceOrderPartsSupply, error) {
	// TODO: Implement me
	return nil, nil
}
