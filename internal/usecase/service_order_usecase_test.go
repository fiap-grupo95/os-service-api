package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
	dto "github.com/fiap-grupo95/os-service-api/internal/infrastructure/database/model"
	mocks "github.com/fiap-grupo95/os-service-api/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func toServiceOrderEntity(model *dto.ServiceOrderModel) *entities.ServiceOrder {
	if model == nil {
		return nil
	}
	return model.ToDomain()
}

const (
	StatusRecebida            = string(valueobject.StatusRecebida)
	StatusEmDiagnostico       = string(valueobject.StatusEmDiagnostico)
	StatusAguardandoAprovacao = string(valueobject.StatusAguardandoAprovacao)
	StatusAprovada            = string(valueobject.StatusAprovada)
	StatusRejeitada           = string(valueobject.StatusRejeitada)
	StatusEmExecucao          = string(valueobject.StatusEmExecucao)
	StatusFinalizada          = string(valueobject.StatusFinalizada)
	StatusEntregue            = string(valueobject.StatusEntregue)
	StatusCancelada           = string(valueobject.StatusCancelada)
)

func UintPointer(value uint) *uint {
	return &value
}

func TestCreateServiceOrder(t *testing.T) {
	t.Run("Success - Valid service order creation", func(t *testing.T) {
		vehicleRepo := new(mocks.MockVehicleGateway)
		customerRepo := new(mocks.MockCustomerGateway)
		serviceOrderRepo := new(mocks.MockServiceOrderGateway)
		serviceRepo := new(mocks.MockServiceGateway)
		partsSupplyRepo := new(mocks.MockPartsSupplyGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo)

		serviceOrder := entities.ServiceOrder{
			CustomerID:         1,
			VehicleID:          1,
			ServiceOrderStatus: valueobject.StatusRecebida,
		}

		vehicleRepo.On("FindByID", uint(1)).Return(&entities.Vehicle{ID: 1}, nil)
		customerRepo.On("GetByID", uint(1)).Return(&entities.Customer{ID: 1}, nil)
		serviceOrderRepo.On("Create", mock.AnythingOfType("*entities.ServiceOrder")).Return(&entities.ServiceOrder{
			ID:         1,
			CustomerID: 1,
			VehicleID:  1,
		}, nil)

		r, err := useCase.CreateServiceOrder(context.Background(), serviceOrder)
		assert.NoError(t, err)
		assert.NotNil(t, r)
	})

	t.Run("Error - Internal DB Error", func(t *testing.T) {
		vehicleRepo := new(mocks.MockVehicleGateway)
		customerRepo := new(mocks.MockCustomerGateway)
		serviceOrderRepo := new(mocks.MockServiceOrderGateway)
		serviceRepo := new(mocks.MockServiceGateway)
		partsSupplyRepo := new(mocks.MockPartsSupplyGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo)

		serviceOrder := entities.ServiceOrder{
			CustomerID: 1,
			VehicleID:  1,
		}

		vehicleRepo.On("FindByID", uint(1)).Return(&entities.Vehicle{ID: 1}, nil)
		customerRepo.On("GetByID", uint(1)).Return(&entities.Customer{ID: 1}, nil)
		serviceOrderRepo.On("Create", mock.AnythingOfType("*entities.ServiceOrder")).Return(nil, errors.New("internal db error"))

		r, err := useCase.CreateServiceOrder(context.Background(), serviceOrder)
		assert.Error(t, err)
		assert.Equal(t, "internal db error", err.Error())
		assert.Nil(t, r)
	})

	t.Run("Error - Vehicle not found", func(t *testing.T) {
		vehicleRepo := new(mocks.MockVehicleGateway)
		customerRepo := new(mocks.MockCustomerGateway)
		serviceOrderRepo := new(mocks.MockServiceOrderGateway)
		serviceRepo := new(mocks.MockServiceGateway)
		partsSupplyRepo := new(mocks.MockPartsSupplyGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo)

		serviceOrder := entities.ServiceOrder{
			CustomerID: 1,
			VehicleID:  1,
		}

		vehicleRepo.On("FindByID", uint(1)).Return(nil, nil)

		r, err := useCase.CreateServiceOrder(context.Background(), serviceOrder)
		assert.Error(t, err)
		assert.Equal(t, ErrVehicleNotFound, err)
		assert.Nil(t, r)
	})

	t.Run("Error - Vehicle Has invalid ID", func(t *testing.T) {
		vehicleRepo := new(mocks.MockVehicleGateway)
		customerRepo := new(mocks.MockCustomerGateway)
		serviceOrderRepo := new(mocks.MockServiceOrderGateway)
		serviceRepo := new(mocks.MockServiceGateway)
		partsSupplyRepo := new(mocks.MockPartsSupplyGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo)

		serviceOrder := entities.ServiceOrder{
			CustomerID: 1,
			VehicleID:  0,
		}

		vehicleRepo.On("FindByID", uint(0)).Return(nil, errors.New("invalid vehicle ID"))

		r, err := useCase.CreateServiceOrder(context.Background(), serviceOrder)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidID, err)
		assert.Nil(t, r)
	})

	t.Run("Error - Customer not found", func(t *testing.T) {
		vehicleRepo := new(mocks.MockVehicleGateway)
		customerRepo := new(mocks.MockCustomerGateway)
		serviceOrderRepo := new(mocks.MockServiceOrderGateway)
		serviceRepo := new(mocks.MockServiceGateway)
		partsSupplyRepo := new(mocks.MockPartsSupplyGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo)

		serviceOrder := entities.ServiceOrder{
			CustomerID: 1,
			VehicleID:  1,
		}

		vehicleRepo.On("FindByID", uint(1)).Return(&entities.Vehicle{ID: 1}, nil)
		customerRepo.On("GetByID", uint(1)).Return(nil, nil)

		r, err := useCase.CreateServiceOrder(context.Background(), serviceOrder)
		assert.Error(t, err)
		assert.Equal(t, ErrCustomerNotFound, err)
		assert.Nil(t, r)
	})

	t.Run("Error - Customer Has invalid ID", func(t *testing.T) {
		vehicleRepo := new(mocks.MockVehicleGateway)
		customerRepo := new(mocks.MockCustomerGateway)
		serviceOrderRepo := new(mocks.MockServiceOrderGateway)
		serviceRepo := new(mocks.MockServiceGateway)
		partsSupplyRepo := new(mocks.MockPartsSupplyGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo)

		serviceOrder := entities.ServiceOrder{
			CustomerID: 0,
			VehicleID:  1,
		}

		vehicleRepo.On("FindByID", uint(1)).Return(&entities.Vehicle{ID: 1}, nil)
		customerRepo.On("GetByID", uint(0)).Return(nil, errors.New("invalid customer ID"))

		r, err := useCase.CreateServiceOrder(context.Background(), serviceOrder)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidID, err)
		assert.Nil(t, r)
	})
}

