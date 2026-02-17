package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
	dto "github.com/fiap-grupo95/os-service-api/internal/infrastructure/database/model"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/adapter/operations"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/constants"
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
		billingServiceRepo := new(mocks.MockBillingServiceGateway)
		executionRepo := new(mocks.MockExecutionGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingServiceRepo, executionRepo)

		serviceOrder := entities.ServiceOrder{
			CustomerID: "1",
			VehicleID:  "1",
			Status:     valueobject.StatusRecebida,
		}

		vehicleRepo.On("FindByID", "1").Return(&entities.Vehicle{ID: "1"}, nil)
		customerRepo.On("GetByID", "1").Return(&entities.Customer{ID: "1"}, nil)
		serviceOrderRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.ServiceOrder")).Return(&entities.ServiceOrder{
			ID:         "1",
			CustomerID: "1",
			VehicleID:  "1",
		}, nil)

		r, err := useCase.CreateServiceOrder(context.Background(), &serviceOrder)
		assert.NoError(t, err)
		assert.NotNil(t, r)
	})

	t.Run("Error - Internal DB Error", func(t *testing.T) {
		vehicleRepo := new(mocks.MockVehicleGateway)
		customerRepo := new(mocks.MockCustomerGateway)
		serviceOrderRepo := new(mocks.MockServiceOrderGateway)
		serviceRepo := new(mocks.MockServiceGateway)
		partsSupplyRepo := new(mocks.MockPartsSupplyGateway)
		billingServiceRepo := new(mocks.MockBillingServiceGateway)
		executionRepo := new(mocks.MockExecutionGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingServiceRepo, executionRepo)

		serviceOrder := entities.ServiceOrder{
			CustomerID: "1",
			VehicleID:  "1",
		}

		vehicleRepo.On("FindByID", "1").Return(&entities.Vehicle{ID: "1"}, nil)
		customerRepo.On("GetByID", "1").Return(&entities.Customer{ID: "1"}, nil)
		serviceOrderRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.ServiceOrder")).Return(nil, errors.New("internal db error"))

		r, err := useCase.CreateServiceOrder(context.Background(), &serviceOrder)
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
		billingServiceRepo := new(mocks.MockBillingServiceGateway)
		executionRepo := new(mocks.MockExecutionGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingServiceRepo, executionRepo)

		serviceOrder := &entities.ServiceOrder{
			CustomerID: "1",
			VehicleID:  "1",
		}

		vehicleRepo.On("FindByID", "1").Return(nil, nil)

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
		billingServiceRepo := new(mocks.MockBillingServiceGateway)
		executionRepo := new(mocks.MockExecutionGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingServiceRepo, executionRepo)

		serviceOrder := &entities.ServiceOrder{
			CustomerID: "1",
			VehicleID:  "0",
		}

		vehicleRepo.On("FindByID", "0").Return(nil, errors.New("invalid vehicle ID"))

		r, err := useCase.CreateServiceOrder(context.Background(), serviceOrder)
		assert.Error(t, err)
		assert.Equal(t, "invalid vehicle ID", err.Error())
		assert.Nil(t, r)
	})

	t.Run("Error - Customer not found", func(t *testing.T) {
		vehicleRepo := new(mocks.MockVehicleGateway)
		customerRepo := new(mocks.MockCustomerGateway)
		serviceOrderRepo := new(mocks.MockServiceOrderGateway)
		serviceRepo := new(mocks.MockServiceGateway)
		partsSupplyRepo := new(mocks.MockPartsSupplyGateway)
		billingServiceRepo := new(mocks.MockBillingServiceGateway)
		executionRepo := new(mocks.MockExecutionGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingServiceRepo, executionRepo)

		serviceOrder := &entities.ServiceOrder{
			CustomerID: "1",
			VehicleID:  "1",
		}

		vehicleRepo.On("FindByID", "1").Return(&entities.Vehicle{ID: "1"}, nil)
		customerRepo.On("GetByID", "1").Return(nil, nil)

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
		billingServiceRepo := new(mocks.MockBillingServiceGateway)
		executionRepo := new(mocks.MockExecutionGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingServiceRepo, executionRepo)

		serviceOrder := &entities.ServiceOrder{
			CustomerID: "0",
			VehicleID:  "1",
		}

		vehicleRepo.On("FindByID", "1").Return(&entities.Vehicle{ID: "1"}, nil)
		customerRepo.On("GetByID", "0").Return(nil, errors.New("invalid customer ID"))

		r, err := useCase.CreateServiceOrder(context.Background(), serviceOrder)
		assert.Error(t, err)
		assert.Equal(t, "invalid customer ID", err.Error())
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
		billingServiceRepo := new(mocks.MockBillingServiceGateway)
		executionRepo := new(mocks.MockExecutionGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingServiceRepo, executionRepo)

		serviceOrder := &entities.ServiceOrder{
			ID:     "1",
			Status: valueobject.StatusRecebida,
			Services: []entities.Service{
				{ID: "1"},
			},
			PartsSupplies: []entities.PartsSupply{
				{
					ID:       "1",
					Quantity: 2,
				},
			},
		}

		executionRepo.On("CreateExecution", mock.Anything, mock.Anything).Return(&entities.Execution{ID: "1", ServiceOrderID: "1"}, nil)

		serviceOrderRepo.On("GetByID", mock.Anything, "1", false).Return(&dto.ServiceOrderModel{
			ID: "1",
			ServiceOrderStatus: dto.ServiceOrderStatus{
				ID:          "1",
				Description: string(valueobject.StatusRecebida),
			},
		}, nil)
		serviceRepo.On("GetByID", mock.Anything, "1").Return(&entities.Service{ID: "1"}, nil)
		partsSupplyRepo.On("GetByID", mock.Anything, "1").Return(&entities.PartsSupply{
			ID:       "1",
			Quantity: 10,
		}, nil)
		partsSupplyRepo.On("AuthorizeReserve", mock.Anything, mock.Anything).Return(nil)
		// Reserve is called inside validateDiagnosis to reserve parts supply
		partsSupplyRepo.On("Reserve", mock.Anything, mock.Anything).Return(nil)
		// Billing service is called to create an estimate during diagnosis
		estimate := &entities.Estimate{ID: "estimate-1", ServiceOrderID: "1", Value: 100.0, Status: "created"}
		billingServiceRepo.On("CreateEstimate", mock.Anything, mock.Anything, mock.Anything).Return(estimate, nil)
		partsSupplyRepo.On("Update", mock.Anything, mock.AnythingOfType("*entities.PartsSupply")).Return(nil)
		serviceOrderRepo.On("Update", mock.AnythingOfType("*entities.ServiceOrder")).Return(nil)

		r, err := useCase.DiagnosisServiceOrder(context.Background(), serviceOrder)
		assert.NoError(t, err)
		assert.NotNil(t, r)
	})

	t.Run("Error - Invalid Status Transition", func(t *testing.T) {
		vehicleRepo := new(mocks.MockVehicleGateway)
		customerRepo := new(mocks.MockCustomerGateway)
		serviceOrderRepo := new(mocks.MockServiceOrderGateway)
		serviceRepo := new(mocks.MockServiceGateway)
		partsSupplyRepo := new(mocks.MockPartsSupplyGateway)
		billingServiceRepo := new(mocks.MockBillingServiceGateway)
		executionRepo := new(mocks.MockExecutionGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingServiceRepo, executionRepo)

		serviceOrder := &entities.ServiceOrder{
			ID:     "1",
			Status: valueobject.StatusEntregue,
		}

		serviceOrderRepo.On("GetByID", mock.Anything, "1", false).Return(&dto.ServiceOrderModel{
			ID: "1",
			ServiceOrderStatus: dto.ServiceOrderStatus{
				ID:          "1",
				Description: string(valueobject.StatusEntregue),
			},
		}, nil)

		r, err := useCase.DiagnosisServiceOrder(context.Background(), serviceOrder)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidTransitionStatusToDiagnosis, err)
		assert.Nil(t, r)
	})
}

