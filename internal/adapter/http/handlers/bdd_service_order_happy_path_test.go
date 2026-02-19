package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	request "github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/request"
	handlers "github.com/fiap-grupo95/os-service-api/internal/adapter/http/handlers"
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
	"github.com/fiap-grupo95/os-service-api/internal/usecase"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type bddServiceOrderFlowSuite struct {
	r               *gin.Engine
	serviceOrderGW  *mocks.MockServiceOrderGateway
	vehicleGW       *mocks.MockVehicleGateway
	customerGW      *mocks.MockCustomerGateway
	serviceGW       *mocks.MockServiceGateway
	partsSupplyGW   *mocks.MockPartsSupplyGateway
	billingGW       *mocks.MockBillingServiceGateway
	executionGW     *mocks.MockExecutionGateway
	serviceOrderUC  *usecase.ServiceOrderUseCase
	serviceOrderHdl *handlers.ServiceOrderHandler

	serviceOrderID string
	estimateID     string
	paymentID      string
}

func setupBDDServiceOrderFlowSuite(t *testing.T) *bddServiceOrderFlowSuite {
	t.Helper()
	gin.SetMode(gin.TestMode)

	s := &bddServiceOrderFlowSuite{}
	s.r = gin.New()

	s.serviceOrderGW = mocks.NewMockServiceOrderGateway()
	s.vehicleGW = &mocks.MockVehicleGateway{}
	s.customerGW = &mocks.MockCustomerGateway{}
	s.serviceGW = &mocks.MockServiceGateway{}
	s.partsSupplyGW = &mocks.MockPartsSupplyGateway{}
	s.billingGW = &mocks.MockBillingServiceGateway{}
	s.executionGW = &mocks.MockExecutionGateway{}

	s.serviceOrderUC = usecase.NewServiceOrderUseCase(
		s.serviceOrderGW,
		s.vehicleGW,
		s.customerGW,
		s.serviceGW,
		s.partsSupplyGW,
		s.billingGW,
		s.executionGW,
	)
	s.serviceOrderHdl = handlers.NewServiceOrderHandler(s.serviceOrderUC)

	// Routes only for service order flow (no auth middleware)
	s.r.POST("/v1/service-orders/create", s.serviceOrderHdl.CreateServiceOrder)
	s.r.POST("/v1/service-orders/:id/diagnosis", s.serviceOrderHdl.DiagnosisServiceOrder)
	s.r.POST("/v1/service-orders/:id/estimate/approve", s.serviceOrderHdl.ApproveServiceOrderEstimate)
	s.r.POST("/v1/service-orders/:id/estimate/reject", s.serviceOrderHdl.RejectServiceOrderEstimate)
	s.r.POST("/v1/service-orders/:id/execution/create", s.serviceOrderHdl.ExecutionServiceOrder)
	s.r.POST("/v1/service-orders/:id/execution/finish", s.serviceOrderHdl.FinishServiceOrderExecution)
	s.r.POST("/v1/service-orders/:id/payment", s.serviceOrderHdl.PaymentServiceOrder)
	s.r.POST("/v1/service-orders/:id/delivery", s.serviceOrderHdl.DeliveryServiceOrder)

	return s
}

