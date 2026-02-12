package repository

import (
	"context"
	"strings"

	dto "github.com/fiap-grupo95/os-service-api/internal/infrastructure/database/model"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/interfaces"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// ServiceOrderRepository implements IServiceOrderRepository interface
type ServiceOrderRepository struct {
	db *gorm.DB
}

var _ interfaces.IServiceOrderRepository = (*ServiceOrderRepository)(nil)

func NewServiceOrderRepository(db *gorm.DB) *ServiceOrderRepository {
	return &ServiceOrderRepository{db: db}
}

func (r *ServiceOrderRepository) Create(ctx context.Context, serviceOrderDto *dto.ServiceOrderModel) (*dto.ServiceOrderModel, error) {
	if serviceOrderDto == nil {
		return nil, gorm.ErrInvalidData
	}

	dtoStatus, err := r.getStatus(ctx, serviceOrderDto.ServiceOrderStatus)
	if err != nil {
		return nil, gorm.ErrInvalidData
	}

	serviceOrderDto.OSStatusID = dtoStatus.ID
	serviceOrderDto.ServiceOrderStatus = *dtoStatus

	// Begin transaction
	tx := r.db.Begin()

	if err := tx.Create(&serviceOrderDto).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return serviceOrderDto, nil
}

func (r *ServiceOrderRepository) GetByID(ctx context.Context, id uint) (*dto.ServiceOrderModel, error) {
	var serviceOrder dto.ServiceOrderModel
	// TODO - Avaliar o que posso tirar do Preload e deixar para serem carregados apenas quando necessÃ¡rio
	err := r.db.Preload("ServiceOrderStatus").
		Preload("AdditionalRepairs").
		Preload("AdditionalRepairs.ARStatus").
		Preload("AdditionalRepairs.Services").
		Preload("AdditionalRepairs.PartsSupplies").
		First(&serviceOrder, id).Error
	if err != nil {
		log.Error().Msgf("Error finding service order with id %d: %v", id, err)
		if strings.EqualFold(err.Error(), gorm.ErrRecordNotFound.Error()) {
			return nil, nil
		}
		return nil, err
	}
	return &serviceOrder, nil
}

func (r *ServiceOrderRepository) UpdateEstimate(ctx context.Context, id uint, estimate float64) error {
	var dtoDB dto.ServiceOrderModel
	if err := r.db.First(&dtoDB, id).Error; err != nil {
		return err
	}
	newEstimate := dtoDB.Estimate + estimate

	return r.db.Model(&dto.ServiceOrderModel{}).
		Where("id = ?", id).
		Update("estimate", newEstimate).Error
}

func (r *ServiceOrderRepository) Update(ctx context.Context, serviceOrder *dto.ServiceOrderModel) error {
	if serviceOrder == nil {
		return gorm.ErrInvalidData
	}

	dtoStatus, err := r.getStatus(ctx, serviceOrder.ServiceOrderStatus)
	if err != nil {
		return gorm.ErrInvalidData
	}

	tx := r.db.Begin()

	serviceOrderDto := dto.ServiceOrderModel{
		ID:                       serviceOrder.ID,
		CustomerID:               serviceOrder.CustomerID,
		VehicleID:                serviceOrder.VehicleID,
		OSStatusID:               dtoStatus.ID,
		Estimate:                 serviceOrder.Estimate,
		StartedExecutionDate:     serviceOrder.StartedExecutionDate,
		FinalExecutionDate:       serviceOrder.FinalExecutionDate,
		ExecutionDurationInHours: serviceOrder.ExecutionDurationInHours,
	}

	if err := tx.Model(&dto.ServiceOrderModel{}).Where("id = ?", serviceOrder.ID).Updates(&serviceOrderDto).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update PartsSupplies relationships
	if serviceOrder.PartsSupplies != nil {
		if err := tx.Where("service_order_id = ?", serviceOrder.ID).Delete(&dto.PartsSupplyServiceOrder{}).Error; err != nil {
			tx.Rollback()
			return err
		}

		for _, partsSupply := range serviceOrder.PartsSupplies {
			relation := dto.PartsSupplyServiceOrder{
				PartsSupplyID:  partsSupply.ID,
				ServiceOrderID: serviceOrder.ID,
				Quantity:       partsSupply.QuantityReserve,
			}
			if err := tx.Create(&relation).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Update Services relationships
	if serviceOrder.Services != nil {
		if err := tx.Where("service_order_id = ?", serviceOrder.ID).Delete(&dto.ServiceServiceOrder{}).Error; err != nil {
			tx.Rollback()
			return err
		}

		for _, service := range serviceOrder.Services {
			relation := dto.ServiceServiceOrder{
				ServiceID:      service.ID,
				ServiceOrderID: serviceOrder.ID,
			}
			if err := tx.Create(&relation).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

func (r *ServiceOrderRepository) List(ctx context.Context) ([]*dto.ServiceOrderModel, error) {
	var serviceOrders []dto.ServiceOrderModel
	// TODO - Avaliar o que posso tirar do Preload e deixar para serem carregados apenas quando necessÃ¡rio
	err := r.db.
		Preload("ServiceOrderStatus").
		Preload("AdditionalRepairs").
		Preload("AdditionalRepairs.ARStatus").
		Preload("AdditionalRepairs.Services").
		Preload("AdditionalRepairs.PartsSupplies").
		Find(&serviceOrders).Error
	if err != nil {
		return nil, err
	}

	result := make([]*dto.ServiceOrderModel, 0, len(serviceOrders))
	for _, so := range serviceOrders {
		result = append(result, &so)
	}
	return result, nil
}

func (r *ServiceOrderRepository) getStatus(ctx context.Context, status dto.ServiceOrderStatus) (*dto.ServiceOrderStatus, error) {
	var serviceOrderStatuses dto.ServiceOrderStatus
	err := r.db.Where("description = ?", status.Description).First(&serviceOrderStatuses).Error
	if err != nil {
		return nil, err
	}
	return &serviceOrderStatuses, nil
}

func (r *ServiceOrderRepository) GetPartsSupplyServiceOrder(ctx context.Context, partsSupplyID uint, serviceOrderID uint) (*dto.PartsSupplyServiceOrder, error) {
	var relation dto.PartsSupplyServiceOrder
	err := r.db.Where("parts_supply_id = ? AND service_order_id = ?", partsSupplyID, serviceOrderID).First(&relation).Error
	if err != nil {
		return nil, err
	}
	return &relation, nil
}