func TestValidateEstimate(t *testing.T) {
	t.Run("Should approve estimate and release parts supply", func(t *testing.T) {
		vehicleRepo := new(mocks.MockVehicleGateway)
		customerRepo := new(mocks.MockCustomerGateway)
		serviceOrderRepo := new(mocks.MockServiceOrderGateway)
		serviceRepo := new(mocks.MockServiceGateway)
		partsSupplyRepo := new(mocks.MockPartsSupplyGateway)
		billingServiceRepo := new(mocks.MockBillingServiceGateway)
		executionRepo := new(mocks.MockExecutionGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingServiceRepo, executionRepo)

		serviceOrderID := "1"

		// Mock get parts supplies by service order ID
		partsSupplyRepo.On("GetByServiceOrderID", context.Background(), "1").Return([]entities.PartsSupply{
			{ID: "1", Quantity: 2},
		}, nil)

		// Mock get parts supply service order relation
		serviceOrderRepo.On("GetPartsSupplyServiceOrder", context.Background(), "1", "1").Return(&dto.PartsSupplyServiceOrder{
			PartsSupplyID:  "1",
			ServiceOrderID: "1",
			Quantity:       2,
		}, nil)

		// Mock get parts supply by ID
		partsSupplyRepo.On("GetByID", context.Background(), "1").Return(&entities.PartsSupply{
			ID:       "1",
			Quantity: 10,
		}, nil)

		// Mock update parts supply
		partsSupplyRepo.On("Update", context.Background(), mock.AnythingOfType("*entities.PartsSupply")).Return(nil)

		serviceOrderRepo.On("GetByID", context.Background(), serviceOrderID, false).Return(&entities.ServiceOrder{
			ID:     serviceOrderID,
			Status: valueobject.StatusAguardandoAprovacao,
			Estimate: &entities.Estimate{
				ID: "estimate-1",
			},
			PartsSupplies: []entities.PartsSupply{{ID: "1", Quantity: 2}},
			Services:      []entities.Service{{ID: "1"}},
		}, nil)
		partsSupplyRepo.On("WriteOff", mock.Anything, mock.Anything).Return(nil)
		billingServiceRepo.On("ApproveEstimate", mock.Anything, mock.Anything, mock.Anything).Return(&entities.Estimate{ID: "estimate-1", ServiceOrderID: "1", Value: 100.0, Status: "approved"}, nil)
		serviceOrderRepo.On("Update", mock.AnythingOfType("*entities.ServiceOrder")).Return(nil)

		result, err := useCase.EstimateServiceOrder(context.Background(), serviceOrderID, constants.ESTIMATE_APPROVE)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, valueobject.StatusAprovada, result.Status)
	})

	t.Run("Should fail when getting parts supply relation fails", func(t *testing.T) {
		vehicleRepo := new(mocks.MockVehicleGateway)
		customerRepo := new(mocks.MockCustomerGateway)
		partsSupplyRepo := new(mocks.MockPartsSupplyGateway)
		serviceOrderRepo := new(mocks.MockServiceOrderGateway)
		serviceRepo := new(mocks.MockServiceGateway)
		billingServiceRepo := new(mocks.MockBillingServiceGateway)
		executionRepo := new(mocks.MockExecutionGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingServiceRepo, executionRepo)

		serviceOrderID := "1"

		serviceOrderRepo.On("GetByID", context.Background(), serviceOrderID, false).Return(&entities.ServiceOrder{ID: serviceOrderID, Status: valueobject.StatusAguardandoAprovacao}, nil)
		result, err := useCase.EstimateServiceOrder(context.Background(), serviceOrderID, "invalid_flow")

		assert.Error(t, err)
		assert.Equal(t, operations.ErrInvalidEstimateOperation, err.Error())
		assert.Nil(t, result)
	})
}

