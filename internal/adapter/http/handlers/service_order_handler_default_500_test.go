package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceOrderHandler_Default500Branches(t *testing.T) {
	mockUC, h, r := setupServiceOrderHandlerTest(t)
	gin.SetMode(gin.TestMode)

	r.GET("/os/:id", h.GetServiceOrder)
	r.GET("/os", h.ListServiceOrders)
	r.POST("/os", h.CreateServiceOrder)
	r.POST("/os/:id/diagnosis", h.DiagnosisServiceOrder)

	t.Run("get default 500", func(t *testing.T) {
		mockUC.On("GetServiceOrder", mock.Anything, mock.MatchedBy(func(so entities.ServiceOrder) bool { return so.ID == "1" }), false).Return((*entities.ServiceOrder)(nil), assert.AnError).Once()
		req, _ := http.NewRequest(http.MethodGet, "/os/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("list default 500", func(t *testing.T) {
		mockUC.On("ListServiceOrders", mock.Anything).Return(([]*entities.ServiceOrder)(nil), assert.AnError).Once()
		req, _ := http.NewRequest(http.MethodGet, "/os", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("create default 500", func(t *testing.T) {
		mockUC.On("CreateServiceOrder", mock.Anything, mock.Anything).Return((*entities.ServiceOrder)(nil), assert.AnError).Once()
		req, _ := http.NewRequest(http.MethodPost, "/os", bytes.NewBufferString(`{"customer_id":"c1","vehicle_id":"v1"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("diagnosis default 500", func(t *testing.T) {
		mockUC.On("DiagnosisServiceOrder", mock.Anything, mock.Anything).Return((*entities.ServiceOrder)(nil), assert.AnError).Once()
		req, _ := http.NewRequest(http.MethodPost, "/os/2/diagnosis", bytes.NewBufferString(`{"services":[],"parts_supplies":[]}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	mockUC.AssertExpectations(t)
}
