package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/constants"
	mocks "github.com/fiap-grupo95/os-service-api/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceOrderUseCase_DeliveryServiceOrder_RepoUpdateFails(t *testing.T) {
	serviceOrderRepo := new(mocks.MockServiceOrderGateway)
	vehicleRepo := new(mocks.MockVehicleGateway)
	customerRepo := new(mocks.MockCustomerGateway)
	serviceRepo := new(mocks.MockServiceGateway)
	partsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	billingRepo := new(mocks.MockBillingServiceGateway)
	executionRepo := new(mocks.MockExecutionGateway)

	uc := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingRepo, executionRepo)

	serviceOrderRepo.On("GetByID", mock.Anything, "1", false).Return(&entities.ServiceOrder{ID: "1", Status: valueobject.StatusFinalizada, Estimate: &entities.Estimate{ID: "est-1"}}, nil)
	billingRepo.On("GetPaymentByEstimateID", mock.Anything, "est-1").Return(&entities.Payment{ID: "pay-1", EstimateID: "est-1"}, nil)
	serviceOrderRepo.On("Update", mock.Anything, mock.AnythingOfType("*entities.ServiceOrder")).Return(nil, errors.New("update error"))

	so, err := uc.DeliveryServiceOrder(context.Background(), "1")
	assert.Error(t, err)
	assert.Equal(t, "update error", err.Error())
	assert.Nil(t, so)
}

func TestServiceOrderUseCase_ExecutionServiceOrder_UpdateFails(t *testing.T) {
	serviceOrderRepo := new(mocks.MockServiceOrderGateway)
	vehicleRepo := new(mocks.MockVehicleGateway)
	customerRepo := new(mocks.MockCustomerGateway)
	serviceRepo := new(mocks.MockServiceGateway)
	partsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	billingRepo := new(mocks.MockBillingServiceGateway)
	executionRepo := new(mocks.MockExecutionGateway)

	uc := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingRepo, executionRepo)

	serviceOrderRepo.On("GetByID", mock.Anything, "1", false).Return(&entities.ServiceOrder{ID: "1", Status: valueobject.StatusAprovada}, nil)
	executionRepo.On("CreateExecution", mock.Anything, mock.Anything).Return(&entities.Execution{ID: "exec-1"}, nil)
	serviceOrderRepo.On("Update", mock.Anything, mock.AnythingOfType("*entities.ServiceOrder")).Return(nil, errors.New("update error"))

	so, err := uc.ExecutionServiceOrder(context.Background(), "1", constants.EXECUTION_START)
	assert.Error(t, err)
	assert.Equal(t, "update error", err.Error())
	assert.Nil(t, so)
}

func TestServiceOrderUseCase_DiagnosisServiceOrder_AuthorizeReserveFails(t *testing.T) {
	serviceOrderRepo := new(mocks.MockServiceOrderGateway)
	vehicleRepo := new(mocks.MockVehicleGateway)
	customerRepo := new(mocks.MockCustomerGateway)
	serviceRepo := new(mocks.MockServiceGateway)
	partsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	billingRepo := new(mocks.MockBillingServiceGateway)
	executionRepo := new(mocks.MockExecutionGateway)

	uc := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingRepo, executionRepo)

	req := &entities.ServiceOrder{ID: "1", Services: []entities.Service{{ID: "s1"}}, PartsSupplies: []entities.PartsSupply{{ID: "p1", Quantity: 1}}}
	serviceOrderRepo.On("GetByID", mock.Anything, "1", false).Return(&entities.ServiceOrder{ID: "1", Status: valueobject.StatusRecebida}, nil)
	serviceRepo.On("GetByID", mock.Anything, "s1").Return(&entities.Service{ID: "s1"}, nil)
	partsSupplyRepo.On("AuthorizeReserve", mock.Anything, mock.Anything).Return(errors.New("auth error"))

	so, err := uc.DiagnosisServiceOrder(context.Background(), req)
	assert.Error(t, err)
	assert.Equal(t, "auth error", err.Error())
	assert.Nil(t, so)
}

func TestServiceOrderUseCase_EstimateServiceOrder_StrategyExecuteFails(t *testing.T) {
	serviceOrderRepo := new(mocks.MockServiceOrderGateway)
	vehicleRepo := new(mocks.MockVehicleGateway)
	customerRepo := new(mocks.MockCustomerGateway)
	serviceRepo := new(mocks.MockServiceGateway)
	partsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	billingRepo := new(mocks.MockBillingServiceGateway)
	executionRepo := new(mocks.MockExecutionGateway)

	uc := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingRepo, executionRepo)

	serviceOrderRepo.On("GetByID", mock.Anything, "1", false).Return(&entities.ServiceOrder{ID: "1", Status: valueobject.StatusAguardandoAprovacao, PartsSupplies: []entities.PartsSupply{{ID: "p1", Quantity: 1}}}, nil)
	// Approve path calls WriteOff first
	partsSupplyRepo.On("WriteOff", mock.Anything, mock.Anything).Return(errors.New("writeoff error"))

	so, err := uc.EstimateServiceOrder(context.Background(), "1", constants.ESTIMATE_APPROVE)
	assert.Error(t, err)
	assert.Equal(t, "writeoff error", err.Error())
	assert.Nil(t, so)
}