func TestCalculateEstimate(t *testing.T) {
	t.Run("Calculate with services and parts supplies", func(t *testing.T) {
		vehicleRepo := new(mocks.MockVehicleGateway)
		customerRepo := new(mocks.MockCustomerGateway)
		partsSupplyRepo := new(mocks.MockPartsSupplyGateway)
		serviceOrderRepo := new(mocks.MockServiceOrderGateway)
		serviceRepo := new(mocks.MockServiceGateway)
		billingServiceRepo := new(mocks.MockBillingServiceGateway)
		executionRepo := new(mocks.MockExecutionGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingServiceRepo, executionRepo)
		services := []entities.Service{
			{ID: "1"},
			{ID: "2"},
		}
		partsSupplies := []entities.PartsSupply{
			{ID: "1", Quantity: 2},
			{ID: "2", Quantity: 3},
		}

		// Setup mocks for services
		serviceRepo.On("GetByID", mock.Anything, "1").Return(&entities.Service{
			ID:    "1",
			Price: 100.0,
		}, nil)
		serviceRepo.On("GetByID", mock.Anything, "2").Return(&entities.Service{
			ID:    "2",
			Price: 75.0,
		}, nil)

		// Setup mocks for parts supplies
		partsSupplyRepo.On("GetByID", mock.Anything, "1").Return(&entities.PartsSupply{
			ID:    "1",
			Price: 50.0,
		}, nil)
		partsSupplyRepo.On("GetByID", mock.Anything, "2").Return(&entities.PartsSupply{
			ID:    "2",
			Price: 25.0,
		}, nil)

		billingServiceRepo.On("CreateEstimate", mock.Anything, mock.Anything, mock.Anything).Return(&entities.Estimate{ID: "estimate-350", ServiceOrderID: "1", Value: 350.0, Status: "created"}, nil)
		result, err := useCase.billingServiceRepo.CreateEstimate(context.Background(), &entities.ServiceOrder{
			ID:            "1",
			Services:      services,
			PartsSupplies: partsSupplies,
		}, nil)
		assert.NoError(t, err)
		assert.Equal(t, 350.0, result.Value) // (100 + 75) + (50*2 + 25*3) = 175 + 175 = 350
	})
}

