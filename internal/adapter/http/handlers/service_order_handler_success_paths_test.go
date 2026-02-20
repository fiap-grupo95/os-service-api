package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceOrderHandler_SuccessPaths(t *testing.T) {
	mockUC, h, r := setupServiceOrderHandlerTest(t)
	gin.SetMode(gin.TestMode)

	r.POST("/os/:id/estimate/approve", h.ApproveServiceOrderEstimate)
	r.POST("/os/:id/estimate/reject", h.RejectServiceOrderEstimate)
	r.POST("/os/:id/estimate/cancel", h.CancelServiceOrderEstimate)
	r.POST("/os/:id/execution/create", h.ExecutionServiceOrder)
	r.POST("/os/:id/execution/finish", h.FinishServiceOrderExecution)
	r.POST("/os/:id/payment", h.PaymentServiceOrder)
	r.POST("/os/:id/delivery", h.DeliveryServiceOrder)

	t.Run("approve estimate success", func(t *testing.T) {
		mockUC.On("EstimateServiceOrder", mock.Anything, "1", mock.Anything).Return(&entities.ServiceOrder{ID: "1"}, nil).Once()
		req, _ := http.NewRequest("POST", "/os/1/estimate/approve", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("reject estimate success", func(t *testing.T) {
		mockUC.On("EstimateServiceOrder", mock.Anything, "2", mock.Anything).Return(&entities.ServiceOrder{ID: "2"}, nil).Once()
		req, _ := http.NewRequest("POST", "/os/2/estimate/reject", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("cancel estimate success", func(t *testing.T) {
		mockUC.On("EstimateServiceOrder", mock.Anything, "3", mock.Anything).Return(&entities.ServiceOrder{ID: "3"}, nil).Once()
		req, _ := http.NewRequest("POST", "/os/3/estimate/cancel", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("execution start success", func(t *testing.T) {
		mockUC.On("ExecutionServiceOrder", mock.Anything, "4", mock.Anything).Return(&entities.ServiceOrder{ID: "4"}, nil).Once()
		req, _ := http.NewRequest("POST", "/os/4/execution/create", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("execution finish success", func(t *testing.T) {
		mockUC.On("ExecutionServiceOrder", mock.Anything, "5", mock.Anything).Return(&entities.ServiceOrder{ID: "5"}, nil).Once()
		req, _ := http.NewRequest("POST", "/os/5/execution/finish", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("payment success", func(t *testing.T) {
		mockUC.On("PaymentServiceOrder", mock.Anything, "6").Return(&entities.Payment{ID: "pay-1"}, nil).Once()
		req, _ := http.NewRequest("POST", "/os/6/payment", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("delivery success", func(t *testing.T) {
		mockUC.On("DeliveryServiceOrder", mock.Anything, "7").Return(&entities.ServiceOrder{ID: "7"}, nil).Once()
		req, _ := http.NewRequest("POST", "/os/7/delivery", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	mockUC.AssertExpectations(t)
}
