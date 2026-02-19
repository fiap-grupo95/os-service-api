package operations

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/constants"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEstimateStrategyFactory_GetStrategy(t *testing.T) {
	partsSupplyRepo := &mocks.MockPartsSupplyGateway{}
	billingRepo := &mocks.MockBillingServiceGateway{}

	factory := NewEstimateStrategyFactory(partsSupplyRepo, billingRepo)

	t.Run("invalid operation", func(t *testing.T) {
		strategy, err := factory.GetStrategy("does-not-exist")
		assert.Error(t, err)
		assert.Nil(t, strategy)
	})

	t.Run("approve operation", func(t *testing.T) {
		strategy, err := factory.GetStrategy(constants.ESTIMATE_APPROVE)
		assert.NoError(t, err)
		assert.NotNil(t, strategy)
		assert.Equal(t, valueobject.StatusAprovada, strategy.GetTargetStatus())
	})

	t.Run("reject operation", func(t *testing.T) {
		strategy, err := factory.GetStrategy(constants.ESTIMATE_REJECT)
		assert.NoError(t, err)
		assert.NotNil(t, strategy)
		assert.Equal(t, valueobject.StatusRejeitada, strategy.GetTargetStatus())
	})

	t.Run("cancel operation", func(t *testing.T) {
		strategy, err := factory.GetStrategy(constants.ESTIMATE_CANCEL)
		assert.NoError(t, err)
		assert.NotNil(t, strategy)
		assert.Equal(t, valueobject.StatusEmDiagnostico, strategy.GetTargetStatus())
	})
}

func TestApproveEstimateStrategy_Execute(t *testing.T) {
	ctx := context.Background()
	partsSupplyRepo := &mocks.MockPartsSupplyGateway{}
	billingRepo := &mocks.MockBillingServiceGateway{}

	strategy := &ApproveEstimateStrategy{partsSupplyRepo: partsSupplyRepo, billingServiceRepo: billingRepo}

	so := &entities.ServiceOrder{ID: "so-1", PartsSupplies: []entities.PartsSupply{{ID: "p1", Quantity: 1}}}

	t.Run("success", func(t *testing.T) {
		partsSupplyRepo.On("WriteOff", mock.Anything, mock.Anything).Return(nil).Once()
		billingRepo.On("ApproveEstimate", mock.Anything, so, (*entities.AdditionalRepair)(nil)).Return(&entities.Estimate{ID: "est-1"}, nil).Once()

		result, err := strategy.Execute(ctx, so)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, valueobject.StatusAprovada, result.Status)
		assert.NotNil(t, result.Estimate)
		assert.Equal(t, "est-1", result.Estimate.ID)

		partsSupplyRepo.AssertExpectations(t)
		billingRepo.AssertExpectations(t)
	})

	t.Run("writeoff fails", func(t *testing.T) {
		partsSupplyRepo2 := &mocks.MockPartsSupplyGateway{}
		billingRepo2 := &mocks.MockBillingServiceGateway{}
		strategy2 := &ApproveEstimateStrategy{partsSupplyRepo: partsSupplyRepo2, billingServiceRepo: billingRepo2}
		so2 := &entities.ServiceOrder{ID: "so-2", PartsSupplies: []entities.PartsSupply{{ID: "p1", Quantity: 1}}}

		partsSupplyRepo2.On("WriteOff", mock.Anything, mock.Anything).Return(errors.New("writeoff error")).Once()

		result, err := strategy2.Execute(ctx, so2)
		assert.Error(t, err)
		assert.Nil(t, result)

		partsSupplyRepo2.AssertExpectations(t)
		billingRepo2.AssertExpectations(t)
	})

	t.Run("approve estimate returns nil", func(t *testing.T) {
		partsSupplyRepo2 := &mocks.MockPartsSupplyGateway{}
		billingRepo2 := &mocks.MockBillingServiceGateway{}
		strategy2 := &ApproveEstimateStrategy{partsSupplyRepo: partsSupplyRepo2, billingServiceRepo: billingRepo2}
		so2 := &entities.ServiceOrder{ID: "so-2", PartsSupplies: []entities.PartsSupply{{ID: "p1", Quantity: 1}}}

		partsSupplyRepo2.On("WriteOff", mock.Anything, mock.Anything).Return(nil).Once()
		billingRepo2.On("ApproveEstimate", mock.Anything, so2, (*entities.AdditionalRepair)(nil)).Return((*entities.Estimate)(nil), nil).Once()

		result, err := strategy2.Execute(ctx, so2)
		assert.Error(t, err)
		assert.Nil(t, result)

		partsSupplyRepo2.AssertExpectations(t)
		billingRepo2.AssertExpectations(t)
	})
}

