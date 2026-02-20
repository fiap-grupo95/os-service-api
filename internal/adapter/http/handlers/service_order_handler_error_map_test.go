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

func TestServiceOrderHandler_ErrorMapping_TableDriven(t *testing.T) {
	mockUC, h, r := setupServiceOrderHandlerTest(t)
	gin.SetMode(gin.TestMode)

	r.POST("/os", h.CreateServiceOrder)
	r.POST("/os/:id/diagnosis", h.DiagnosisServiceOrder)
	r.POST("/os/:id/estimate/approve", h.ApproveServiceOrderEstimate)
	r.POST("/os/:id/execution/create", h.ExecutionServiceOrder)
	r.POST("/os/:id/payment", h.PaymentServiceOrder)
	r.POST("/os/:id/delivery", h.DeliveryServiceOrder)

	t.Run("CreateServiceOrder - vehicle not found -> 404", func(t *testing.T) {
		mockUC.On("CreateServiceOrder", mock.Anything, mock.Anything).Return((*entities.ServiceOrder)(nil), usecase.ErrVehicleNotFound).Once()
		req, _ := http.NewRequest("POST", "/os", bytes.NewBufferString(`{"customer_id":"1","vehicle_id":"1"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("CreateServiceOrder - invalid id -> 400", func(t *testing.T) {
		mockUC.On("CreateServiceOrder", mock.Anything, mock.Anything).Return((*entities.ServiceOrder)(nil), usecase.ErrInvalidID).Once()
		req, _ := http.NewRequest("POST", "/os", bytes.NewBufferString(`{"customer_id":"1","vehicle_id":"1"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DiagnosisServiceOrder - insufficient parts -> 400", func(t *testing.T) {
		mockUC.On("DiagnosisServiceOrder", mock.Anything, mock.Anything).Return((*entities.ServiceOrder)(nil), usecase.ErrInsufficientPartsSupply).Once()
		req, _ := http.NewRequest("POST", "/os/1/diagnosis", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Approve estimate - invalid flow -> 400", func(t *testing.T) {
		mockUC.On("EstimateServiceOrder", mock.Anything, "1", mock.Anything).Return((*entities.ServiceOrder)(nil), usecase.ErrInvalidFlow).Once()
		req, _ := http.NewRequest("POST", "/os/1/estimate/approve", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Execution start - invalid status -> 400", func(t *testing.T) {
		mockUC.On("ExecutionServiceOrder", mock.Anything, "1", mock.Anything).Return((*entities.ServiceOrder)(nil), usecase.ErrInvalidStatus).Once()
		req, _ := http.NewRequest("POST", "/os/1/execution/create", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Payment - service order not found -> 404", func(t *testing.T) {
		mockUC.On("PaymentServiceOrder", mock.Anything, "1").Return((*entities.Payment)(nil), usecase.ErrServiceOrderNotFound).Once()
		req, _ := http.NewRequest("POST", "/os/1/payment", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Delivery - internal error -> 500", func(t *testing.T) {
		mockUC.On("DeliveryServiceOrder", mock.Anything, "1").Return((*entities.ServiceOrder)(nil), assert.AnError).Once()
		req, _ := http.NewRequest("POST", "/os/1/delivery", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	mockUC.AssertExpectations(t)
}
