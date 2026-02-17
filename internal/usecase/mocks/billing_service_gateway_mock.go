package mocks

import (
	"context"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/stretchr/testify/mock"
)

type MockBillingServiceGateway struct {
	mock.Mock
}

func (m *MockBillingServiceGateway) CreateEstimate(ctx context.Context, serviceOrder *entities.ServiceOrder, additionalRepair *entities.AdditionalRepair) (*entities.Estimate, error) {
	args := m.Called(ctx, serviceOrder, additionalRepair)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Estimate), args.Error(1)
}

func (m *MockBillingServiceGateway) ApproveEstimate(ctx context.Context, serviceOrder *entities.ServiceOrder, additionalRepair *entities.AdditionalRepair) (*entities.Estimate, error) {
	args := m.Called(ctx, serviceOrder, additionalRepair)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Estimate), args.Error(1)
}

func (m *MockBillingServiceGateway) RejectEstimate(ctx context.Context, serviceOrder *entities.ServiceOrder, additionalRepair *entities.AdditionalRepair) (*entities.Estimate, error) {
	args := m.Called(ctx, serviceOrder, additionalRepair)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Estimate), args.Error(1)
}

func (m *MockBillingServiceGateway) CancelEstimate(ctx context.Context, serviceOrder *entities.ServiceOrder, additionalRepair *entities.AdditionalRepair) (*entities.Estimate, error) {
	args := m.Called(ctx, serviceOrder, additionalRepair)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Estimate), args.Error(1)
}

func (m *MockBillingServiceGateway) GetPaymentByEstimateID(ctx context.Context, estimateID string) (*entities.Payment, error) {
	args := m.Called(ctx, estimateID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Payment), args.Error(1)
}

func (m *MockBillingServiceGateway) CreatePayment(ctx context.Context, estimateID string) (*entities.Payment, error) {
	args := m.Called(ctx, estimateID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Payment), args.Error(1)
}
