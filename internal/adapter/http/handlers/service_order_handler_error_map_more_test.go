package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceOrderHandler_ErrorMapping_MoreCases(t *testing.T) {
	mockUC, h, r := setupServiceOrderHandlerTest(t)
	gin.SetMode(gin.TestMode)

	r.POST("/os", h.CreateServiceOrder)
	r.POST("/os/:id/diagnosis", h.DiagnosisServiceOrder)
	r.POST("/os/:id/estimate/approve", h.ApproveServiceOrderEstimate)
	r.POST("/os/:id/execution/create", h.ExecutionServiceOrder)
	r.POST("/os/:id/delivery", h.DeliveryServiceOrder)

	t.Run("create - customer not found -> 404", func(t *testing.T) {
		mockUC.On("CreateServiceOrder", mock.Anything, mock.Anything).Return((*entities.ServiceOrder)(nil), usecase.ErrCustomerNotFound).Once()
		req, _ := http.NewRequest("POST", "/os", bytes.NewBufferString(`{"customer_id":"1","vehicle_id":"1"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("diagnosis - invalid transition -> 400", func(t *testing.T) {
		mockUC.On("DiagnosisServiceOrder", mock.Anything, mock.Anything).Return((*entities.ServiceOrder)(nil), usecase.ErrInvalidTransitionStatusToDiagnosis).Once()
		req, _ := http.NewRequest("POST", "/os/1/diagnosis", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("diagnosis - service not found -> 404", func(t *testing.T) {
		mockUC.On("DiagnosisServiceOrder", mock.Anything, mock.Anything).Return((*entities.ServiceOrder)(nil), usecase.ErrServiceNotFound).Once()
		req, _ := http.NewRequest("POST", "/os/2/diagnosis", bytes.NewBufferString(`{"services":[{"id":"1"}],"parts_supplies":[{"id":"p1","quantity":1}]}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("diagnosis - parts supply not found -> 404", func(t *testing.T) {
		mockUC.On("DiagnosisServiceOrder", mock.Anything, mock.Anything).Return((*entities.ServiceOrder)(nil), usecase.ErrPartsSupplyNotFound).Once()
		req, _ := http.NewRequest("POST", "/os/3/diagnosis", bytes.NewBufferString(`{"services":[{"id":"1"}],"parts_supplies":[{"id":"p1","quantity":1}]}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("estimate - invalid transition -> 400", func(t *testing.T) {
		mockUC.On("EstimateServiceOrder", mock.Anything, "4", mock.Anything).Return((*entities.ServiceOrder)(nil), usecase.ErrInvalidTransitionStatusToEstimate).Once()
		req, _ := http.NewRequest("POST", "/os/4/estimate/approve", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("execution - invalid transition -> 400", func(t *testing.T) {
		mockUC.On("ExecutionServiceOrder", mock.Anything, "5", mock.Anything).Return((*entities.ServiceOrder)(nil), usecase.ErrInvalidTransitionStatusToExecution).Once()
		req, _ := http.NewRequest("POST", "/os/5/execution/create", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("delivery - invalid transition -> 400", func(t *testing.T) {
		mockUC.On("DeliveryServiceOrder", mock.Anything, "6").Return((*entities.ServiceOrder)(nil), usecase.ErrInvalidTransitionStatusToDelivery).Once()
		req, _ := http.NewRequest("POST", "/os/6/delivery", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	mockUC.AssertExpectations(t)
}
