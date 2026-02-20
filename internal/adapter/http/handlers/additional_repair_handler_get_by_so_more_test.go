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
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAdditionalRepairHandler_GetBySO_EmptyAndError(t *testing.T) {
	mockUC := mocks.NewMockIAdditionalRepairUseCase()
	h := handler.NewAdditionalRepairHandler(mockUC)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/additional-repair/service-orders/:id", h.GetAdditionalRepairBySO)

	t.Run("empty list -> 200", func(t *testing.T) {
		mockUC.On("GetAdditionalRepairBySO", mock.Anything, "1").Return([]entities.AdditionalRepair{}, nil).Once()
		req, _ := http.NewRequest(http.MethodGet, "/additional-repair/service-orders/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("non-empty list -> 200", func(t *testing.T) {
		mockUC.On("GetAdditionalRepairBySO", mock.Anything, "2").Return([]entities.AdditionalRepair{{ID: "ar-1", ServiceOrderID: "2", Status: valueobject.StatusARAberta}}, nil).Once()
		req, _ := http.NewRequest(http.MethodGet, "/additional-repair/service-orders/2", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("error -> 500", func(t *testing.T) {
		mockUC.On("GetAdditionalRepairBySO", mock.Anything, "3").Return(([]entities.AdditionalRepair)(nil), errors.New("fail")).Once()
		req, _ := http.NewRequest(http.MethodGet, "/additional-repair/service-orders/3", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	mockUC.AssertExpectations(t)
}