func TestUpdateServiceOrder(t *testing.T) {
	t.Run("Success - Update to EmDiagnostico", func(t *testing.T) {
		vehicleRepo := new(mocks.MockVehicleGateway)
		customerRepo := new(mocks.MockCustomerGateway)
		serviceOrderRepo := new(mocks.MockServiceOrderGateway)
		serviceRepo := new(mocks.MockServiceGateway)
		partsSupplyRepo := new(mocks.MockPartsSupplyGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo)

		serviceOrder := entities.ServiceOrder{
			ID:                 1,
			ServiceOrderStatus: valueobject.StatusEmDiagnostico,
			Services: []entities.Service{
				{ID: 1},
			},
			PartsSupplies: []entities.PartsSupply{
				{
					ID:              1,
					QuantityReserve: 2,
				},
			},
		}
		flow := DIAGNOSIS

		serviceOrderRepo.On("GetByID", uint(1), false).Return(&dto.ServiceOrderModel{
			ID: 1,
			ServiceOrderStatus: dto.ServiceOrderStatus{
				ID:          1,
				Description: string(valueobject.StatusRecebida),
			},
		}, nil)
		serviceRepo.On("GetByID", mock.Anything, uint(1)).Return(&entities.Service{ID: 1}, nil)
		partsSupplyRepo.On("GetByID", mock.Anything, uint(1)).Return(&entities.PartsSupply{
			ID:              1,
			QuantityTotal:   10,
			QuantityReserve: 2,
		}, nil)
		partsSupplyRepo.On("Update", mock.Anything, mock.AnythingOfType("*entities.PartsSupply")).Return(nil)
		serviceOrderRepo.On("Update", mock.AnythingOfType("*entities.ServiceOrder")).Return(nil)

		r, err := useCase.UpdateServiceOrder(context.Background(), serviceOrder, flow)
		assert.NoError(t, err)
		assert.NotNil(t, r)
	})

	t.Run("Error - Invalid Status Transition", func(t *testing.T) {
		vehicleRepo := new(mocks.MockVehicleGateway)
		customerRepo := new(mocks.MockCustomerGateway)
		serviceOrderRepo := new(mocks.MockServiceOrderGateway)
		serviceRepo := new(mocks.MockServiceGateway)
		partsSupplyRepo := new(mocks.MockPartsSupplyGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo)

		serviceOrder := entities.ServiceOrder{
			ID:                 1,
			ServiceOrderStatus: valueobject.StatusEntregue,
		}
		flow := DIAGNOSIS

		serviceOrderRepo.On("GetByID", uint(1), false).Return(&dto.ServiceOrderModel{
			ID: 1,
			ServiceOrderStatus: dto.ServiceOrderStatus{
				ID:          1,
				Description: string(valueobject.StatusRecebida),
			},
		}, nil)

		r, err := useCase.UpdateServiceOrder(context.Background(), serviceOrder, flow)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidTransitionStatusToDiagnosis, err)
		assert.Nil(t, r)
	})
}