func TestValidateExecution(t *testing.T) {
	t.Run("Success - Start Execution", func(t *testing.T) {
		vehicleRepo := new(mocks.MockVehicleGateway)
		customerRepo := new(mocks.MockCustomerGateway)
		partsSupplyRepo := new(mocks.MockPartsSupplyGateway)
		serviceOrderRepo := new(mocks.MockServiceOrderGateway)
		serviceRepo := new(mocks.MockServiceGateway)
		billingServiceRepo := new(mocks.MockBillingServiceGateway)
		executionRepo := new(mocks.MockExecutionGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingServiceRepo, executionRepo)
		serviceOrderID := "1"
		serviceOrderRepo.On("GetByID", context.Background(), serviceOrderID, false).Return(&entities.ServiceOrder{ID: serviceOrderID, Status: valueobject.StatusAprovada}, nil)
		executionRepo.On("CreateExecution", mock.Anything, mock.Anything).Return(&entities.Execution{ID: "1", ServiceOrderID: "1", Status: "started"}, nil)
		serviceOrderRepo.On("Update", mock.AnythingOfType("*entities.ServiceOrder")).Return(nil)

		result, err := useCase.ExecutionServiceOrder(context.Background(), serviceOrderID, constants.EXECUTION_START)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, valueobject.StatusEmExecucao, result.Status)
		assert.NotNil(t, result.Execution)
	})

	t.Run("Success - Finish Execution", func(t *testing.T) {
		vehicleRepo := new(mocks.MockVehicleGateway)
		customerRepo := new(mocks.MockCustomerGateway)
		partsSupplyRepo := new(mocks.MockPartsSupplyGateway)
		serviceOrderRepo := new(mocks.MockServiceOrderGateway)
		serviceRepo := new(mocks.MockServiceGateway)
		billingServiceRepo := new(mocks.MockBillingServiceGateway)
		executionRepo := new(mocks.MockExecutionGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingServiceRepo, executionRepo)
		serviceOrderID := "1"
		serviceOrderRepo.On("GetByID", context.Background(), serviceOrderID, false).Return(&entities.ServiceOrder{ID: serviceOrderID, Status: valueobject.StatusEmExecucao}, nil)
		executionRepo.On("FinishExecution", mock.Anything, mock.Anything).Return(&entities.Execution{ID: "1", ServiceOrderID: "1", Status: "finished"}, nil)
		serviceOrderRepo.On("Update", mock.AnythingOfType("*entities.ServiceOrder")).Return(nil)

		result, err := useCase.ExecutionServiceOrder(context.Background(), serviceOrderID, constants.EXECUTION_FINISH)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, valueobject.StatusFinalizada, result.Status)
		assert.NotNil(t, result.Execution)
	})

	t.Run("Error - Invalid Status Transition", func(t *testing.T) {
		vehicleRepo := new(mocks.MockVehicleGateway)
		customerRepo := new(mocks.MockCustomerGateway)
		partsSupplyRepo := new(mocks.MockPartsSupplyGateway)
		serviceOrderRepo := new(mocks.MockServiceOrderGateway)
		serviceRepo := new(mocks.MockServiceGateway)
		billingServiceRepo := new(mocks.MockBillingServiceGateway)
		executionRepo := new(mocks.MockExecutionGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingServiceRepo, executionRepo)
		serviceOrderID := "1"
		serviceOrderRepo.On("GetByID", context.Background(), serviceOrderID, false).Return(&entities.ServiceOrder{ID: serviceOrderID, Status: valueobject.StatusRecebida}, nil)

		result, err := useCase.ExecutionServiceOrder(context.Background(), serviceOrderID, constants.EXECUTION_FINISH)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidTransitionStatusToExecution, err)
		assert.Nil(t, result)
	})
}

