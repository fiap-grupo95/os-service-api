package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
	mocks "github.com/fiap-grupo95/os-service-api/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCancelServiceOrder_GetByIDFails(t *testing.T) {
	serviceOrderRepo := new(mocks.MockServiceOrderGateway)
	vehicleRepo := new(mocks.MockVehicleGateway)
	customerRepo := new(mocks.MockCustomerGateway)
	serviceRepo := new(mocks.MockServiceGateway)
	partsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	billingRepo := new(mocks.MockBillingServiceGateway)
	executionRepo := new(mocks.MockExecutionGateway)

	uc := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingRepo, executionRepo)

	serviceOrderRepo.On("GetByID", mock.Anything, "1", false).Return((*entities.ServiceOrder)(nil), errors.New("get error"))

	so, err := uc.CancelServiceOrder(context.Background(), "1")
	assert.Error(t, err)
	assert.Equal(t, "get error", err.Error())
	assert.Nil(t, so)
}

func TestCancelServiceOrder_Success(t *testing.T) {
	serviceOrderRepo := new(mocks.MockServiceOrderGateway)
	vehicleRepo := new(mocks.MockVehicleGateway)
	customerRepo := new(mocks.MockCustomerGateway)
	serviceRepo := new(mocks.MockServiceGateway)
	partsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	billingRepo := new(mocks.MockBillingServiceGateway)
	executionRepo := new(mocks.MockExecutionGateway)

	uc := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingRepo, executionRepo)

	serviceOrderRepo.On("GetByID", mock.Anything, "1", false).Return(&entities.ServiceOrder{ID: "1", Status: valueobject.StatusRecebida}, nil)
	serviceOrderRepo.On("Update", mock.Anything, mock.AnythingOfType("*entities.ServiceOrder")).Return(&entities.ServiceOrder{ID: "1", Status: valueobject.StatusCancelada}, nil)

	so, err := uc.CancelServiceOrder(context.Background(), "1")
	assert.NoError(t, err)
	assert.NotNil(t, so)
	assert.True(t, so.Status.IsCancelada())
}

func TestCancelServiceOrder_UpdateFails(t *testing.T) {
	serviceOrderRepo := new(mocks.MockServiceOrderGateway)
	vehicleRepo := new(mocks.MockVehicleGateway)
	customerRepo := new(mocks.MockCustomerGateway)
	serviceRepo := new(mocks.MockServiceGateway)
	partsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	billingRepo := new(mocks.MockBillingServiceGateway)
	executionRepo := new(mocks.MockExecutionGateway)

	uc := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingRepo, executionRepo)

	serviceOrderRepo.On("GetByID", mock.Anything, "1", false).Return(&entities.ServiceOrder{ID: "1", Status: valueobject.StatusRecebida}, nil)
	serviceOrderRepo.On("Update", mock.Anything, mock.AnythingOfType("*entities.ServiceOrder")).Return((*entities.ServiceOrder)(nil), errors.New("update error"))

	so, err := uc.CancelServiceOrder(context.Background(), "1")
	assert.Error(t, err)
	assert.Equal(t, "update error", err.Error())
	assert.Nil(t, so)
}
