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

func (m *MockAdditionalRepairGateway) Create(ctx context.Context, additionalRepair entities.AdditionalRepair) (entities.AdditionalRepair, error) {
	args := m.Called(ctx, additionalRepair)
	var result entities.AdditionalRepair
	if v := args.Get(0); v != nil {
		if cast, ok := v.(entities.AdditionalRepair); ok {
			result = cast
		}
	}
	return result, args.Error(1)
}

func (m *MockAdditionalRepairGateway) GetByID(ctx context.Context, id uint) (entities.AdditionalRepair, error) {
	args := m.Called(ctx, id)
	var result entities.AdditionalRepair
	if v := args.Get(0); v != nil {
		if cast, ok := v.(entities.AdditionalRepair); ok {
			result = cast
		}
	}
	return result, args.Error(1)
}

func (m *MockAdditionalRepairGateway) AddPartSupplyAndService(ctx context.Context, additionalRepairID uint, services []entities.Service, partsSupplies []entities.PartsSupply, newEstimate float64) error {
	args := m.Called(ctx, additionalRepairID, services, partsSupplies, newEstimate)
	return args.Error(0)
}

func (m *MockAdditionalRepairGateway) ReplacePartSupplyAndService(ctx context.Context, additionalRepairID uint, services []entities.Service, partsSupplies []entities.PartsSupply, newEstimate float64) error {
	args := m.Called(ctx, additionalRepairID, services, partsSupplies, newEstimate)
	return args.Error(0)
}

func (m *MockAdditionalRepairGateway) GetByServiceOrder(ctx context.Context, serviceOrderId uint) ([]entities.AdditionalRepair, error) {
	args := m.Called(ctx, serviceOrderId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.AdditionalRepair), args.Error(1)
}

func (m *MockAdditionalRepairGateway) UpdateStatus(ctx context.Context, id uint, status entities.AdditionalRepairStatusDTO) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockAdditionalRepairGateway) GetPartsSupplyQuantity(ctx context.Context, partsSupplyID uint, additionalRepairID uint) (int, error) {
	args := m.Called(ctx, partsSupplyID, additionalRepairID)
	return args.Int(0), args.Error(1)
}