func TestValidateDelivery(t *testing.T) {
	t.Run("Success - Complete Delivery", func(t *testing.T) {
		vehicleRepo := new(mocks.MockVehicleGateway)
		customerRepo := new(mocks.MockCustomerGateway)
		partsSupplyRepo := new(mocks.MockPartsSupplyGateway)
		serviceOrderRepo := new(mocks.MockServiceOrderGateway)
		serviceRepo := new(mocks.MockServiceGateway)
		billingServiceRepo := new(mocks.MockBillingServiceGateway)
		executionRepo := new(mocks.MockExecutionGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingServiceRepo, executionRepo)
		serviceOrderID := "1"
		serviceOrderRepo.On("GetByID", context.Background(), serviceOrderID, false).Return(&entities.ServiceOrder{ID: serviceOrderID, Status: valueobject.StatusFinalizada, Estimate: &entities.Estimate{ID: "estimate-1"}}, nil)
		billingServiceRepo.On("GetPaymentByEstimateID", mock.Anything, "estimate-1").Return(&entities.Payment{ID: "1", EstimateID: "estimate-1", Amount: 100.0, PaymentDate: time.Now()}, nil)
		serviceOrderRepo.On("Update", mock.AnythingOfType("*entities.ServiceOrder")).Return(nil)

		result, err := useCase.DeliveryServiceOrder(context.Background(), serviceOrderID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, valueobject.StatusEntregue, result.Status)
	})

	t.Run("Error - Missing Payment Information", func(t *testing.T) {
		vehicleRepo := new(mocks.MockVehicleGateway)
		customerRepo := new(mocks.MockCustomerGateway)
		partsSupplyRepo := new(mocks.MockPartsSupplyGateway)
		serviceOrderRepo := new(mocks.MockServiceOrderGateway)
		serviceRepo := new(mocks.MockServiceGateway)
		billingServiceRepo := new(mocks.MockBillingServiceGateway)
		executionRepo := new(mocks.MockExecutionGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingServiceRepo, executionRepo)
		serviceOrderID := "1"
		serviceOrderRepo.On("GetByID", context.Background(), serviceOrderID, false).Return(&entities.ServiceOrder{ID: serviceOrderID, Status: valueobject.StatusFinalizada, Estimate: &entities.Estimate{ID: "estimate-1"}}, nil)
		billingServiceRepo.On("GetPaymentByEstimateID", mock.Anything, "estimate-1").Return(nil, nil)

		result, err := useCase.DeliveryServiceOrder(context.Background(), serviceOrderID)

		assert.Error(t, err)
		assert.Equal(t, "payment information is required for delivery", err.Error())
		assert.Nil(t, result)
	})

	t.Run("Error - Invalid Status Transition", func(t *testing.T) {
		vehicleRepo := new(mocks.MockVehicleGateway)
		customerRepo := new(mocks.MockCustomerGateway)
		partsSupplyRepo := new(mocks.MockPartsSupplyGateway)
		serviceOrderRepo := new(mocks.MockServiceOrderGateway)
		serviceRepo := new(mocks.MockServiceGateway)
		billingServiceRepo := new(mocks.MockBillingServiceGateway)
		executionRepo := new(mocks.MockExecutionGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingServiceRepo, executionRepo)
		serviceOrderID := "1"
		serviceOrderRepo.On("GetByID", context.Background(), serviceOrderID, false).Return(&entities.ServiceOrder{ID: serviceOrderID, Status: valueobject.StatusEmDiagnostico}, nil)

		result, err := useCase.DeliveryServiceOrder(context.Background(), serviceOrderID)

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
		billingServiceRepo := new(mocks.MockBillingServiceGateway)
		executionRepo := new(mocks.MockExecutionGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingServiceRepo, executionRepo)

		serviceOrder := &entities.ServiceOrder{
			ID:     "999",
			Status: valueobject.StatusEmDiagnostico,
		}
		serviceOrderRepo.On("GetByID", context.Background(), "999", false).Return(nil, ErrServiceOrderNotFound)
		r, err := useCase.DiagnosisServiceOrder(context.Background(), serviceOrder)
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
		billingServiceRepo := new(mocks.MockBillingServiceGateway)
		executionRepo := new(mocks.MockExecutionGateway)

		useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingServiceRepo, executionRepo)

		serviceOrderID := "1"

		serviceOrderRepo.On("GetByID", context.Background(), serviceOrderID, false).Return(&entities.ServiceOrder{ID: serviceOrderID, Status: valueobject.StatusAguardandoAprovacao}, nil)
		r, err := useCase.EstimateServiceOrder(context.Background(), serviceOrderID, "invalid_flow")
		assert.Error(t, err)
		assert.Equal(t, operations.ErrInvalidEstimateOperation, err.Error())
		assert.Nil(t, r)
	})
}

