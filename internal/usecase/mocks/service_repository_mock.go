package mocks

import (
	context "context"

	entities "github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	mock "github.com/stretchr/testify/mock"
)

type MockServiceGateway struct {
	mock.Mock
}

func (m *MockServiceGateway) GetByID(ctx context.Context, id string) (*entities.Service, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Service), args.Error(1)
}
