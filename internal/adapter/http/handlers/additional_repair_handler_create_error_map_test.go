package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	handler "github.com/fiap-grupo95/os-service-api/internal/adapter/http/handlers"
	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/handlers/mocks"
	"github.com/fiap-grupo95/os-service-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAdditionalRepairHandler_Create_ErrorMapping(t *testing.T) {
	mockUC := mocks.NewMockIAdditionalRepairUseCase()
	h := handler.NewAdditionalRepairHandler(mockUC)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/additional-repair", h.CreateAdditionalRepair)

	t.Run("service not found -> 404", func(t *testing.T) {
		mockUC.On("CreateAdditionalRepair", mock.Anything, mock.Anything).Return(nil, usecase.ErrServiceNotFound).Once()
		req, _ := http.NewRequest(http.MethodPost, "/additional-repair", bytes.NewBufferString(`{"service_order_id":"1","description":"x","services":[{"id":"s1"}]}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("parts supply not found -> 404", func(t *testing.T) {
		mockUC.On("CreateAdditionalRepair", mock.Anything, mock.Anything).Return(nil, usecase.ErrPartsSupplyNotFound).Once()
		req, _ := http.NewRequest(http.MethodPost, "/additional-repair", bytes.NewBufferString(`{"service_order_id":"1","description":"x","parts_supplies":[{"id":"p1","quantity":1}]}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("internal error -> 500", func(t *testing.T) {
		mockUC.On("CreateAdditionalRepair", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		req, _ := http.NewRequest(http.MethodPost, "/additional-repair", bytes.NewBufferString(`{"service_order_id":"1","description":"x"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	mockUC.AssertExpectations(t)
}