func TestValidateEstimate(t *testing.T) {
	t.Run("Should approve estimate and release parts supply", func(t *testing.T) {
		partsSupplyRepo := new(mocks.MockPartsSupplyGateway)
		serviceOrderRepo := new(mocks.MockServiceOrderGateway)

		request := &entities.ServiceOrder{
			ID:                 1,
			ServiceOrderStatus: valueobject.StatusAprovada,
		}
		serviceOrderDTO := &dto.ServiceOrderModel{
			ID: 1,
			ServiceOrderStatus: dto.ServiceOrderStatus{
				Description: StatusAguardandoAprovacao,
			},
		}

		// Mock get parts supplies by service order ID
		partsSupplyRepo.On("GetByServiceOrderID", context.Background(), uint(1)).Return([]entities.PartsSupply{
			{ID: 1, QuantityTotal: 10, QuantityReserve: 2},
		}, nil)

		// Mock get parts supply service order relation
		serviceOrderRepo.On("GetPartsSupplyServiceOrder", uint(1), uint(1)).Return(&dto.PartsSupplyServiceOrder{
			PartsSupplyID:  1,
			ServiceOrderID: 1,
			Quantity:       2,
		}, nil)

		// Mock get parts supply by ID
		partsSupplyRepo.On("GetByID", context.Background(), uint(1)).Return(&entities.PartsSupply{
			ID:              1,
			QuantityTotal:   10,
			QuantityReserve: 2,
		}, nil)

		// Mock update parts supply
		partsSupplyRepo.On("Update", context.Background(), mock.AnythingOfType("*entities.PartsSupply")).Return(nil)

		update := &entities.ServiceOrder{}
		current := toServiceOrderEntity(serviceOrderDTO)
		result, err := ValidateEstimate(context.Background(), request, current, update, partsSupplyRepo, serviceOrderRepo)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, valueobject.StatusAprovada, result.ServiceOrderStatus)
	})

	t.Run("Should fail when getting parts supply relation fails", func(t *testing.T) {
		partsSupplyRepo := new(mocks.MockPartsSupplyGateway)
		serviceOrderRepo := new(mocks.MockServiceOrderGateway)

		request := &entities.ServiceOrder{
			ID:                 1,
			ServiceOrderStatus: valueobject.StatusAprovada,
		}
		serviceOrderDTO := &dto.ServiceOrderModel{
			ID: 1,
			ServiceOrderStatus: dto.ServiceOrderStatus{
				Description: StatusAguardandoAprovacao,
			},
		}

		partsSupplyRepo.On("GetByServiceOrderID", context.Background(), uint(1)).Return([]entities.PartsSupply{
			{ID: 1},
		}, nil)

		serviceOrderRepo.On("GetPartsSupplyServiceOrder", uint(1), uint(1)).Return(nil, errors.New("error getting relation"))

		update := &entities.ServiceOrder{}
		current := toServiceOrderEntity(serviceOrderDTO)
		result, err := ValidateEstimate(context.Background(), request, current, update, partsSupplyRepo, serviceOrderRepo)

		assert.Error(t, err)
		assert.Equal(t, "error getting relation", err.Error())
		assert.Nil(t, result)
	})
}

