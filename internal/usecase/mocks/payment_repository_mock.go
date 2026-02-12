package mocks

import (
	"context"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/stretchr/testify/mock"
)

// MockPaymentGateway provides a gomock mock for the payment repository.
type MockPaymentGateway struct {
	mock.Mock
}

// Create mocks Create.
func (m *MockPaymentGateway) Create(ctx context.Context, payment entities.Payment) (entities.Payment, error) {
	args := m.Called(ctx, payment)
	if args.Get(0) == nil {
		return entities.Payment{}, args.Error(1)
	}
	return args.Get(0).(entities.Payment), args.Error(1)
}

// GetByID mocks GetByID.
func (m *MockPaymentGateway) GetByID(ctx context.Context, id uint) (entities.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return entities.Payment{}, args.Error(1)
	}
	return args.Get(0).(entities.Payment), args.Error(1)
}

// GetByServiceOrderID mocks GetByServiceOrderID.
func (m *MockPaymentGateway) GetByServiceOrderID(ctx context.Context, serviceOrderID uint) (entities.Payment, error) {
	args := m.Called(ctx, serviceOrderID)
	if args.Get(0) == nil {
		return entities.Payment{}, args.Error(1)
	}
	return args.Get(0).(entities.Payment), args.Error(1)
}