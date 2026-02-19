package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
	"github.com/fiap-grupo95/os-service-api/internal/usecase"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateAdditionalRepair_ServiceOrderNotFound(t *testing.T) {
	mockRepo := new(mocks.MockAdditionalRepairGateway)
	mockRepoOS := new(mocks.MockServiceOrderGateway)
	mockServiceRepo := new(mocks.MockServiceGateway)
	mockPartsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	mockBillingServiceRepo := new(mocks.MockBillingServiceGateway)

	uc := usecase.NewAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo, mockBillingServiceRepo)

	adr := entities.AdditionalRepair{ServiceOrderID: "2"}

	mockRepoOS.On("GetByID", mock.Anything, "2", false).Return(nil, errors.New("not found"))

	result, err := uc.CreateAdditionalRepair(context.Background(), adr)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCreateAdditionalRepair_Success(t *testing.T) {
	mockRepo := new(mocks.MockAdditionalRepairGateway)
	mockRepoOS := new(mocks.MockServiceOrderGateway)
	mockServiceRepo := new(mocks.MockServiceGateway)
	mockPartsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	mockBillingServiceRepo := new(mocks.MockBillingServiceGateway)

	uc := usecase.NewAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo, mockBillingServiceRepo)

	adr := entities.AdditionalRepair{
		ServiceOrderID: "1",
		Description:    "desc",
		Services:       []entities.Service{{ID: "s1"}},
		PartsSupplies:  []entities.PartsSupply{{ID: "p1", Quantity: 1}},
	}

	mockRepoOS.On("GetByID", mock.Anything, "1", false).Return(&entities.ServiceOrder{ID: "1"}, nil)
	mockServiceRepo.On("GetByID", mock.Anything, "s1").Return(&entities.Service{ID: "s1"}, nil)
	mockPartsSupplyRepo.On("GetByID", mock.Anything, "p1").Return(&entities.PartsSupply{ID: "p1", Quantity: 10}, nil)
	mockPartsSupplyRepo.On("AuthorizeReserve", mock.Anything, mock.Anything).Return(nil)

	created := &entities.AdditionalRepair{ID: "ar-1", ServiceOrderID: "1", Status: valueobject.StatusARAberta, PartsSupplies: []entities.PartsSupply{{ID: "p1", Quantity: 1}}}
	mockRepo.On("CreateAdditionalRepair", mock.Anything, mock.AnythingOfType("*entities.AdditionalRepair")).Return(created, nil)
	mockPartsSupplyRepo.On("Reserve", mock.Anything, mock.Anything).Return(nil)
	mockBillingServiceRepo.On("CreateEstimate", mock.Anything, (*entities.ServiceOrder)(nil), created).Return(&entities.Estimate{ID: "est-1"}, nil)
	mockRepo.On("UpdateAdditionalRepair", mock.Anything, mock.AnythingOfType("*entities.AdditionalRepair")).Return(created, nil)

	result, err := uc.CreateAdditionalRepair(context.Background(), adr)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, valueobject.StatusARAguardandoAprovacao, result.Status)
}

func TestCreateAdditionalRepair_AuthorizeReserveFails(t *testing.T) {
	mockRepo := new(mocks.MockAdditionalRepairGateway)
	mockRepoOS := new(mocks.MockServiceOrderGateway)
	mockServiceRepo := new(mocks.MockServiceGateway)
	mockPartsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	mockBillingServiceRepo := new(mocks.MockBillingServiceGateway)

	uc := usecase.NewAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo, mockBillingServiceRepo)

	adr := entities.AdditionalRepair{ServiceOrderID: "1", PartsSupplies: []entities.PartsSupply{{ID: "p1"}}}

	mockRepoOS.On("GetByID", mock.Anything, "1", false).Return(&entities.ServiceOrder{ID: "1"}, nil)
	mockPartsSupplyRepo.On("GetByID", mock.Anything, "p1").Return(&entities.PartsSupply{ID: "p1"}, nil)
	mockPartsSupplyRepo.On("AuthorizeReserve", mock.Anything, mock.Anything).Return(errors.New("auth error"))

	result, err := uc.CreateAdditionalRepair(context.Background(), adr)
	assert.Error(t, err)
	assert.Equal(t, "auth error", err.Error())
	assert.Nil(t, result)
}

