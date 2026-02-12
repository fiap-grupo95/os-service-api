package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	dto "github.com/fiap-grupo95/os-service-api/internal/infrastructure/database/model"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/mocks"
)

func TestPaymentUseCase_CreatePayment(t *testing.T) {
	mockDate := time.Date(2026, time.February, 12, 14, 22, 27, 0, time.Local)

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		mockPaymentRepo := new(mocks.MockPaymentGateway)
		mockServiceOrderRepo := new(mocks.MockServiceOrderGateway)
		u := NewPaymentUseCase(mockPaymentRepo, mockServiceOrderRepo)

		mockServiceOrderDTO := &dto.ServiceOrderModel{ID: 1, Estimate: 100.0}
		payment := &entities.Payment{ID: 1, ServiceOrderID: 1, Amount: 100.0, PaymentDate: mockDate}

		mockServiceOrderRepo.On("GetByID", uint(1), false).Return(mockServiceOrderDTO, nil)
		mockPaymentRepo.On("GetByServiceOrderID", ctx, uint(1)).Return(entities.Payment{}, nil)
		mockPaymentRepo.On("Create", ctx, mock.AnythingOfType("entities.Payment")).Return(entities.Payment{
			ID:             1,
			ServiceOrderID: 1,
			Amount:         100,
			PaymentDate:    mockDate,
		}, nil)
		result, err := u.CreatePayment(ctx, payment)
		assert.NoError(t, err)
		assert.Equal(t, payment, result)
	})

	t.Run("service order not found", func(t *testing.T) {
		ctx := context.Background()
		mockPaymentRepo := new(mocks.MockPaymentGateway)
		mockServiceOrderRepo := new(mocks.MockServiceOrderGateway)
		u := NewPaymentUseCase(mockPaymentRepo, mockServiceOrderRepo)

		badPayment := &entities.Payment{ID: 2, ServiceOrderID: 2, Amount: 100.0}
		mockServiceOrderRepo.On("GetByID", uint(2), false).Return(nil, errors.New("not found"))
		_, err := u.CreatePayment(ctx, badPayment)
		assert.Error(t, err)
		assert.EqualError(t, err, "not found")
	})

	t.Run("amount does not match", func(t *testing.T) {
		ctx := context.Background()
		mockPaymentRepo := new(mocks.MockPaymentGateway)
		mockServiceOrderRepo := new(mocks.MockServiceOrderGateway)
		u := NewPaymentUseCase(mockPaymentRepo, mockServiceOrderRepo)

		mockServiceOrderDTO := &dto.ServiceOrderModel{ID: 1, Estimate: 100.0}
		badPayment := &entities.Payment{ID: 3, ServiceOrderID: 1, Amount: 200.0}
		mockServiceOrderRepo.On("GetByID", uint(1), false).Return(mockServiceOrderDTO, nil)
		_, err := u.CreatePayment(ctx, badPayment)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrPaymentAmountDoesNotMatch)
	})

	t.Run("payment already exists", func(t *testing.T) {
		ctx := context.Background()
		mockPaymentRepo := new(mocks.MockPaymentGateway)
		mockServiceOrderRepo := new(mocks.MockServiceOrderGateway)
		u := NewPaymentUseCase(mockPaymentRepo, mockServiceOrderRepo)

		mockServiceOrderDTO := &dto.ServiceOrderModel{ID: 1, Estimate: 100.0}
		mockServiceOrderRepo.On("GetByID", uint(1), false).Return(mockServiceOrderDTO, nil)
		mockPaymentRepo.On("GetByServiceOrderID", ctx, uint(1)).Return(entities.Payment{ID: 99}, nil)
		badPayment := &entities.Payment{ID: 4, ServiceOrderID: 1, Amount: 100.0}
		_, err := u.CreatePayment(ctx, badPayment)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrPaymentAlreadyExists)
	})

	t.Run("repo create error", func(t *testing.T) {
		ctx := context.Background()
		mockPaymentRepo := new(mocks.MockPaymentGateway)
		mockServiceOrderRepo := new(mocks.MockServiceOrderGateway)
		u := NewPaymentUseCase(mockPaymentRepo, mockServiceOrderRepo)

		mockServiceOrderDTO := &dto.ServiceOrderModel{ID: 1, Estimate: 100.0}
		payment := &entities.Payment{ID: 1, ServiceOrderID: 1, Amount: 100.0}
		mockServiceOrderRepo.On("GetByID", uint(1), false).Return(mockServiceOrderDTO, nil)
		mockPaymentRepo.On("GetByServiceOrderID", ctx, uint(1)).Return(entities.Payment{}, nil)
		mockPaymentRepo.On("Create", ctx, mock.AnythingOfType("entities.Payment")).Return(entities.Payment{PaymentDate: mockDate}, errors.New("db error"))
		_, err := u.CreatePayment(ctx, payment)
		assert.Error(t, err)
		assert.EqualError(t, err, "db error")
	})
}

func TestPaymentUseCase_GetPaymentByID(t *testing.T) {
	ctx := context.Background()
	mockPaymentRepo := new(mocks.MockPaymentGateway)
	mockServiceOrderRepo := new(mocks.MockServiceOrderGateway)
	u := NewPaymentUseCase(mockPaymentRepo, mockServiceOrderRepo)

	t.Run("success", func(t *testing.T) {
		mockPaymentRepo.On("GetByID", ctx, uint(1)).Return(entities.Payment{ID: 1}, nil)
		result, err := u.GetPaymentByID(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), result.ID)
	})

	t.Run("not found", func(t *testing.T) {
		mockPaymentRepo.On("GetByID", ctx, uint(2)).Return(entities.Payment{}, nil)
		_, err := u.GetPaymentByID(ctx, 2)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrorPaymentNotFound)
	})

	t.Run("repo error", func(t *testing.T) {
		mockPaymentRepo.On("GetByID", ctx, uint(3)).Return(entities.Payment{}, errors.New("db error"))
		_, err := u.GetPaymentByID(ctx, 3)
		assert.Error(t, err)
		assert.EqualError(t, err, "db error")
	})
}