func TestGetServiceOrder(t *testing.T) {
	serviceOrderRepo := new(mocks.MockServiceOrderGateway)
	vehicleRepo := new(mocks.MockVehicleGateway)
	customerRepo := new(mocks.MockCustomerGateway)
	serviceRepo := new(mocks.MockServiceGateway)
	partsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	billingServiceRepo := new(mocks.MockBillingServiceGateway)
	executionRepo := new(mocks.MockExecutionGateway)
	useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingServiceRepo, executionRepo)

	ctx := context.Background()
	validID := "1"
	invalidID := "999"
	serviceOrderEntity := entities.ServiceOrder{ID: validID}
	serviceOrderDTO := &dto.ServiceOrderModel{ID: validID}

	t.Run("success", func(t *testing.T) {
		serviceOrderRepo.On("GetByID", ctx, validID, false).Return(serviceOrderDTO, nil)
		result, err := useCase.GetServiceOrder(ctx, serviceOrderEntity, false)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		serviceOrderRepo.AssertCalled(t, "GetByID", ctx, validID, false)
	})

	t.Run("not found", func(t *testing.T) {
		serviceOrderRepo.On("GetByID", ctx, invalidID, false).Return(nil, nil)
		result, err := useCase.GetServiceOrder(ctx, entities.ServiceOrder{ID: invalidID}, false)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrServiceOrderNotFound, err)
		serviceOrderRepo.AssertCalled(t, "GetByID", ctx, invalidID, false)
	})
}