func TestCalculateEstimate(t *testing.T) {
	t.Run("Calculate with services and parts supplies", func(t *testing.T) {
		services := []entities.Service{
			{ID: 1},
			{ID: 2},
		}
		partsSupplies := []entities.PartsSupply{
			{ID: 1, QuantityReserve: 2},
			{ID: 2, QuantityReserve: 3},
		}

		serviceRepo := new(mocks.MockServiceGateway)
		partsSupplyRepo := new(mocks.MockPartsSupplyGateway)

		// Setup mocks for services
		serviceRepo.On("GetByID", mock.Anything, uint(1)).Return(&entities.Service{
			ID:    1,
			Price: 100.0,
		}, nil)
		serviceRepo.On("GetByID", mock.Anything, uint(2)).Return(&entities.Service{
			ID:    2,
			Price: 75.0,
		}, nil)

		// Setup mocks for parts supplies
		partsSupplyRepo.On("GetByID", mock.Anything, uint(1)).Return(&entities.PartsSupply{
			ID:    1,
			Price: 50.0,
		}, nil)
		partsSupplyRepo.On("GetByID", mock.Anything, uint(2)).Return(&entities.PartsSupply{
			ID:    2,
			Price: 25.0,
		}, nil)

		result, err := CalculateEstimate(context.Background(), services, partsSupplies, serviceRepo, partsSupplyRepo)
		assert.NoError(t, err)
		assert.Equal(t, 350.0, result) // (100 + 75) + (50*2 + 25*3) = 175 + 175 = 350
	})
}

func TestValidateExecution(t *testing.T) {
	t.Run("Success - Start Execution", func(t *testing.T) {
		request := &entities.ServiceOrder{
			ID:                 1,
			ServiceOrderStatus: valueobject.StatusEmExecucao,
		}
		serviceOrderDTO := &dto.ServiceOrderModel{
			ID: 1,
			ServiceOrderStatus: dto.ServiceOrderStatus{
				Description: string(valueobject.StatusAprovada),
			},
		}

		update := &entities.ServiceOrder{}
		current := toServiceOrderEntity(serviceOrderDTO)
		result, err := ValidateExecution(context.Background(), request, current, update)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, valueobject.StatusEmExecucao, result.ServiceOrderStatus)
		assert.NotNil(t, result.StartedExecutionDate)
	})

	t.Run("Success - Finish Execution", func(t *testing.T) {
		request := &entities.ServiceOrder{
			ID:                 1,
			ServiceOrderStatus: valueobject.StatusFinalizada,
		}
		start := time.Now().Add(-75 * time.Minute) // 1.25 horas no passado
		serviceOrderDTO := &dto.ServiceOrderModel{
			ID: 1,
			ServiceOrderStatus: dto.ServiceOrderStatus{
				Description: string(valueobject.StatusEmExecucao),
			},
			StartedExecutionDate: &start,
		}

		update := &entities.ServiceOrder{}
		current := toServiceOrderEntity(serviceOrderDTO)
		result, err := ValidateExecution(context.Background(), request, current, update)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, valueobject.StatusFinalizada, result.ServiceOrderStatus)
		assert.InDelta(t, 1.25, result.ExecutionDurationInHours, 0.01)
	})

	t.Run("Error - Invalid Status Transition", func(t *testing.T) {
		request := &entities.ServiceOrder{
			ID:                 1,
			ServiceOrderStatus: valueobject.StatusFinalizada,
		}
		serviceOrderDTO := &dto.ServiceOrderModel{
			ID: 1,
			ServiceOrderStatus: dto.ServiceOrderStatus{
				Description: string(valueobject.StatusRecebida),
			},
		}

		update := &entities.ServiceOrder{}
		current := toServiceOrderEntity(serviceOrderDTO)
		result, err := ValidateExecution(context.Background(), request, current, update)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidTransitionStatusToExecution, err)
		assert.Nil(t, result)
	})
}

