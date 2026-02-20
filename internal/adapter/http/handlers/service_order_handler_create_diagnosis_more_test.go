package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/request"
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceOrderHandler_CreateDiagnosisAndGetFullData(t *testing.T) {
	mockUC, h, r := setupServiceOrderHandlerTest(t)
	gin.SetMode(gin.TestMode)

	r.POST("/os", h.CreateServiceOrder)
	r.POST("/os/:id/diagnosis", h.DiagnosisServiceOrder)
	r.GET("/os/:id", h.GetServiceOrder)

	t.Run("create success -> 201", func(t *testing.T) {
		payload := request.ServiceOrderCreateRequest{CustomerID: "c1", VehicleID: "v1"}
		body, _ := json.Marshal(payload)

		mockUC.On("CreateServiceOrder", mock.Anything, mock.Anything).Return(&entities.ServiceOrder{ID: "1"}, nil).Once()
		req, _ := http.NewRequest(http.MethodPost, "/os", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("diagnosis success -> 200", func(t *testing.T) {
		mockUC.On("DiagnosisServiceOrder", mock.Anything, mock.Anything).Return(&entities.ServiceOrder{ID: "2"}, nil).Once()
		req, _ := http.NewRequest(http.MethodPost, "/os/2/diagnosis", bytes.NewBufferString(`{"services":[],"parts_supplies":[]}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("diagnosis not found -> 404", func(t *testing.T) {
		mockUC.On("DiagnosisServiceOrder", mock.Anything, mock.Anything).Return((*entities.ServiceOrder)(nil), usecase.ErrServiceOrderNotFound).Once()
		req, _ := http.NewRequest(http.MethodPost, "/os/3/diagnosis", bytes.NewBufferString(`{"services":[],"parts_supplies":[]}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("get full_data=true -> 200", func(t *testing.T) {
		mockUC.On("GetServiceOrder", mock.Anything, mock.MatchedBy(func(so entities.ServiceOrder) bool { return so.ID == "4" }), true).Return(&entities.ServiceOrder{ID: "4"}, nil).Once()
		req, _ := http.NewRequest(http.MethodGet, "/os/4?full_data=true", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	mockUC.AssertExpectations(t)
}