func TestCreateAdditionalRepair_ReserveFails_RollbackCalled(t *testing.T) {
	mockRepo := new(mocks.MockAdditionalRepairGateway)
	mockRepoOS := new(mocks.MockServiceOrderGateway)
	mockServiceRepo := new(mocks.MockServiceGateway)
	mockPartsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	mockBillingServiceRepo := new(mocks.MockBillingServiceGateway)

	uc := usecase.NewAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo, mockBillingServiceRepo)

	adr := entities.AdditionalRepair{ServiceOrderID: "1", PartsSupplies: []entities.PartsSupply{{ID: "p1", Quantity: 1}}}

	mockRepoOS.On("GetByID", mock.Anything, "1", false).Return(&entities.ServiceOrder{ID: "1"}, nil)
	mockPartsSupplyRepo.On("GetByID", mock.Anything, "p1").Return(&entities.PartsSupply{ID: "p1", Quantity: 10}, nil)
	mockPartsSupplyRepo.On("AuthorizeReserve", mock.Anything, mock.Anything).Return(nil)

	created := &entities.AdditionalRepair{ID: "ar-1", ServiceOrderID: "1", Status: valueobject.StatusARAberta, PartsSupplies: []entities.PartsSupply{{ID: "p1", Quantity: 1}}}
	mockRepo.On("CreateAdditionalRepair", mock.Anything, mock.AnythingOfType("*entities.AdditionalRepair")).Return(created, nil)

	mockPartsSupplyRepo.On("Reserve", mock.Anything, mock.Anything).Return(errors.New("reserve error"))

	// rollback path: CancelAdditionalRepair -> UpdateAdditionalRepair
	mockRepo.On("GetByID", mock.Anything, "ar-1").Return(created, nil)
	mockRepo.On("UpdateAdditionalRepair", mock.Anything, mock.AnythingOfType("*entities.AdditionalRepair")).Return(created, nil)

	result, err := uc.CreateAdditionalRepair(context.Background(), adr)
	assert.Error(t, err)
	assert.Equal(t, "reserve error", err.Error())
	assert.Nil(t, result)
}

func TestCreateAdditionalRepair_CreateEstimateFails_RollbackCalled(t *testing.T) {
	mockRepo := new(mocks.MockAdditionalRepairGateway)
	mockRepoOS := new(mocks.MockServiceOrderGateway)
	mockServiceRepo := new(mocks.MockServiceGateway)
	mockPartsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	mockBillingServiceRepo := new(mocks.MockBillingServiceGateway)

	uc := usecase.NewAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo, mockBillingServiceRepo)

	adr := entities.AdditionalRepair{ServiceOrderID: "1", PartsSupplies: []entities.PartsSupply{{ID: "p1", Quantity: 1}}}

	mockRepoOS.On("GetByID", mock.Anything, "1", false).Return(&entities.ServiceOrder{ID: "1"}, nil)
	mockPartsSupplyRepo.On("GetByID", mock.Anything, "p1").Return(&entities.PartsSupply{ID: "p1", Quantity: 10}, nil)
	mockPartsSupplyRepo.On("AuthorizeReserve", mock.Anything, mock.Anything).Return(nil)

	created := &entities.AdditionalRepair{ID: "ar-1", ServiceOrderID: "1", Status: valueobject.StatusARAberta, PartsSupplies: []entities.PartsSupply{{ID: "p1", Quantity: 1}}}
	mockRepo.On("CreateAdditionalRepair", mock.Anything, mock.AnythingOfType("*entities.AdditionalRepair")).Return(created, nil)
	mockPartsSupplyRepo.On("Reserve", mock.Anything, mock.Anything).Return(nil)

	mockBillingServiceRepo.On("CreateEstimate", mock.Anything, (*entities.ServiceOrder)(nil), created).Return((*entities.Estimate)(nil), errors.New("billing error"))

	// rollback path: CancelAdditionalRepair
	mockRepo.On("GetByID", mock.Anything, "ar-1").Return(created, nil)
	mockRepo.On("UpdateAdditionalRepair", mock.Anything, mock.AnythingOfType("*entities.AdditionalRepair")).Return(created, nil)

	result, err := uc.CreateAdditionalRepair(context.Background(), adr)
	assert.Error(t, err)
	assert.Equal(t, "billing error", err.Error())
	assert.Nil(t, result)
}

func TestCancelAdditionalRepair_AlreadyCanceled(t *testing.T) {
	mockRepo := new(mocks.MockAdditionalRepairGateway)
	mockRepoOS := new(mocks.MockServiceOrderGateway)
	mockServiceRepo := new(mocks.MockServiceGateway)
	mockPartsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	mockBillingServiceRepo := new(mocks.MockBillingServiceGateway)

	uc := usecase.NewAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo, mockBillingServiceRepo)

	canceled := &entities.AdditionalRepair{ID: "ar-1", Status: valueobject.StatusARCancelada}
	mockRepo.On("GetByID", mock.Anything, "ar-1").Return(canceled, nil)

	result, err := uc.CancelAdditionalRepair(context.Background(), "ar-1")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetAdditionalRepair_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockAdditionalRepairGateway)
	mockRepoOS := new(mocks.MockServiceOrderGateway)
	mockServiceRepo := new(mocks.MockServiceGateway)
	mockPartsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	mockBillingServiceRepo := new(mocks.MockBillingServiceGateway)

	uc := usecase.NewAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo, mockBillingServiceRepo)

	mockRepo.On("GetByID", mock.Anything, "x").Return(&entities.AdditionalRepair{}, nil)

	result, err := uc.GetAdditionalRepair(context.Background(), "x")
	assert.Error(t, err)
	assert.Equal(t, usecase.ErrAdditionalRepairNotFound, err)
	assert.Nil(t, result)
}

