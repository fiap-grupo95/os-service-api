package gateway

import (
	"errors"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
	dto "github.com/fiap-grupo95/os-service-api/internal/infrastructure/database/model"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/logs"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/interfaces"
)

const (
	ErrInvalidStatus = "invalid status to create service order"
)

type ServiceOrderGateway struct {
	repo interfaces.IServiceOrderRepository
}

func NewServiceOrderGateway(repository interfaces.IServiceOrderRepository) *ServiceOrderGateway {
	return &ServiceOrderGateway{
		repo: repository,
	}
}

func (s *ServiceOrderGateway) Create(serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error){
	logger := logs.Logger()	
	if serviceOrder.ServiceOrderStatus != valueobject.StatusRecebida {
		return nil, errors.New(ErrInvalidStatus)
	}

	serviceOrderDto := &dto.ServiceOrderModel{
		ID:                   serviceOrder.ID,
		CustomerID:           serviceOrder.CustomerID,
		VehicleID:            serviceOrder.VehicleID,
		Estimate:             serviceOrder.Estimate,
		StartedExecutionDate: serviceOrder.StartedExecutionDate,
		FinalExecutionDate:   serviceOrder.FinalExecutionDate,
		CreatedAt:            serviceOrder.CreatedAt,
		UpdatedAt:            serviceOrder.UpdatedAt,
	}

	createdServiceOrder, err := s.repo.Create(serviceOrderDto)
	if err != nil {
		logger.Error().Msg(err.Error())
		return nil, err
	}
	return createdServiceOrder.ToDomain(), nil
}

func (s *ServiceOrderGateway) GetByID(id uint) (*entities.ServiceOrder, error){
	// TODO: Implement me
	return nil, nil
}

func (s *ServiceOrderGateway) Update(serviceOrder *entities.ServiceOrder) error{
	// TODO: Implement me
	return nil
}

func (s *ServiceOrderGateway) List() ([]*entities.ServiceOrder, error){
	// TODO: Implement me
	return nil, nil
}

func (s *ServiceOrderGateway) UpdateEstimate(id uint, estimate float64) error{
	// TODO: Implement me
	return nil
}

func (s *ServiceOrderGateway) GetPartsSupplyServiceOrder(partsSupplyID uint, serviceOrderID uint) (*entities.ServiceOrderPartsSupply, error){
	// TODO: Implement me
	return nil, nil
}