func (s *bddServiceOrderFlowSuite) stepCreateServiceOrder(t *testing.T) {
	t.Helper()

	const (
		soID       = "so-1"
		vehicleID  = "veh-1"
		customerID = "cus-1"
	)

	s.serviceOrderID = soID

	s.vehicleGW.On("FindByID", vehicleID).Return(&entities.Vehicle{ID: vehicleID}, nil).Once()
	s.customerGW.On("GetByID", customerID).Return(&entities.Customer{ID: customerID}, nil).Once()

	s.serviceOrderGW.On(
		"Create",
		mock.Anything,
		mock.MatchedBy(func(so *entities.ServiceOrder) bool {
			return so != nil && so.VehicleID == vehicleID && so.CustomerID == customerID && so.Status.IsRecebida()
		}),
	).Return(&entities.ServiceOrder{ID: soID, VehicleID: vehicleID, CustomerID: customerID, Status: valueobject.StatusRecebida}, nil).Once()

	payload := request.ServiceOrderCreateRequest{CustomerID: customerID, VehicleID: vehicleID}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/v1/service-orders/create", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func (s *bddServiceOrderFlowSuite) stepDiagnosisWithServicesAndParts(t *testing.T) {
	t.Helper()

	const estimateID = "est-1"
	s.estimateID = estimateID

	serviceItems := []entities.Service{{ID: "svc-1"}, {ID: "svc-2"}}
	partsItems := []entities.PartsSupply{{ID: "part-1", Quantity: 2}, {ID: "part-2", Quantity: 1}}

	s.serviceOrderGW.On("GetByID", mock.Anything, s.serviceOrderID, false).Return(&entities.ServiceOrder{ID: s.serviceOrderID, Status: valueobject.StatusRecebida}, nil).Once()

	for _, svc := range serviceItems {
		copySvc := svc
		s.serviceGW.On("GetByID", mock.Anything, copySvc.ID).Return(&entities.Service{ID: copySvc.ID}, nil).Once()
	}

	s.partsSupplyGW.On("AuthorizeReserve", mock.Anything, mock.Anything).Return(nil).Once()
	s.partsSupplyGW.On("Reserve", mock.Anything, mock.Anything).Return(nil).Once()

	s.billingGW.On("CreateEstimate", mock.Anything, mock.Anything, (*entities.AdditionalRepair)(nil)).Return(&entities.Estimate{ID: estimateID, Value: 100}, nil).Once()

	s.serviceOrderGW.On(
		"Update",
		mock.Anything,
		mock.MatchedBy(func(so *entities.ServiceOrder) bool {
			return so != nil && so.ID == s.serviceOrderID && so.Status.IsAguardandoAprovacao() && so.Estimate != nil && so.Estimate.ID == estimateID
		}),
	).Return(nil, nil).Once()

	payload := request.ServiceOrderDiagnosisUpdateRequest{
		Services:      []request.ServiceOrderServiceItem{{ID: serviceItems[0].ID}, {ID: serviceItems[1].ID}},
		PartsSupplies: []request.ServiceOrderPartsSupplyItem{{ID: partsItems[0].ID, Quantity: partsItems[0].Quantity}, {ID: partsItems[1].ID, Quantity: partsItems[1].Quantity}},
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/v1/service-orders/"+s.serviceOrderID+"/diagnosis", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func (s *bddServiceOrderFlowSuite) stepDiagnosisPending(t *testing.T) {
	t.Helper()

	s.serviceOrderGW.On("GetByID", mock.Anything, s.serviceOrderID, false).Return(
		&entities.ServiceOrder{ID: s.serviceOrderID, Status: valueobject.StatusRecebida},
		nil,
	).Once()

	s.partsSupplyGW.On("Reserve", mock.Anything, mock.Anything).Return(nil).Once()

	s.serviceOrderGW.On(
		"Update",
		mock.Anything,
		mock.MatchedBy(func(so *entities.ServiceOrder) bool {
			return so != nil && so.ID == s.serviceOrderID && so.Status.IsEmDiagnostico()
		}),
	).Return(nil, nil).Once()

	req, _ := http.NewRequest(http.MethodPost, "/v1/service-orders/"+s.serviceOrderID+"/diagnosis", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func (s *bddServiceOrderFlowSuite) stepRejectEstimate(t *testing.T) {
	t.Helper()

	s.serviceOrderGW.On("GetByID", mock.Anything, s.serviceOrderID, false).Return(
		&entities.ServiceOrder{ID: s.serviceOrderID, Status: valueobject.StatusAguardandoAprovacao, Estimate: &entities.Estimate{ID: s.estimateID}},
		nil,
	).Once()

	s.partsSupplyGW.On("Release", mock.Anything, mock.Anything).Return(nil).Once()
	s.billingGW.On("RejectEstimate", mock.Anything, mock.Anything, (*entities.AdditionalRepair)(nil)).Return(
		&entities.Estimate{ID: s.estimateID, Value: 100},
		nil,
	).Once()

	s.serviceOrderGW.On("Update", mock.Anything, mock.MatchedBy(func(so *entities.ServiceOrder) bool {
		return so != nil && so.ID == s.serviceOrderID && so.Status.IsRejeitada()
	})).Return(nil, nil).Once()

	req, _ := http.NewRequest(http.MethodPost, "/v1/service-orders/"+s.serviceOrderID+"/estimate/reject", nil)
	w := httptest.NewRecorder()
	s.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func (s *bddServiceOrderFlowSuite) stepApproveEstimate(t *testing.T) {
	t.Helper()
	s.serviceOrderGW.On("GetByID", mock.Anything, s.serviceOrderID, false).Return(&entities.ServiceOrder{ID: s.serviceOrderID, Status: valueobject.StatusAguardandoAprovacao, Estimate: &entities.Estimate{ID: s.estimateID}}, nil).Once()
	// Estimate approve strategy will call billing gateway approve and then parts supply write-off
	s.billingGW.On("ApproveEstimate", mock.Anything, mock.Anything, (*entities.AdditionalRepair)(nil)).Return(&entities.Estimate{ID: s.estimateID, Value: 100}, nil).Once()
	s.partsSupplyGW.On("WriteOff", mock.Anything, mock.Anything).Return(nil).Once()

	s.serviceOrderGW.On("Update", mock.Anything, mock.MatchedBy(func(so *entities.ServiceOrder) bool {
		return so != nil && so.ID == s.serviceOrderID && so.Status.IsAprovada()
	})).Return(nil, nil).Once()

	req, _ := http.NewRequest(http.MethodPost, "/v1/service-orders/"+s.serviceOrderID+"/estimate/approve", nil)
	w := httptest.NewRecorder()
	s.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func (s *bddServiceOrderFlowSuite) stepStartExecution(t *testing.T) {
	t.Helper()
	s.serviceOrderGW.On("GetByID", mock.Anything, s.serviceOrderID, false).Return(&entities.ServiceOrder{ID: s.serviceOrderID, Status: valueobject.StatusAprovada, Estimate: &entities.Estimate{ID: s.estimateID}}, nil).Once()
	s.executionGW.On("CreateExecution", mock.Anything, mock.Anything).Return(&entities.Execution{ID: "exec-1"}, nil).Once()
	s.serviceOrderGW.On("Update", mock.Anything, mock.MatchedBy(func(so *entities.ServiceOrder) bool {
		return so != nil && so.ID == s.serviceOrderID && so.Status.IsEmExecucao() && so.Execution != nil && so.Execution.ID == "exec-1"
	})).Return(nil, nil).Once()

	req, _ := http.NewRequest(http.MethodPost, "/v1/service-orders/"+s.serviceOrderID+"/execution/create", nil)
	w := httptest.NewRecorder()
	s.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func (s *bddServiceOrderFlowSuite) stepFinishExecution(t *testing.T) {
	t.Helper()
	s.serviceOrderGW.On("GetByID", mock.Anything, s.serviceOrderID, false).Return(&entities.ServiceOrder{ID: s.serviceOrderID, Status: valueobject.StatusEmExecucao, Estimate: &entities.Estimate{ID: s.estimateID}, Execution: &entities.Execution{ID: "exec-1"}}, nil).Once()
	s.executionGW.On("FinishExecution", mock.Anything, mock.Anything).Return(&entities.Execution{ID: "exec-1"}, nil).Once()
	s.serviceOrderGW.On("Update", mock.Anything, mock.MatchedBy(func(so *entities.ServiceOrder) bool {
		return so != nil && so.ID == s.serviceOrderID && so.Status.IsFinalizada()
	})).Return(nil, nil).Once()

	req, _ := http.NewRequest(http.MethodPost, "/v1/service-orders/"+s.serviceOrderID+"/execution/finish", nil)
	w := httptest.NewRecorder()
	s.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func (s *bddServiceOrderFlowSuite) stepCreatePayment(t *testing.T) {
	t.Helper()
	const paymentID = "pay-1"
	s.paymentID = paymentID

	s.serviceOrderGW.On("GetByID", mock.Anything, s.serviceOrderID, false).Return(&entities.ServiceOrder{ID: s.serviceOrderID, Status: valueobject.StatusFinalizada, Estimate: &entities.Estimate{ID: s.estimateID}}, nil).Once()
	s.billingGW.On("CreatePayment", mock.Anything, s.estimateID).Return(&entities.Payment{ID: paymentID, EstimateID: s.estimateID, Amount: 100}, nil).Once()

	req, _ := http.NewRequest(http.MethodPost, "/v1/service-orders/"+s.serviceOrderID+"/payment", nil)
	w := httptest.NewRecorder()
	s.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func (s *bddServiceOrderFlowSuite) stepDeliverServiceOrder(t *testing.T) {
	t.Helper()
	s.serviceOrderGW.On("GetByID", mock.Anything, s.serviceOrderID, false).Return(&entities.ServiceOrder{ID: s.serviceOrderID, Status: valueobject.StatusFinalizada, Estimate: &entities.Estimate{ID: s.estimateID}}, nil).Once()
	s.billingGW.On("GetPaymentByEstimateID", mock.Anything, s.estimateID).Return(&entities.Payment{ID: s.paymentID, EstimateID: s.estimateID, Amount: 100}, nil).Once()
	s.serviceOrderGW.On("Update", mock.Anything, mock.MatchedBy(func(so *entities.ServiceOrder) bool {
		return so != nil && so.ID == s.serviceOrderID && so.Status.IsEntregue()
	})).Return(nil, nil).Once()

	req, _ := http.NewRequest(http.MethodPost, "/v1/service-orders/"+s.serviceOrderID+"/delivery", nil)
	w := httptest.NewRecorder()
	s.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBDD_ServiceOrder_HappyPath_FullFlow(t *testing.T) {
	s := setupBDDServiceOrderFlowSuite(t)

	s.stepCreateServiceOrder(t)
	s.stepDiagnosisWithServicesAndParts(t)
	s.stepApproveEstimate(t)
	s.stepStartExecution(t)
	s.stepFinishExecution(t)
	s.stepCreatePayment(t)
	s.stepDeliverServiceOrder(t)

	s.vehicleGW.AssertExpectations(t)
	s.customerGW.AssertExpectations(t)
	s.serviceGW.AssertExpectations(t)
	s.partsSupplyGW.AssertExpectations(t)
	s.billingGW.AssertExpectations(t)
	s.executionGW.AssertExpectations(t)
	s.serviceOrderGW.AssertExpectations(t)
}

func TestBDD_ServiceOrder_RejectFlow_Create_DiagnosisTwoSteps_Reject(t *testing.T) {
	s := setupBDDServiceOrderFlowSuite(t)

	s.stepCreateServiceOrder(t)
	s.stepDiagnosisPending(t)
	s.stepDiagnosisWithServicesAndParts(t)
	s.stepRejectEstimate(t)

	s.vehicleGW.AssertExpectations(t)
	s.customerGW.AssertExpectations(t)
	s.serviceGW.AssertExpectations(t)
	s.partsSupplyGW.AssertExpectations(t)
	s.billingGW.AssertExpectations(t)
	s.executionGW.AssertExpectations(t)
	s.serviceOrderGW.AssertExpectations(t)
}
