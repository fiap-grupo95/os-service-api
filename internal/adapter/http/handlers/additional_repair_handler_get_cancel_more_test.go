package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	handler "github.com/fiap-grupo95/os-service-api/internal/adapter/http/handlers"
	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/handlers/mocks"
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
	"github.com/fiap-grupo95/os-service-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAdditionalRepairHandler_GetAndCancel_More(t *testing.T) {
	mockUC := mocks.NewMockIAdditionalRepairUseCase()
	h := handler.NewAdditionalRepairHandler(mockUC)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/additional-repair/:id", h.GetAdditionalRepair)
	r.POST("/additional-repair/:id/cancel", h.CancelAdditionalRepair)

	t.Run("get not found -> 404", func(t *testing.T) {
		mockUC.On("GetAdditionalRepair", mock.Anything, "1").Return((*entities.AdditionalRepair)(nil), usecase.ErrAdditionalRepairNotFound).Once()
		req, _ := http.NewRequest(http.MethodGet, "/additional-repair/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("get internal error -> 500", func(t *testing.T) {
		mockUC.On("GetAdditionalRepair", mock.Anything, "2").Return((*entities.AdditionalRepair)(nil), errors.New("fail")).Once()
		req, _ := http.NewRequest(http.MethodGet, "/additional-repair/2", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("cancel internal error -> 500", func(t *testing.T) {
		mockUC.On("CancelAdditionalRepair", mock.Anything, "3").Return((*entities.AdditionalRepair)(nil), errors.New("fail")).Once()
		req, _ := http.NewRequest(http.MethodPost, "/additional-repair/3/cancel", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("cancel success", func(t *testing.T) {
		ar := &entities.AdditionalRepair{ID: "4", ServiceOrderID: "10", Status: valueobject.StatusARCancelada}
		mockUC.On("CancelAdditionalRepair", mock.Anything, "4").Return(ar, nil).Once()
		req, _ := http.NewRequest(http.MethodPost, "/additional-repair/4/cancel", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	mockUC.AssertExpectations(t)
}