func TestValidateDelivery(t *testing.T) {
	t.Run("Success - Complete Delivery", func(t *testing.T) {
		request := &entities.ServiceOrder{
			ID:                 1,
			ServiceOrderStatus: valueobject.StatusEntregue,
		}
		serviceOrderDTO := &dto.ServiceOrderModel{
			ID: 1,
			ServiceOrderStatus: dto.ServiceOrderStatus{
				Description: string(valueobject.StatusFinalizada),
			},
			PaymentID: UintPointer(1),
		}

		update := &entities.ServiceOrder{}
		current := toServiceOrderEntity(serviceOrderDTO)
		result, err := ValidateDelivery(context.Background(), request, current, update)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, valueobject.StatusEntregue, result.ServiceOrderStatus)
	})

	t.Run("Error - Missing Payment Information", func(t *testing.T) {
		request := &entities.ServiceOrder{
			ID:                 1,
			ServiceOrderStatus: valueobject.StatusEntregue,
		}
		serviceOrderDTO := &dto.ServiceOrderModel{
			ID: 1,
			ServiceOrderStatus: dto.ServiceOrderStatus{
				Description: string(valueobject.StatusFinalizada),
			},
		}

		update := &entities.ServiceOrder{}
		current := toServiceOrderEntity(serviceOrderDTO)
		result, err := ValidateDelivery(context.Background(), request, current, update)

		assert.Error(t, err)
		assert.Equal(t, "payment information is required for delivery", err.Error())
		assert.Nil(t, result)
	})

	t.Run("Error - Invalid Status Transition", func(t *testing.T) {
		request := &entities.ServiceOrder{
			ID:                 1,
			ServiceOrderStatus: valueobject.StatusEntregue,
			PaymentID:          UintPointer(1),
		}
		serviceOrderDTO := &dto.ServiceOrderModel{
			ID: 1,
			ServiceOrderStatus: dto.ServiceOrderStatus{
				Description: string(valueobject.StatusEmDiagnostico),
			},
		}

		update := &entities.ServiceOrder{}
		current := toServiceOrderEntity(serviceOrderDTO)
		result, err := ValidateDelivery(context.Background(), request, current, update)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidTransitionStatusToDelivery, err)
		assert.Nil(t, result)
	})
}

func TestInvalidServiceOrder(t *testing.T) {
	t.Run("Error - Service Order Not Found", func(t *testing.T) {
		vehicleRepo := new(mocks.MockVehicleGateway)
		customerRepo := new(mocks.MockCustomerGateway)
		serviceOrderRepo := new(mocks.MockServiceOrderGateway)
		serviceRepo := new(mocks.MockServiceGateway)
		partsSupplyRepo := new(mocks.MockPartsSupplyGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo)

		serviceOrder := entities.ServiceOrder{
			ID:                 999,
			ServiceOrderStatus: valueobject.StatusEmDiagnostico,
		}
		flow := DIAGNOSIS

		serviceOrderRepo.On("GetByID", uint(999), false).Return(nil, ErrServiceOrderNotFound)

		r, err := useCase.UpdateServiceOrder(context.Background(), serviceOrder, flow)
		assert.Error(t, err)
		assert.Equal(t, ErrServiceOrderNotFound, err)
		assert.Nil(t, r)
	})

	t.Run("Error - Invalid Flow Type", func(t *testing.T) {
		vehicleRepo := new(mocks.MockVehicleGateway)
		customerRepo := new(mocks.MockCustomerGateway)
		serviceOrderRepo := new(mocks.MockServiceOrderGateway)
		serviceRepo := new(mocks.MockServiceGateway)
		partsSupplyRepo := new(mocks.MockPartsSupplyGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo)

		serviceOrder := entities.ServiceOrder{
			ID:                 1,
			ServiceOrderStatus: valueobject.StatusEmDiagnostico,
		}
		flow := "invalid_flow"

		serviceOrderRepo.On("GetByID", uint(1), false).Return(&dto.ServiceOrderModel{
			ID: 1,
			ServiceOrderStatus: dto.ServiceOrderStatus{
				Description: string(valueobject.StatusRecebida),
			},
		}, nil)

		r, err := useCase.UpdateServiceOrder(context.Background(), serviceOrder, flow)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidFlow, err)
		assert.Nil(t, r)
	})
}

