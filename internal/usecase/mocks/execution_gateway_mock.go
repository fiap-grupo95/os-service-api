package mocks

import (
	"context"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/stretchr/testify/mock"
)

type MockExecutionGateway struct {
	mock.Mock
}

func (m *MockExecutionGateway) CreateExecution(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.Execution, error) {
	args := m.Called(ctx, serviceOrder)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Execution), args.Error(1)
}

func (m *MockExecutionGateway) FinishExecution(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.Execution, error) {
	args := m.Called(ctx, serviceOrder)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Execution), args.Error(1)
}
