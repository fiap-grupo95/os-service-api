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

func setupADRRouterMore(h *handler.AdditionalRepairHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/additional-repair/:id", h.GetAdditionalRepair)
	r.GET("/additional-repair/service-orders/:id", h.GetAdditionalRepairBySO)
	r.POST("/additional-repair/:id/approve", h.ApproveAdditionalRepair)
	r.POST("/additional-repair/:id/reject", h.RejectAdditionalRepair)
	r.POST("/additional-repair/:id/cancel", h.CancelAdditionalRepair)
	return r
}

func TestAdditionalRepairHandler_ErrorMapping(t *testing.T) {
	mockUC := mocks.NewMockIAdditionalRepairUseCase()
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouterMore(h)

	t.Run("get - not found maps to 404", func(t *testing.T) {
		mockUC.On("GetAdditionalRepair", mock.Anything, "1").Return((*entities.AdditionalRepair)(nil), usecase.ErrAdditionalRepairNotFound).Once()
		req, _ := http.NewRequest(http.MethodGet, "/additional-repair/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("get - internal error maps to 500", func(t *testing.T) {
		mockUC.On("GetAdditionalRepair", mock.Anything, "2").Return((*entities.AdditionalRepair)(nil), errors.New("fail")).Once()
		req, _ := http.NewRequest(http.MethodGet, "/additional-repair/2", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("get by so - internal error maps to 500", func(t *testing.T) {
		mockUC.On("GetAdditionalRepairBySO", mock.Anything, "10").Return(([]entities.AdditionalRepair)(nil), errors.New("fail")).Once()
		req, _ := http.NewRequest(http.MethodGet, "/additional-repair/service-orders/10", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("approve - status not permitted maps to 400", func(t *testing.T) {
		mockUC.On("CustomerApprovalStatus", mock.Anything, "3", "APPROVED").Return((*entities.AdditionalRepair)(nil), usecase.ErrStatusNotPermitted).Once()
		req, _ := http.NewRequest(http.MethodPost, "/additional-repair/3/approve", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("reject - parts supply not found maps to 404", func(t *testing.T) {
		mockUC.On("CustomerApprovalStatus", mock.Anything, "4", "REJECTED").Return((*entities.AdditionalRepair)(nil), usecase.ErrPartsSupplyNotFound).Once()
		req, _ := http.NewRequest(http.MethodPost, "/additional-repair/4/reject", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("cancel - success", func(t *testing.T) {
		ar := &entities.AdditionalRepair{ID: "5", ServiceOrderID: "10", Status: valueobject.StatusARCancelada}
		mockUC.On("CancelAdditionalRepair", mock.Anything, "5").Return(ar, nil).Once()
		req, _ := http.NewRequest(http.MethodPost, "/additional-repair/5/cancel", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	mockUC.AssertExpectations(t)
}
