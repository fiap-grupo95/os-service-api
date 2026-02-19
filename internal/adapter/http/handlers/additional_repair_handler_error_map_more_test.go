package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	handler "github.com/fiap-grupo95/os-service-api/internal/adapter/http/handlers"
	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/handlers/mocks"
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAdditionalRepairHandler_ErrorMapping_MoreCases(t *testing.T) {
	mockUC := mocks.NewMockIAdditionalRepairUseCase()
	h := handler.NewAdditionalRepairHandler(mockUC)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/additional-repair/:id/approve", h.ApproveAdditionalRepair)
	r.POST("/additional-repair/:id/reject", h.RejectAdditionalRepair)
	r.POST("/additional-repair/:id/cancel", h.CancelAdditionalRepair)

	t.Run("approve - not found -> 404", func(t *testing.T) {
		mockUC.On("CustomerApprovalStatus", mock.Anything, "1", "APPROVED").Return((*entities.AdditionalRepair)(nil), usecase.ErrAdditionalRepairNotFound).Once()
		req, _ := http.NewRequest(http.MethodPost, "/additional-repair/1/approve", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("reject - not found -> 404", func(t *testing.T) {
		mockUC.On("CustomerApprovalStatus", mock.Anything, "2", "REJECTED").Return((*entities.AdditionalRepair)(nil), usecase.ErrAdditionalRepairNotFound).Once()
		req, _ := http.NewRequest(http.MethodPost, "/additional-repair/2/reject", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("cancel - not found -> 404", func(t *testing.T) {
		mockUC.On("CancelAdditionalRepair", mock.Anything, "3").Return((*entities.AdditionalRepair)(nil), usecase.ErrAdditionalRepairNotFound).Once()
		req, _ := http.NewRequest(http.MethodPost, "/additional-repair/3/cancel", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	mockUC.AssertExpectations(t)
}
