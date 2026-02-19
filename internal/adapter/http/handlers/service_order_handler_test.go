package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockServiceOrderUseCase struct {
	mock.Mock
}

func (m *MockServiceOrderUseCase) CreateServiceOrder(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error) {
	args := m.Called(ctx, serviceOrder)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ServiceOrder), args.Error(1)
}

func (m *MockServiceOrderUseCase) DiagnosisServiceOrder(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error) {
	args := m.Called(ctx, serviceOrder)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ServiceOrder), args.Error(1)
}

func (m *MockServiceOrderUseCase) EstimateServiceOrder(ctx context.Context, serviceOrderID string, operation string) (*entities.ServiceOrder, error) {
	args := m.Called(ctx, serviceOrderID, operation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ServiceOrder), args.Error(1)
}

func (m *MockServiceOrderUseCase) ExecutionServiceOrder(ctx context.Context, serviceOrderID string, operation string) (*entities.ServiceOrder, error) {
	args := m.Called(ctx, serviceOrderID, operation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ServiceOrder), args.Error(1)
}

func (m *MockServiceOrderUseCase) PaymentServiceOrder(ctx context.Context, serviceOrderID string) (*entities.Payment, error) {
	args := m.Called(ctx, serviceOrderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Payment), args.Error(1)
}

func (m *MockServiceOrderUseCase) DeliveryServiceOrder(ctx context.Context, serviceOrderID string) (*entities.ServiceOrder, error) {
	args := m.Called(ctx, serviceOrderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ServiceOrder), args.Error(1)
}

func (m *MockServiceOrderUseCase) CancelServiceOrder(ctx context.Context, serviceOrderID string) (*entities.ServiceOrder, error) {
	args := m.Called(ctx, serviceOrderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ServiceOrder), args.Error(1)
}

func (m *MockServiceOrderUseCase) UpdateServiceOrder(ctx context.Context, serviceOrder entities.ServiceOrder, flow string) (*entities.ServiceOrder, error) {
	args := m.Called(ctx, serviceOrder, flow)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ServiceOrder), args.Error(1)
}

func (m *MockServiceOrderUseCase) GetServiceOrder(ctx context.Context, serviceOrder entities.ServiceOrder, isFullData bool) (*entities.ServiceOrder, error) {
	args := m.Called(ctx, serviceOrder, isFullData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ServiceOrder), args.Error(1)
}

func (m *MockServiceOrderUseCase) ListServiceOrders(ctx context.Context) ([]*entities.ServiceOrder, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.ServiceOrder), args.Error(1)
}

func setupServiceOrderHandlerTest(t *testing.T) (*MockServiceOrderUseCase, *ServiceOrderHandler, *gin.Engine) {
	mockUC := new(MockServiceOrderUseCase)
	h := &ServiceOrderHandler{serviceOrderUseCase: mockUC}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return mockUC, h, r
}

