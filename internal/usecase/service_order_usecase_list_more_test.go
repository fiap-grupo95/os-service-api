package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
	mocks "github.com/fiap-grupo95/os-service-api/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceOrderUseCase_ListServiceOrders_FilterAndSort(t *testing.T) {
	serviceOrderRepo := new(mocks.MockServiceOrderGateway)
	vehicleRepo := new(mocks.MockVehicleGateway)
	customerRepo := new(mocks.MockCustomerGateway)
	serviceRepo := new(mocks.MockServiceGateway)
	partsSupplyRepo := new(mocks.MockPartsSupplyGateway)
	billingRepo := new(mocks.MockBillingServiceGateway)
	executionRepo := new(mocks.MockExecutionGateway)

	uc := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo, billingRepo, executionRepo)

	now := time.Now()
	later := now.Add(1 * time.Hour)

	soNil := (*entities.ServiceOrder)(nil)
	soFinalized := &entities.ServiceOrder{ID: "f", Status: valueobject.StatusFinalizada, CreatedAt: &now}
	soDelivered := &entities.ServiceOrder{ID: "d", Status: valueobject.StatusEntregue, CreatedAt: &now}
	soPendingNoTime := &entities.ServiceOrder{ID: "a", Status: valueobject.StatusRecebida, CreatedAt: nil}
	soDiagOlder := &entities.ServiceOrder{ID: "b", Status: valueobject.StatusEmDiagnostico, CreatedAt: &now}
	soDiagNewer := &entities.ServiceOrder{ID: "c", Status: valueobject.StatusEmDiagnostico, CreatedAt: &later}

	serviceOrderRepo.On("List", mock.Anything).Return([]*entities.ServiceOrder{soNil, soFinalized, soDelivered, soPendingNoTime, soDiagNewer, soDiagOlder}, nil).Once()

	list, err := uc.ListServiceOrders(context.Background())
	assert.NoError(t, err)

	// finalized and delivered filtered out, nil filtered out
	assert.Len(t, list, 3)

	// EmDiagnostico has higher priority than Recebida, so it comes first.
	assert.Equal(t, "b", list[0].ID)
	assert.Equal(t, "c", list[1].ID)

	// CreatedAt nil sorts last within same priority.
	assert.Equal(t, "a", list[2].ID)
}