func TestRejectEstimateStrategy_Execute(t *testing.T) {
	ctx := context.Background()
	partsSupplyRepo := &mocks.MockPartsSupplyGateway{}
	billingRepo := &mocks.MockBillingServiceGateway{}

	strategy := &RejectEstimateStrategy{partsSupplyRepo: partsSupplyRepo, billingServiceRepo: billingRepo}

	t.Run("success", func(t *testing.T) {
		so := &entities.ServiceOrder{ID: "so-1", PartsSupplies: []entities.PartsSupply{{ID: "p1", Quantity: 1}}}
		partsSupplyRepo.On("Release", mock.Anything, mock.Anything).Return(nil).Once()
		billingRepo.On("RejectEstimate", mock.Anything, so, (*entities.AdditionalRepair)(nil)).Return(&entities.Estimate{ID: "est-1"}, nil).Once()

		result, err := strategy.Execute(ctx, so)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, valueobject.StatusRejeitada, result.Status)
		assert.NotNil(t, result.Estimate)

		partsSupplyRepo.AssertExpectations(t)
		billingRepo.AssertExpectations(t)
	})

	t.Run("release fails", func(t *testing.T) {
		partsSupplyRepo2 := &mocks.MockPartsSupplyGateway{}
		billingRepo2 := &mocks.MockBillingServiceGateway{}
		strategy2 := &RejectEstimateStrategy{partsSupplyRepo: partsSupplyRepo2, billingServiceRepo: billingRepo2}
		so2 := &entities.ServiceOrder{ID: "so-2", PartsSupplies: []entities.PartsSupply{{ID: "p1", Quantity: 1}}}

		partsSupplyRepo2.On("Release", mock.Anything, mock.Anything).Return(errors.New("release error")).Once()

		result, err := strategy2.Execute(ctx, so2)
		assert.Error(t, err)
		assert.Nil(t, result)

		partsSupplyRepo2.AssertExpectations(t)
		billingRepo2.AssertExpectations(t)
	})

	t.Run("reject estimate returns nil", func(t *testing.T) {
		partsSupplyRepo2 := &mocks.MockPartsSupplyGateway{}
		billingRepo2 := &mocks.MockBillingServiceGateway{}
		strategy2 := &RejectEstimateStrategy{partsSupplyRepo: partsSupplyRepo2, billingServiceRepo: billingRepo2}
		so2 := &entities.ServiceOrder{ID: "so-2", PartsSupplies: []entities.PartsSupply{{ID: "p1", Quantity: 1}}}

		partsSupplyRepo2.On("Release", mock.Anything, mock.Anything).Return(nil).Once()
		billingRepo2.On("RejectEstimate", mock.Anything, so2, (*entities.AdditionalRepair)(nil)).Return((*entities.Estimate)(nil), nil).Once()

		result, err := strategy2.Execute(ctx, so2)
		assert.Error(t, err)
		assert.Nil(t, result)

		partsSupplyRepo2.AssertExpectations(t)
		billingRepo2.AssertExpectations(t)
	})
}

func TestCancelEstimateStrategy_Execute(t *testing.T) {
	ctx := context.Background()
	partsSupplyRepo := &mocks.MockPartsSupplyGateway{}
	billingRepo := &mocks.MockBillingServiceGateway{}

	strategy := &CancelEstimateStrategy{partsSupplyRepo: partsSupplyRepo, billingServiceRepo: billingRepo}

	t.Run("success", func(t *testing.T) {
		so := &entities.ServiceOrder{ID: "so-1", PartsSupplies: []entities.PartsSupply{{ID: "p1", Quantity: 1}}}
		partsSupplyRepo.On("Release", mock.Anything, mock.Anything).Return(nil).Once()
		billingRepo.On("CancelEstimate", mock.Anything, so, (*entities.AdditionalRepair)(nil)).Return(&entities.Estimate{ID: "est-1"}, nil).Once()

		result, err := strategy.Execute(ctx, so)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, valueobject.StatusEmDiagnostico, result.Status)
		assert.NotNil(t, result.Estimate)

		partsSupplyRepo.AssertExpectations(t)
		billingRepo.AssertExpectations(t)
	})

	t.Run("release fails", func(t *testing.T) {
		partsSupplyRepo2 := &mocks.MockPartsSupplyGateway{}
		billingRepo2 := &mocks.MockBillingServiceGateway{}
		strategy2 := &CancelEstimateStrategy{partsSupplyRepo: partsSupplyRepo2, billingServiceRepo: billingRepo2}
		so2 := &entities.ServiceOrder{ID: "so-2", PartsSupplies: []entities.PartsSupply{{ID: "p1", Quantity: 1}}}

		partsSupplyRepo2.On("Release", mock.Anything, mock.Anything).Return(errors.New("release error")).Once()

		result, err := strategy2.Execute(ctx, so2)
		assert.Error(t, err)
		assert.Nil(t, result)

		partsSupplyRepo2.AssertExpectations(t)
		billingRepo2.AssertExpectations(t)
	})

	t.Run("cancel estimate returns nil", func(t *testing.T) {
		partsSupplyRepo2 := &mocks.MockPartsSupplyGateway{}
		billingRepo2 := &mocks.MockBillingServiceGateway{}
		strategy2 := &CancelEstimateStrategy{partsSupplyRepo: partsSupplyRepo2, billingServiceRepo: billingRepo2}
		so2 := &entities.ServiceOrder{ID: "so-2", PartsSupplies: []entities.PartsSupply{{ID: "p1", Quantity: 1}}}

		partsSupplyRepo2.On("Release", mock.Anything, mock.Anything).Return(nil).Once()
		billingRepo2.On("CancelEstimate", mock.Anything, so2, (*entities.AdditionalRepair)(nil)).Return((*entities.Estimate)(nil), nil).Once()

		result, err := strategy2.Execute(ctx, so2)
		assert.Error(t, err)
		assert.Nil(t, result)

		partsSupplyRepo2.AssertExpectations(t)
		billingRepo2.AssertExpectations(t)
	})
}