func TestCreateServiceOrder(t *testing.T) {
	mockUC, h, r := setupServiceOrderHandlerTest(t)
	r.POST("/os", h.CreateServiceOrder)

	// Success
	mockUC.On("CreateServiceOrder", mock.Anything, mock.Anything).Return(&entities.ServiceOrder{}, nil).Once()
	jsonBody := `{"customer_id":"1","vehicle_id":"1"}`
	req, _ := http.NewRequest("POST", "/os", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Invalid input
	req, _ = http.NewRequest("POST", "/os", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Failure
	mockUC.On("CreateServiceOrder", mock.Anything, mock.Anything).Return((*entities.ServiceOrder)(nil), errors.New("fail")).Once()
	req, _ = http.NewRequest("POST", "/os", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockUC.AssertExpectations(t)
}

func TestEstimateHandlers_ServiceOrder(t *testing.T) {
	mockUC, h, r := setupServiceOrderHandlerTest(t)
	r.POST("/os/:id/estimate/approve", h.ApproveServiceOrderEstimate)
	r.POST("/os/:id/estimate/reject", h.RejectServiceOrderEstimate)
	r.POST("/os/:id/estimate/cancel", h.CancelServiceOrderEstimate)

	t.Run("approve - not found", func(t *testing.T) {
		mockUC.On("EstimateServiceOrder", mock.Anything, "1", mock.Anything).Return((*entities.ServiceOrder)(nil), usecase.ErrServiceOrderNotFound).Once()
		req, _ := http.NewRequest("POST", "/os/1/estimate/approve", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("reject - bad request", func(t *testing.T) {
		mockUC.On("EstimateServiceOrder", mock.Anything, "2", mock.Anything).Return((*entities.ServiceOrder)(nil), usecase.ErrInvalidTransitionStatusToEstimate).Once()
		req, _ := http.NewRequest("POST", "/os/2/estimate/reject", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("cancel - internal error", func(t *testing.T) {
		mockUC.On("EstimateServiceOrder", mock.Anything, "3", mock.Anything).Return((*entities.ServiceOrder)(nil), errors.New("fail")).Once()
		req, _ := http.NewRequest("POST", "/os/3/estimate/cancel", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	mockUC.AssertExpectations(t)
}

func TestExecutionHandlers_ServiceOrder(t *testing.T) {
	mockUC, h, r := setupServiceOrderHandlerTest(t)
	r.POST("/os/:id/execution/create", h.ExecutionServiceOrder)
	r.POST("/os/:id/execution/finish", h.FinishServiceOrderExecution)

	t.Run("start - bad request", func(t *testing.T) {
		mockUC.On("ExecutionServiceOrder", mock.Anything, "1", mock.Anything).Return((*entities.ServiceOrder)(nil), usecase.ErrInvalidTransitionStatusToExecution).Once()
		req, _ := http.NewRequest("POST", "/os/1/execution/create", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("finish - not found", func(t *testing.T) {
		mockUC.On("ExecutionServiceOrder", mock.Anything, "2", mock.Anything).Return((*entities.ServiceOrder)(nil), usecase.ErrServiceOrderNotFound).Once()
		req, _ := http.NewRequest("POST", "/os/2/execution/finish", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("finish - internal error", func(t *testing.T) {
		mockUC.On("ExecutionServiceOrder", mock.Anything, "3", mock.Anything).Return((*entities.ServiceOrder)(nil), errors.New("fail")).Once()
		req, _ := http.NewRequest("POST", "/os/3/execution/finish", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	mockUC.AssertExpectations(t)
}

func TestPaymentAndDeliveryHandlers_ServiceOrder(t *testing.T) {
	mockUC, h, r := setupServiceOrderHandlerTest(t)
	r.POST("/os/:id/payment", h.PaymentServiceOrder)
	r.POST("/os/:id/delivery", h.DeliveryServiceOrder)

	t.Run("payment - bad request", func(t *testing.T) {
		mockUC.On("PaymentServiceOrder", mock.Anything, "1").Return((*entities.Payment)(nil), usecase.ErrInvalidTransitionStatusToDelivery).Once()
		req, _ := http.NewRequest("POST", "/os/1/payment", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("payment - internal error", func(t *testing.T) {
		mockUC.On("PaymentServiceOrder", mock.Anything, "2").Return((*entities.Payment)(nil), errors.New("fail")).Once()
		req, _ := http.NewRequest("POST", "/os/2/payment", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("delivery - not found", func(t *testing.T) {
		mockUC.On("DeliveryServiceOrder", mock.Anything, "3").Return((*entities.ServiceOrder)(nil), usecase.ErrServiceOrderNotFound).Once()
		req, _ := http.NewRequest("POST", "/os/3/delivery", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("delivery - bad request", func(t *testing.T) {
		mockUC.On("DeliveryServiceOrder", mock.Anything, "4").Return((*entities.ServiceOrder)(nil), usecase.ErrInvalidTransitionStatusToDelivery).Once()
		req, _ := http.NewRequest("POST", "/os/4/delivery", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("delivery - internal error", func(t *testing.T) {
		mockUC.On("DeliveryServiceOrder", mock.Anything, "5").Return((*entities.ServiceOrder)(nil), errors.New("fail")).Once()
		req, _ := http.NewRequest("POST", "/os/5/delivery", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	mockUC.AssertExpectations(t)
}

func TestCancelServiceOrder(t *testing.T) {
	mockUC, h, r := setupServiceOrderHandlerTest(t)
	r.POST("/os/:id/cancel", h.CancelServiceOrder)

	// Success
	mockUC.On("CancelServiceOrder", mock.Anything, "1").Return(&entities.ServiceOrder{ID: "1"}, nil).Once()
	req, _ := http.NewRequest("POST", "/os/1/cancel", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Failure
	mockUC.On("CancelServiceOrder", mock.Anything, "2").Return((*entities.ServiceOrder)(nil), errors.New("fail")).Once()
	req, _ = http.NewRequest("POST", "/os/2/cancel", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockUC.AssertExpectations(t)
}

func TestCancelServiceOrderEstimate(t *testing.T) {
	mockUC, h, r := setupServiceOrderHandlerTest(t)
	r.POST("/os/:id/estimate/cancel", h.CancelServiceOrderEstimate)

	mockUC.On("EstimateServiceOrder", mock.Anything, "1", mock.Anything).Return(&entities.ServiceOrder{ID: "1"}, nil).Once()
	req, _ := http.NewRequest("POST", "/os/1/estimate/cancel", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mockUC.On("EstimateServiceOrder", mock.Anything, "2", mock.Anything).Return((*entities.ServiceOrder)(nil), errors.New("fail")).Once()
	req, _ = http.NewRequest("POST", "/os/2/estimate/cancel", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockUC.AssertExpectations(t)
}

func TestDiagnosisServiceOrder_InvalidJSON(t *testing.T) {
	_, h, r := setupServiceOrderHandlerTest(t)
	r.POST("/os/:id/diagnosis", h.DiagnosisServiceOrder)

	req, _ := http.NewRequest("POST", "/os/1/diagnosis", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetServiceOrderFullData_QueryParam(t *testing.T) {
	mockUC, h, r := setupServiceOrderHandlerTest(t)
	r.GET("/os/:id", h.GetServiceOrder)

	so := &entities.ServiceOrder{ID: "1"}
	mockUC.On("GetServiceOrder", mock.Anything, entities.ServiceOrder{ID: "1"}, true).Return(so, nil).Once()

	req, _ := http.NewRequest("GET", "/os/1?full_data=true", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var got map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &got)

	mockUC.AssertExpectations(t)
}

func TestGetServiceOrder(t *testing.T) {
	mockUC, h, r := setupServiceOrderHandlerTest(t)
	r.GET("/os/:id", h.GetServiceOrder)

	mockUC.On("GetServiceOrder", mock.Anything, entities.ServiceOrder{ID: "1"}, false).Return(&entities.ServiceOrder{ID: "1"}, nil).Once()
	req, _ := http.NewRequest("GET", "/os/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mockUC.On("GetServiceOrder", mock.Anything, entities.ServiceOrder{ID: "2"}, false).Return((*entities.ServiceOrder)(nil), errors.New("fail")).Once()
	req, _ = http.NewRequest("GET", "/os/2", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockUC.On("GetServiceOrder", mock.Anything, entities.ServiceOrder{ID: "abc"}, false).Return((*entities.ServiceOrder)(nil), errors.New("fail")).Once()
	req, _ = http.NewRequest("GET", "/os/abc", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockUC.AssertExpectations(t)
}

func TestGetServiceOrderFullData(t *testing.T) {
	mockUC, h, r := setupServiceOrderHandlerTest(t)
	r.GET("/os/:id", h.GetServiceOrder)

	mockUC.On("GetServiceOrder", mock.Anything, entities.ServiceOrder{ID: "1"}, true).Return(&entities.ServiceOrder{ID: "1", Customer: &entities.Customer{ID: "1"}, Vehicle: &entities.Vehicle{ID: "1"}, PartsSupplies: []entities.PartsSupply{{ID: "1"}}, Services: []entities.Service{{ID: "1"}}}, nil).Once()
	req, _ := http.NewRequest("GET", "/os/1?full_data=true", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mockUC.On("GetServiceOrder", mock.Anything, entities.ServiceOrder{ID: "2"}, true).Return((*entities.ServiceOrder)(nil), errors.New("fail")).Once()
	req, _ = http.NewRequest("GET", "/os/2?full_data=true", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockUC.On("GetServiceOrder", mock.Anything, entities.ServiceOrder{ID: "abc"}, true).Return((*entities.ServiceOrder)(nil), errors.New("fail")).Once()
	req, _ = http.NewRequest("GET", "/os/abc?full_data=true", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockUC.AssertExpectations(t)
}

func TestListServiceOrders(t *testing.T) {
	mockUC, h, r := setupServiceOrderHandlerTest(t)
	r.GET("/os", h.ListServiceOrders)
	mockUC.On("ListServiceOrders", mock.Anything).Return([]*entities.ServiceOrder{{ID: "1"}}, nil).Once()
	req, _ := http.NewRequest("GET", "/os", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mockUC.On("ListServiceOrders", mock.Anything).Return(([]*entities.ServiceOrder)(nil), errors.New("fail")).Once()
	req, _ = http.NewRequest("GET", "/os", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockUC.AssertExpectations(t)
}
