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

func TestCustomerApprovalStatus_ApprovedFlow_Success(t *testing.T) {
	mockRepo := new(mocks.MockAdditionalRepairGateway)
	mockRepoOS := new(mocks.MockServiceOrderGateway)
	mockServiceRepo := new(mocks.MockServiceGateway)
	mockPartsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	mockBillingServiceRepo := new(mocks.MockBillingServiceGateway)

	uc := usecase.NewAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo, mockBillingServiceRepo)

	ar := &entities.AdditionalRepair{ID: "ar-1", Status: valueobject.StatusARAguardandoAprovacao, PartsSupplies: []entities.PartsSupply{{ID: "p1", Quantity: 1}}}
	mockRepo.On("GetByID", mock.Anything, "ar-1").Return(ar, nil)
	mockPartsSupplyRepo.On("WriteOff", mock.Anything, mock.Anything).Return(nil)
	mockBillingServiceRepo.On("ApproveEstimate", mock.Anything, (*entities.ServiceOrder)(nil), ar).Return(&entities.Estimate{ID: "est-1"}, nil)
	mockRepo.On("UpdateAdditionalRepair", mock.Anything, mock.AnythingOfType("*entities.AdditionalRepair")).Return(ar, nil)

	result, err := uc.CustomerApprovalStatus(context.Background(), "ar-1", usecase.APPROVED_FLOW)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Status.IsAprovada())
}

func TestCustomerApprovalStatus_RejectedFlow_Success(t *testing.T) {
	mockRepo := new(mocks.MockAdditionalRepairGateway)
	mockRepoOS := new(mocks.MockServiceOrderGateway)
	mockServiceRepo := new(mocks.MockServiceGateway)
	mockPartsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	mockBillingServiceRepo := new(mocks.MockBillingServiceGateway)

	uc := usecase.NewAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo, mockBillingServiceRepo)

	ar := &entities.AdditionalRepair{ID: "ar-1", Status: valueobject.StatusARAguardandoAprovacao, PartsSupplies: []entities.PartsSupply{{ID: "p1", Quantity: 1}}}
	mockRepo.On("GetByID", mock.Anything, "ar-1").Return(ar, nil)
	mockPartsSupplyRepo.On("Release", mock.Anything, mock.Anything).Return(nil)
	mockBillingServiceRepo.On("RejectEstimate", mock.Anything, (*entities.ServiceOrder)(nil), ar).Return(&entities.Estimate{ID: "est-1"}, nil)
	mockRepo.On("UpdateAdditionalRepair", mock.Anything, mock.AnythingOfType("*entities.AdditionalRepair")).Return(ar, nil)

	result, err := uc.CustomerApprovalStatus(context.Background(), "ar-1", usecase.REJECTED_FLOW)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Status.IsRejeitada())
}

func TestCustomerApprovalStatus_ApprovedFlow_WriteOffFails(t *testing.T) {
	mockRepo := new(mocks.MockAdditionalRepairGateway)
	mockRepoOS := new(mocks.MockServiceOrderGateway)
	mockServiceRepo := new(mocks.MockServiceGateway)
	mockPartsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	mockBillingServiceRepo := new(mocks.MockBillingServiceGateway)

	uc := usecase.NewAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo, mockBillingServiceRepo)

	ar := &entities.AdditionalRepair{ID: "ar-1", Status: valueobject.StatusARAguardandoAprovacao, PartsSupplies: []entities.PartsSupply{{ID: "p1", Quantity: 1}}}
	mockRepo.On("GetByID", mock.Anything, "ar-1").Return(ar, nil)
	mockPartsSupplyRepo.On("WriteOff", mock.Anything, mock.Anything).Return(errors.New("writeoff error"))

	result, err := uc.CustomerApprovalStatus(context.Background(), "ar-1", usecase.APPROVED_FLOW)
	assert.Error(t, err)
	assert.Equal(t, "writeoff error", err.Error())
	assert.Nil(t, result)
}

func TestCustomerApprovalStatus_RejectedFlow_ReleaseFails(t *testing.T) {
	mockRepo := new(mocks.MockAdditionalRepairGateway)
	mockRepoOS := new(mocks.MockServiceOrderGateway)
	mockServiceRepo := new(mocks.MockServiceGateway)
	mockPartsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	mockBillingServiceRepo := new(mocks.MockBillingServiceGateway)

	uc := usecase.NewAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo, mockBillingServiceRepo)

	ar := &entities.AdditionalRepair{ID: "ar-1", Status: valueobject.StatusARAguardandoAprovacao, PartsSupplies: []entities.PartsSupply{{ID: "p1", Quantity: 1}}}
	mockRepo.On("GetByID", mock.Anything, "ar-1").Return(ar, nil)
	mockPartsSupplyRepo.On("Release", mock.Anything, mock.Anything).Return(errors.New("release error"))

	result, err := uc.CustomerApprovalStatus(context.Background(), "ar-1", usecase.REJECTED_FLOW)
	assert.Error(t, err)
	assert.Equal(t, "release error", err.Error())
	assert.Nil(t, result)
}