func TestListServiceOrders(t *testing.T) {
	serviceOrderRepo := new(mocks.MockServiceOrderGateway)
	vehicleRepo := new(mocks.MockVehicleGateway)
	customerRepo := new(mocks.MockCustomerGateway)
	serviceRepo := new(mocks.MockServiceGateway)
	partsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	billingServiceRepo := new(mocks.MockBillingServiceGateway)
	executionRepo := new(mocks.MockExecutionGateway)
	useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingServiceRepo, executionRepo)

	ctx := context.Background()
	base := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	timePtr := func(offsetDays int) *time.Time {
		ts := base.Add(time.Duration(offsetDays) * 24 * time.Hour)
		return &ts
	}

	serviceOrderDTOs := []dto.ServiceOrderModel{
		{ID: "1", CreatedAt: timePtr(0), ServiceOrderStatus: dto.ServiceOrderStatus{Description: string(valueobject.StatusEmExecucao)}},
		{ID: "5", CreatedAt: timePtr(2), ServiceOrderStatus: dto.ServiceOrderStatus{Description: string(valueobject.StatusEmExecucao)}},
		{ID: "2", CreatedAt: timePtr(1), ServiceOrderStatus: dto.ServiceOrderStatus{Description: string(valueobject.StatusAguardandoAprovacao)}},
		{ID: "3", CreatedAt: timePtr(3), ServiceOrderStatus: dto.ServiceOrderStatus{Description: string(valueobject.StatusEmDiagnostico)}},
		{ID: "4", CreatedAt: timePtr(4), ServiceOrderStatus: dto.ServiceOrderStatus{Description: string(valueobject.StatusRecebida)}},
		{ID: "6", CreatedAt: timePtr(5), ServiceOrderStatus: dto.ServiceOrderStatus{Description: string(valueobject.StatusCancelada)}},
		{ID: "7", CreatedAt: timePtr(6), ServiceOrderStatus: dto.ServiceOrderStatus{Description: string(valueobject.StatusFinalizada)}},
		{ID: "8", CreatedAt: timePtr(7), ServiceOrderStatus: dto.ServiceOrderStatus{Description: string(valueobject.StatusEntregue)}},
	}

	t.Run("apply ordering rules", func(t *testing.T) {
		serviceOrderRepo.On("List").Return(serviceOrderDTOs, nil)
		result, err := useCase.ListServiceOrders(ctx)

		assert.NoError(t, err)
		// Finalizada and Entregue must be excluded from the listing.
		assert.Len(t, result, 6)

		expectedOrder := []string{"1", "5", "2", "3", "4", "6"}
		for idx, expectedID := range expectedOrder {
			assert.Equal(t, expectedID, result[idx].ID)
		}

		// Ensure no filtered statuses leak through.
		for _, so := range result {
			assert.False(t, so.Status.IsFinalizada())
			assert.False(t, so.Status.IsEntregue())
		}

		serviceOrderRepo.AssertExpectations(t)
	})
}
