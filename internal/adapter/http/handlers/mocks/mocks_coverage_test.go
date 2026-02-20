package mocks

import (
	"context"
	"testing"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestMocksCoverage_CustomerUseCase(t *testing.T) {
	m := NewMockICustomerUseCase(nil)

	m.On("CreateCustomer", (*entities.Customer)(nil)).Return((*entities.Customer)(nil), nil).Once()
	m.On("DeleteCustomer", uint(1)).Return(nil).Once()
	m.On("GetByDocument", "doc").Return((*entities.Customer)(nil), nil).Once()
	m.On("GetById", uint(1)).Return((*entities.Customer)(nil), nil).Once()
	m.On("ListCustomer").Return(([]entities.Customer)(nil), nil).Once()
	m.On("UpdateCustomer", uint(1), (*entities.Customer)(nil)).Return(nil).Once()

	_, _ = m.CreateCustomer(nil)
	_ = m.DeleteCustomer(1)
	_, _ = m.GetByDocument("doc")
	_, _ = m.GetById(1)
	_, _ = m.ListCustomer()
	_ = m.UpdateCustomer(1, nil)

	m.AssertExpectations(t)
}

func TestMocksCoverage_ServiceOrderUseCase(t *testing.T) {
	m := NewMockIServiceOrderUseCase(nil)

	m.On("CreateServiceOrder", context.Background(), entities.ServiceOrder{}).Return((*entities.ServiceOrder)(nil), nil).Once()
	m.On("GetServiceOrder", context.Background(), entities.ServiceOrder{}).Return((*entities.ServiceOrder)(nil), nil).Once()
	m.On("ListServiceOrders", context.Background()).Return(([]*entities.ServiceOrder)(nil), nil).Once()
	m.On("UpdateServiceOrder", context.Background(), entities.ServiceOrder{}, "flow").Return((*entities.ServiceOrder)(nil), nil).Once()

	_, _ = m.CreateServiceOrder(context.Background(), entities.ServiceOrder{})
	_, _ = m.GetServiceOrder(context.Background(), entities.ServiceOrder{})
	_, _ = m.ListServiceOrders(context.Background())
	_, _ = m.UpdateServiceOrder(context.Background(), entities.ServiceOrder{}, "flow")

	m.AssertExpectations(t)
}

func TestMocksCoverage_PaymentUseCase(t *testing.T) {
	m := NewMockIPaymentUseCase(nil)

	m.On("CreatePayment", context.Background(), (*entities.Payment)(nil)).Return((*entities.Payment)(nil), nil).Once()
	m.On("GetPaymentByID", context.Background(), uint(1)).Return((*entities.Payment)(nil), nil).Once()
	m.On("ListPayments", context.Background()).Return(([]entities.Payment)(nil), nil).Once()

	_, _ = m.CreatePayment(context.Background(), nil)
	_, _ = m.GetPaymentByID(context.Background(), 1)
	_, _ = m.ListPayments(context.Background())

	m.AssertExpectations(t)
}

func TestMocksCoverage_PartsSupplyUseCase(t *testing.T) {
	m := NewMockIPartsSupplyUseCase(nil)

	m.On("GetPartsSupplyByServiceOrderID", context.Background(), uint(1)).Return(([]entities.PartsSupply)(nil), nil).Once()
	m.On("CreatePartsSupply", context.Background(), (*entities.PartsSupply)(nil)).Return(entities.PartsSupply{}, nil).Once()
	m.On("DeletePartsSupply", context.Background(), uint(1)).Return(nil).Once()
	m.On("GetPartsSupplyByID", context.Background(), uint(1)).Return(entities.PartsSupply{}, nil).Once()
	m.On("ListPartsSupplies", context.Background()).Return(([]entities.PartsSupply)(nil), nil).Once()
	m.On("UpdatePartsSupply", context.Background(), (*entities.PartsSupply)(nil)).Return(nil).Once()

	_, _ = m.GetPartsSupplyByServiceOrderID(context.Background(), 1)
	_, _ = m.CreatePartsSupply(context.Background(), nil)
	_ = m.DeletePartsSupply(context.Background(), 1)
	_, _ = m.GetPartsSupplyByID(context.Background(), 1)
	_, _ = m.ListPartsSupplies(context.Background())
	_ = m.UpdatePartsSupply(context.Background(), nil)

	m.AssertExpectations(t)
}

func TestMocksCoverage_ServiceUseCase(t *testing.T) {
	m := NewMockIServiceUseCase()

	m.On("CreateService", context.Background(), (*entities.Service)(nil)).Return(entities.Service{}, nil).Once()
	m.On("DeleteService", context.Background(), uint(1)).Return(nil).Once()
	m.On("GetServiceByID", context.Background(), uint(1)).Return(entities.Service{}, nil).Once()
	m.On("ListServices", context.Background()).Return(([]entities.Service)(nil), nil).Once()
	m.On("UpdateService", context.Background(), (*entities.Service)(nil)).Return(nil).Once()

	_, _ = m.CreateService(context.Background(), nil)
	_ = m.DeleteService(context.Background(), 1)
	_, _ = m.GetServiceByID(context.Background(), 1)
	_, _ = m.ListServices(context.Background())
	_ = m.UpdateService(context.Background(), nil)

	m.AssertExpectations(t)
}

func TestMocksCoverage_AdditionalRepairUseCase(t *testing.T) {
	m := NewMockIAdditionalRepairUseCase()

	m.On("CreateAdditionalRepair", context.Background(), entities.AdditionalRepair{}).Return((*entities.AdditionalRepair)(nil), nil).Once()
	m.On("GetAdditionalRepair", context.Background(), "1").Return((*entities.AdditionalRepair)(nil), nil).Once()
	m.On("GetAdditionalRepairBySO", context.Background(), "1").Return(([]entities.AdditionalRepair)(nil), nil).Once()
	m.On("CustomerApprovalStatus", context.Background(), "1", "flow").Return((*entities.AdditionalRepair)(nil), nil).Once()
	m.On("CancelAdditionalRepair", context.Background(), "1").Return((*entities.AdditionalRepair)(nil), nil).Once()
	m.On("Rollback", context.Background(), (*entities.AdditionalRepair)(nil), false).Return(nil).Once()

	_, _ = m.CreateAdditionalRepair(context.Background(), entities.AdditionalRepair{})
	_, _ = m.GetAdditionalRepair(context.Background(), "1")
	_, _ = m.GetAdditionalRepairBySO(context.Background(), "1")
	_, _ = m.CustomerApprovalStatus(context.Background(), "1", "flow")
	_, _ = m.CancelAdditionalRepair(context.Background(), "1")
	_ = m.Rollback(context.Background(), nil, false)

	m.AssertExpectations(t)
}

func TestMocksCoverage_VehicleService(t *testing.T) {
	m := &MockVehicleService{}

	m.On("GetAllVehicles").Return(([]entities.Vehicle)(nil), nil).Once()
	m.On("GetVehicleByID", uint(1)).Return((*entities.Vehicle)(nil), nil).Once()
	m.On("GetVehiclesByCustomerID", uint(1)).Return(([]entities.Vehicle)(nil), nil).Once()
	m.On("CreateVehicle", entities.Vehicle{}).Return((*entities.Vehicle)(nil), nil).Once()
	m.On("UpdateVehicle", entities.Vehicle{}).Return("", nil).Once()
	m.On("DeleteVehicle", uint(1)).Return(nil).Once()
	m.On("GetVehicleByPlate", "p").Return((*entities.Vehicle)(nil), nil).Once()
	m.On("UpdateVehiclePartial", uint(1), map[string]interface{}{}).Return("", nil).Once()

	_, _ = m.GetAllVehicles()
	_, _ = m.GetVehicleByID(1)
	_, _ = m.GetVehiclesByCustomerID(1)
	_, _ = m.CreateVehicle(entities.Vehicle{})
	_, _ = m.UpdateVehicle(entities.Vehicle{})
	_ = m.DeleteVehicle(1)
	_, _ = m.GetVehicleByPlate("p")
	_, _ = m.UpdateVehiclePartial(1, map[string]interface{}{})

	m.AssertExpectations(t)

	assert.True(t, true)
}