func TestCustomerApprovalStatus_ApproveEstimateFails(t *testing.T) {
	mockRepo := new(mocks.MockAdditionalRepairGateway)
	mockRepoOS := new(mocks.MockServiceOrderGateway)
	mockServiceRepo := new(mocks.MockServiceGateway)
	mockPartsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	mockBillingServiceRepo := new(mocks.MockBillingServiceGateway)

	uc := usecase.NewAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo, mockBillingServiceRepo)

	ar := &entities.AdditionalRepair{ID: "ar-1", Status: valueobject.StatusARAguardandoAprovacao, PartsSupplies: []entities.PartsSupply{{ID: "p1", Quantity: 1}}}
	mockRepo.On("GetByID", mock.Anything, "ar-1").Return(ar, nil)
	mockPartsSupplyRepo.On("WriteOff", mock.Anything, mock.Anything).Return(nil)
	mockBillingServiceRepo.On("ApproveEstimate", mock.Anything, (*entities.ServiceOrder)(nil), ar).Return((*entities.Estimate)(nil), errors.New("billing error"))

	result, err := uc.CustomerApprovalStatus(context.Background(), "ar-1", usecase.APPROVED_FLOW)
	assert.Error(t, err)
	assert.Equal(t, "billing error", err.Error())
	assert.Nil(t, result)
}

func TestCustomerApprovalStatus_UpdateFails(t *testing.T) {
	mockRepo := new(mocks.MockAdditionalRepairGateway)
	mockRepoOS := new(mocks.MockServiceOrderGateway)
	mockServiceRepo := new(mocks.MockServiceGateway)
	mockPartsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	mockBillingServiceRepo := new(mocks.MockBillingServiceGateway)

	uc := usecase.NewAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo, mockBillingServiceRepo)

	ar := &entities.AdditionalRepair{ID: "ar-1", Status: valueobject.StatusARAguardandoAprovacao, PartsSupplies: []entities.PartsSupply{{ID: "p1", Quantity: 1}}}
	mockRepo.On("GetByID", mock.Anything, "ar-1").Return(ar, nil)
	mockPartsSupplyRepo.On("WriteOff", mock.Anything, mock.Anything).Return(nil)
	mockBillingServiceRepo.On("ApproveEstimate", mock.Anything, (*entities.ServiceOrder)(nil), ar).Return(&entities.Estimate{ID: "est-1"}, nil)
	mockRepo.On("UpdateAdditionalRepair", mock.Anything, mock.AnythingOfType("*entities.AdditionalRepair")).Return((*entities.AdditionalRepair)(nil), errors.New("update error"))

	result, err := uc.CustomerApprovalStatus(context.Background(), "ar-1", usecase.APPROVED_FLOW)
	assert.Error(t, err)
	assert.Equal(t, "update error", err.Error())
	assert.Nil(t, result)
}

func TestRollback_WithReservedPartsSupplyAndEstimate(t *testing.T) {
	mockRepo := new(mocks.MockAdditionalRepairGateway)
	mockRepoOS := new(mocks.MockServiceOrderGateway)
	mockServiceRepo := new(mocks.MockServiceGateway)
	mockPartsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	mockBillingServiceRepo := new(mocks.MockBillingServiceGateway)

	uc := usecase.NewAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo, mockBillingServiceRepo)

	ar := &entities.AdditionalRepair{ID: "ar-1", Status: valueobject.StatusARAberta, PartsSupplies: []entities.PartsSupply{{ID: "p1", Quantity: 1}}, Estimate: &entities.Estimate{ID: "est-1"}}

	mockPartsSupplyRepo.On("Release", mock.Anything, mock.Anything).Return(nil)
	mockBillingServiceRepo.On("CancelEstimate", mock.Anything, (*entities.ServiceOrder)(nil), ar).Return(&entities.Estimate{ID: "est-1"}, nil)
	mockRepo.On("UpdateAdditionalRepair", mock.Anything, mock.AnythingOfType("*entities.AdditionalRepair")).Return(ar, nil)

	// CancelAdditionalRepair path
	mockRepo.On("GetByID", mock.Anything, "ar-1").Return(ar, nil)
	mockRepo.On("UpdateAdditionalRepair", mock.Anything, mock.AnythingOfType("*entities.AdditionalRepair")).Return(ar, nil)

	err := uc.Rollback(context.Background(), ar, true)
	assert.NoError(t, err)
}
