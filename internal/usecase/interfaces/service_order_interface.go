package interfaces

import (
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	dto "github.com/fiap-grupo95/os-service-api/internal/infrastructure/database/model"
)

type IServiceOrderGateway interface {
	Create(serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error)
	GetByID(id uint) (*entities.ServiceOrder, error)
	Update(serviceOrder *entities.ServiceOrder) error
	List() ([]*entities.ServiceOrder, error)
	UpdateEstimate(id uint, estimate float64) error
	GetPartsSupplyServiceOrder(partsSupplyID uint, serviceOrderID uint) (*entities.ServiceOrderPartsSupply, error)
}

type IServiceOrderRepository interface {
	Create(serviceOrderDto *dto.ServiceOrderModel) (*dto.ServiceOrderModel, error)
	GetByID(id uint) (*dto.ServiceOrderModel, error)
	Update(serviceOrderDto *dto.ServiceOrderModel) error
	List() ([]*dto.ServiceOrderModel, error)
	UpdateEstimate(id uint, estimate float64) error
	GetPartsSupplyServiceOrder(partsSupplyID uint, serviceOrderID uint) (*dto.PartsSupplyServiceOrder, error)
}

