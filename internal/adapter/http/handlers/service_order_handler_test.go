package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockServiceOrderUseCase struct {
	mock.Mock
}

func (m *MockServiceOrderUseCase) CreateServiceOrder(ctx context.Context, serviceOrder entities.ServiceOrder) (*entities.ServiceOrder, error) {
	args := m.Called(ctx, serviceOrder)
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
	jsonBody := `{"customer_id":1,"vehicle_id":1}`
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

func TestGetServiceOrder(t *testing.T) {
	mockUC, h, r := setupServiceOrderHandlerTest(t)
	r.GET("/os/:id", h.GetServiceOrder)

	mockUC.On("GetServiceOrder", mock.Anything, entities.ServiceOrder{ID: 1}, false).Return(&entities.ServiceOrder{ID: 1}, nil).Once()
	req, _ := http.NewRequest("GET", "/os/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mockUC.On("GetServiceOrder", mock.Anything, entities.ServiceOrder{ID: 2}, false).Return((*entities.ServiceOrder)(nil), errors.New("fail")).Once()
	req, _ = http.NewRequest("GET", "/os/2", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	req, _ = http.NewRequest("GET", "/os/abc", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockUC.AssertExpectations(t)
}

func TestListServiceOrders(t *testing.T) {
	mockUC, h, r := setupServiceOrderHandlerTest(t)
	r.GET("/os", h.ListServiceOrders)
	mockUC.On("ListServiceOrders", mock.Anything).Return([]*entities.ServiceOrder{{ID: 1}}, nil).Once()
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
