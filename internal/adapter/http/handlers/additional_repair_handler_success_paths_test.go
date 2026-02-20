package handlers_test

import (
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

func TestAdditionalRepairHandler_SuccessPaths(t *testing.T) {
	mockUC := mocks.NewMockIAdditionalRepairUseCase()
	h := handler.NewAdditionalRepairHandler(mockUC)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/additional-repair/service-orders/:id", h.GetAdditionalRepairBySO)
	r.POST("/additional-repair/:id/approve", h.ApproveAdditionalRepair)
	r.POST("/additional-repair/:id/reject", h.RejectAdditionalRepair)

	t.Run("get by service order success", func(t *testing.T) {
		mockUC.On("GetAdditionalRepairBySO", mock.Anything, "10").Return([]entities.AdditionalRepair{{ID: "ar-1", ServiceOrderID: "10", Status: valueobject.StatusARAberta}}, nil).Once()
		req, _ := http.NewRequest(http.MethodGet, "/additional-repair/service-orders/10", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("approve success", func(t *testing.T) {
		mockUC.On("CustomerApprovalStatus", mock.Anything, "1", "APPROVED").Return(&entities.AdditionalRepair{ID: "1", Status: valueobject.StatusAAprovada}, nil).Once()
		req, _ := http.NewRequest(http.MethodPost, "/additional-repair/1/approve", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("reject success", func(t *testing.T) {
		mockUC.On("CustomerApprovalStatus", mock.Anything, "2", "REJECTED").Return(&entities.AdditionalRepair{ID: "2", Status: valueobject.StatusARRejeitada}, nil).Once()
		req, _ := http.NewRequest(http.MethodPost, "/additional-repair/2/reject", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	mockUC.AssertExpectations(t)
}
