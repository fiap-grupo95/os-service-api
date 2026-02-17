// Code generated manually to match updated repository interface.
// Package mocks provides gomock implementations for use case dependencies.
package mocks

import (
	"context"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	mock "github.com/stretchr/testify/mock"
)

// MockAdditionalRepairGateway is a testify-based mock that implements the
// IAdditionalRepairRepository interface used by the AdditionalRepairUseCase.
type MockAdditionalRepairGateway struct {
	mock.Mock
}

func (m *MockAdditionalRepairGateway) CreateAdditionalRepair(ctx context.Context, additionalRepair *entities.AdditionalRepair) (*entities.AdditionalRepair, error) {
	args := m.Called(ctx, additionalRepair)
	var result *entities.AdditionalRepair
	if v := args.Get(0); v != nil {
		if cast, ok := v.(*entities.AdditionalRepair); ok {
			result = cast
		}
	}
	return result, args.Error(1)
}

func (m *MockAdditionalRepairGateway) GetByID(ctx context.Context, id string) (*entities.AdditionalRepair, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.AdditionalRepair), args.Error(1)
}

func (m *MockAdditionalRepairGateway) GetByServiceOrderID(ctx context.Context, serviceOrderID string) ([]entities.AdditionalRepair, error) {
	args := m.Called(ctx, serviceOrderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.AdditionalRepair), args.Error(1)
}

func (m *MockAdditionalRepairGateway) UpdateAdditionalRepair(ctx context.Context, additionalRepair *entities.AdditionalRepair) (*entities.AdditionalRepair, error) {
	args := m.Called(ctx, additionalRepair)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.AdditionalRepair), args.Error(1)
}