func TestGetServiceOrder(t *testing.T) {
	serviceOrderRepo := new(mocks.MockServiceOrderGateway)
	vehicleRepo := new(mocks.MockVehicleGateway)
	customerRepo := new(mocks.MockCustomerGateway)
	serviceRepo := new(mocks.MockServiceGateway)
	partsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo)

	ctx := context.Background()
	validID := uint(1)
	invalidID := uint(999)
	serviceOrderEntity := entities.ServiceOrder{ID: validID}
	serviceOrderDTO := &dto.ServiceOrderModel{ID: validID}

	t.Run("success", func(t *testing.T) {
		serviceOrderRepo.On("GetByID", validID, false).Return(serviceOrderDTO, nil)
		result, err := useCase.GetServiceOrder(ctx, serviceOrderEntity, false)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		serviceOrderRepo.AssertCalled(t, "GetByID", validID, false)
	})

	t.Run("not found", func(t *testing.T) {
		serviceOrderRepo.On("GetByID", invalidID, false).Return(nil, nil)
		result, err := useCase.GetServiceOrder(ctx, entities.ServiceOrder{ID: invalidID}, false)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrServiceOrderNotFound, err)
		serviceOrderRepo.AssertCalled(t, "GetByID", invalidID, false)
	})
}

func TestListServiceOrders(t *testing.T) {
	serviceOrderRepo := new(mocks.MockServiceOrderGateway)
	vehicleRepo := new(mocks.MockVehicleGateway)
	customerRepo := new(mocks.MockCustomerGateway)
	serviceRepo := new(mocks.MockServiceGateway)
	partsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo)

	ctx := context.Background()
	base := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	timePtr := func(offsetDays int) *time.Time {
		ts := base.Add(time.Duration(offsetDays) * 24 * time.Hour)
		return &ts
	}

	serviceOrderDTOs := []dto.ServiceOrderModel{
		{ID: 1, CreatedAt: timePtr(0), ServiceOrderStatus: dto.ServiceOrderStatus{Description: string(valueobject.StatusEmExecucao)}},
		{ID: 5, CreatedAt: timePtr(2), ServiceOrderStatus: dto.ServiceOrderStatus{Description: string(valueobject.StatusEmExecucao)}},
		{ID: 2, CreatedAt: timePtr(1), ServiceOrderStatus: dto.ServiceOrderStatus{Description: string(valueobject.StatusAguardandoAprovacao)}},
		{ID: 3, CreatedAt: timePtr(3), ServiceOrderStatus: dto.ServiceOrderStatus{Description: string(valueobject.StatusEmDiagnostico)}},
		{ID: 4, CreatedAt: timePtr(4), ServiceOrderStatus: dto.ServiceOrderStatus{Description: string(valueobject.StatusRecebida)}},
		{ID: 6, CreatedAt: timePtr(5), ServiceOrderStatus: dto.ServiceOrderStatus{Description: string(valueobject.StatusCancelada)}},
		{ID: 7, CreatedAt: timePtr(6), ServiceOrderStatus: dto.ServiceOrderStatus{Description: string(valueobject.StatusFinalizada)}},
		{ID: 8, CreatedAt: timePtr(7), ServiceOrderStatus: dto.ServiceOrderStatus{Description: string(valueobject.StatusEntregue)}},
	}

	t.Run("apply ordering rules", func(t *testing.T) {
		serviceOrderRepo.On("List").Return(serviceOrderDTOs, nil)
		result, err := useCase.ListServiceOrders(ctx)

		assert.NoError(t, err)
		// Finalizada and Entregue must be excluded from the listing.
		assert.Len(t, result, 6)

		expectedOrder := []uint{1, 5, 2, 3, 4, 6}
		for idx, expectedID := range expectedOrder {
			assert.Equal(t, expectedID, result[idx].ID)
		}

		// Ensure no filtered statuses leak through.
		for _, so := range result {
			assert.False(t, so.ServiceOrderStatus.IsFinalizada())
			assert.False(t, so.ServiceOrderStatus.IsEntregue())
		}

		serviceOrderRepo.AssertExpectations(t)
	})
}