func TestGetAdditionalRepairBySO_NilSliceBecomesEmpty(t *testing.T) {
	mockRepo := new(mocks.MockAdditionalRepairGateway)
	mockRepoOS := new(mocks.MockServiceOrderGateway)
	mockServiceRepo := new(mocks.MockServiceGateway)
	mockPartsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	mockBillingServiceRepo := new(mocks.MockBillingServiceGateway)

	uc := usecase.NewAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo, mockBillingServiceRepo)

	mockRepo.On("GetByServiceOrderID", mock.Anything, "1").Return(([]entities.AdditionalRepair)(nil), nil)

	list, err := uc.GetAdditionalRepairBySO(context.Background(), "1")
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Len(t, list, 0)
}

func TestCustomerApprovalStatus_InvalidStatus(t *testing.T) {
	mockRepo := new(mocks.MockAdditionalRepairGateway)
	mockRepoOS := new(mocks.MockServiceOrderGateway)
	mockServiceRepo := new(mocks.MockServiceGateway)
	mockPartsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	mockBillingServiceRepo := new(mocks.MockBillingServiceGateway)

	uc := usecase.NewAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo, mockBillingServiceRepo)

	ar := &entities.AdditionalRepair{ID: "ar-1", Status: valueobject.StatusARAberta}
	mockRepo.On("GetByID", mock.Anything, "ar-1").Return(ar, nil)

	result, err := uc.CustomerApprovalStatus(context.Background(), "ar-1", usecase.APPROVED_FLOW)
	assert.Error(t, err)
	assert.Equal(t, usecase.ErrStatusNotPermitted, err)
	assert.Nil(t, result)
}

func TestRollback_NilAdditionalRepair(t *testing.T) {
	mockRepo := new(mocks.MockAdditionalRepairGateway)
	mockRepoOS := new(mocks.MockServiceOrderGateway)
	mockServiceRepo := new(mocks.MockServiceGateway)
	mockPartsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	mockBillingServiceRepo := new(mocks.MockBillingServiceGateway)

	uc := usecase.NewAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo, mockBillingServiceRepo)

	err := uc.Rollback(context.Background(), nil, false)
	assert.Error(t, err)
	assert.Equal(t, usecase.ErrAdditionalRepairNotFound, err)
}

func TestCreateAdditionalRepair_ServiceNotFound(t *testing.T) {
	mockRepo := new(mocks.MockAdditionalRepairGateway)
	mockRepoOS := new(mocks.MockServiceOrderGateway)
	mockServiceRepo := new(mocks.MockServiceGateway)
	mockPartsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	mockBillingServiceRepo := new(mocks.MockBillingServiceGateway)

	uc := usecase.NewAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo, mockBillingServiceRepo)

	adr := entities.AdditionalRepair{
		ServiceOrderID: "1",
		Services:       []entities.Service{{ID: "99"}},
	}

	mockRepoOS.On("GetByID", mock.Anything, "1", false).Return(&entities.ServiceOrder{ID: "1"}, nil)
	mockServiceRepo.On("GetByID", mock.Anything, "99").Return(&entities.Service{}, nil)

	result, err := uc.CreateAdditionalRepair(context.Background(), adr)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCreateAdditionalRepair_PartsSupplyNotFound(t *testing.T) {
	mockRepo := new(mocks.MockAdditionalRepairGateway)
	mockRepoOS := new(mocks.MockServiceOrderGateway)
	mockServiceRepo := new(mocks.MockServiceGateway)
	mockPartsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	mockBillingServiceRepo := new(mocks.MockBillingServiceGateway)

	uc := usecase.NewAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo, mockBillingServiceRepo)

	adr := entities.AdditionalRepair{
		ServiceOrderID: "1",
		PartsSupplies:  []entities.PartsSupply{{ID: "88"}},
	}

	mockRepoOS.On("GetByID", mock.Anything, "1", false).Return(&entities.ServiceOrder{ID: "1"}, nil)
	mockPartsSupplyRepo.On("GetByID", mock.Anything, "88").Return(&entities.PartsSupply{}, nil)

	result, err := uc.CreateAdditionalRepair(context.Background(), adr)
	assert.Error(t, err)
	assert.Nil(t, result)
}
