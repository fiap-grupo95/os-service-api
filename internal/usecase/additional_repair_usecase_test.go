package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	dto "github.com/fiap-grupo95/os-service-api/internal/infrastructure/database/model"
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

	uc := usecase.NewSOAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo)

	adr := entities.AdditionalRepair{ServiceOrderID: 2}

	mockRepoOS.On("GetByID", uint(2), false).Return(nil, errors.New("not found"))

	result, err := uc.CreateAdditionalRepair(context.Background(), adr)
	assert.Error(t, err)
	assert.Equal(t, entities.AdditionalRepair{}, result)
}

func TestCreateAdditionalRepair_ServiceNotFound(t *testing.T) {
	mockRepo := new(mocks.MockAdditionalRepairGateway)
	mockRepoOS := new(mocks.MockServiceOrderGateway)
	mockServiceRepo := new(mocks.MockServiceGateway)
	mockPartsSupplyRepo := new(mocks.MockPartsSupplyGateway)

	uc := usecase.NewSOAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo)

	adr := entities.AdditionalRepair{
		ServiceOrderID: 1,
		Services:       []entities.Service{{ID: 99}},
	}

	mockRepoOS.On("GetByID", uint(1), false).Return((&dto.ServiceOrderModel{ID: 1}).ToDomain(), nil)
	mockServiceRepo.On("GetByID", mock.Anything, uint(99)).Return(&entities.Service{}, nil)

	result, err := uc.CreateAdditionalRepair(context.Background(), adr)
	assert.Error(t, err)
	assert.Equal(t, entities.AdditionalRepair{}, result)
}

func TestCreateAdditionalRepair_PartsSupplyNotFound(t *testing.T) {
	mockRepo := new(mocks.MockAdditionalRepairGateway)
	mockRepoOS := new(mocks.MockServiceOrderGateway)
	mockServiceRepo := new(mocks.MockServiceGateway)
	mockPartsSupplyRepo := new(mocks.MockPartsSupplyGateway)

	uc := usecase.NewSOAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo)

	adr := entities.AdditionalRepair{
		ServiceOrderID: 1,
		PartsSupplies:  []entities.PartsSupply{{ID: 88}},
	}

	mockRepoOS.On("GetByID", uint(1), false).Return((&dto.ServiceOrderModel{ID: 1}).ToDomain(), nil)
	mockPartsSupplyRepo.On("GetByID", mock.Anything, uint(88)).Return(&entities.PartsSupply{}, nil)

	result, err := uc.CreateAdditionalRepair(context.Background(), adr)
	assert.Error(t, err)
	assert.Equal(t, entities.AdditionalRepair{}, result)
}
