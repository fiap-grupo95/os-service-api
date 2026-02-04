package usecase

import (
	"errors"
	"testing"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/mocks"

	"github.com/stretchr/testify/assert"
)

func TestVehicleService_GetAllVehicles(t *testing.T) {
	mockRepo := mocks.NewMockIVehicleGateway()
	service := NewVehicleService(mockRepo)

	expected := []entities.Vehicle{
		{ID: 1, Plate: valueobject.ParsePlate("ABC1234")},
		{ID: 2, Plate: valueobject.ParsePlate("XYZ5678")},
	}

	mockRepo.On("FindAll").Return(expected, nil)

	result, err := service.GetAllVehicles()

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestVehicleService_GetAllVehicles_Error(t *testing.T) {
	mockRepo := mocks.NewMockIVehicleGateway()
	service := NewVehicleService(mockRepo)

	mockRepo.On("FindAll").Return(nil, errors.New("db error"))

	result, err := service.GetAllVehicles()

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestVehicleService_GetVehicleByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := mocks.NewMockIVehicleGateway()
		service := NewVehicleService(mockRepo)

		vehicle := &entities.Vehicle{ID: 1, Plate: valueobject.ParsePlate("ABC1234")}

		mockRepo.On("FindByID", uint(1)).Return(vehicle, nil)

		result, err := service.GetVehicleByID(1)

		assert.NoError(t, err)
		assert.Equal(t, vehicle, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := mocks.NewMockIVehicleGateway()
		service := NewVehicleService(mockRepo)

		mockRepo.On("FindByID", uint(1)).Return((*entities.Vehicle)(nil), nil)

		result, err := service.GetVehicleByID(1)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrVehicleNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := mocks.NewMockIVehicleGateway()
		service := NewVehicleService(mockRepo)

		mockRepo.On("FindByID", uint(1)).Return((*entities.Vehicle)(nil), errors.New("db error"))

		result, err := service.GetVehicleByID(1)

		assert.Nil(t, result)
		assert.EqualError(t, err, "db error")
		mockRepo.AssertExpectations(t)
	})
}

func TestVehicleService_GetVehicleByPlate(t *testing.T) {
	t.Run("invalid plate format", func(t *testing.T) {
		mockRepo := mocks.NewMockIVehicleGateway()
		service := NewVehicleService(mockRepo)

		result, err := service.GetVehicleByPlate("INVALID")

		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrInvalidPlateFormat)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo := mocks.NewMockIVehicleGateway()
		service := NewVehicleService(mockRepo)
		vehicle := &entities.Vehicle{ID: 1, Plate: valueobject.ParsePlate("ABC1234")}

		mockRepo.On("FindByPlate", valueobject.ParsePlate("ABC1234")).Return(vehicle, nil)

		result, err := service.GetVehicleByPlate("ABC1234")

		assert.NoError(t, err)
		assert.Equal(t, vehicle, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := mocks.NewMockIVehicleGateway()
		service := NewVehicleService(mockRepo)

		mockRepo.On("FindByPlate", valueobject.ParsePlate("ABC1234")).Return((*entities.Vehicle)(nil), nil)

		result, err := service.GetVehicleByPlate("ABC1234")

		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrVehicleNotFound)
		mockRepo.AssertExpectations(t)
	})
}

func TestVehicleService_GetVehiclesByCustomerID(t *testing.T) {
	mockRepo := mocks.NewMockIVehicleGateway()
	service := NewVehicleService(mockRepo)

	vehicles := []entities.Vehicle{{ID: 1}, {ID: 2}}
	mockRepo.On("FindByCustomerID", uint(10)).Return(vehicles, nil)

	result, err := service.GetVehiclesByCustomerID(10)

	assert.NoError(t, err)
	assert.Equal(t, vehicles, result)
	mockRepo.AssertExpectations(t)
}