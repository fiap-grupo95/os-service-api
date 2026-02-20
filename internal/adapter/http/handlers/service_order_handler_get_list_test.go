package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceOrderHandler_GetAndList(t *testing.T) {
	mockUC, h, r := setupServiceOrderHandlerTest(t)
	gin.SetMode(gin.TestMode)

	r.GET("/os/:id", h.GetServiceOrder)
	r.GET("/os", h.ListServiceOrders)

	t.Run("get success", func(t *testing.T) {
		mockUC.On("GetServiceOrder", mock.Anything, mock.MatchedBy(func(so entities.ServiceOrder) bool { return so.ID == "1" }), false).Return(&entities.ServiceOrder{ID: "1"}, nil).Once()
		req, _ := http.NewRequest(http.MethodGet, "/os/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("get not found -> 404", func(t *testing.T) {
		mockUC.On("GetServiceOrder", mock.Anything, mock.MatchedBy(func(so entities.ServiceOrder) bool { return so.ID == "2" }), false).Return((*entities.ServiceOrder)(nil), usecase.ErrServiceOrderNotFound).Once()
		req, _ := http.NewRequest(http.MethodGet, "/os/2", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("get internal error -> 500", func(t *testing.T) {
		mockUC.On("GetServiceOrder", mock.Anything, mock.MatchedBy(func(so entities.ServiceOrder) bool { return so.ID == "3" }), false).Return((*entities.ServiceOrder)(nil), assert.AnError).Once()
		req, _ := http.NewRequest(http.MethodGet, "/os/3", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("list success", func(t *testing.T) {
		mockUC.On("ListServiceOrders", mock.Anything).Return([]*entities.ServiceOrder{{ID: "1"}, {ID: "2"}}, nil).Once()
		req, _ := http.NewRequest(http.MethodGet, "/os", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("list internal error -> 500", func(t *testing.T) {
		mockUC.On("ListServiceOrders", mock.Anything).Return(([]*entities.ServiceOrder)(nil), assert.AnError).Once()
		req, _ := http.NewRequest(http.MethodGet, "/os", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	mockUC.AssertExpectations(t)
}
