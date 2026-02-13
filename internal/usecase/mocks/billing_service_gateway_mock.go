package mocks

import (
	"context"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/stretchr/testify/mock"
)

type MockBillingServiceGateway struct {
	mock.Mock
}

func (m *MockBillingServiceGateway) CreateEstimate(ctx context.Context, serviceOrder *entities.ServiceOrder) (*float64, error) {
	args := m.Called(ctx, serviceOrder)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	value, ok := args.Get(0).(float64)
	if !ok {
		return nil, args.Error(1)
	}
	return &value, args.Error(1)
}
