package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fiap-grupo95/os-service-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestServiceOrderHandler_ParseIsFullDataParam_InvalidValue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/os/:id", func(c *gin.Context) {
		assert.False(t, parseIsFullDataParam(c))
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodGet, "/os/1?full_data=TRUE", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestServiceOrderHandler_IsBadRequestServiceOrderError(t *testing.T) {
	assert.True(t, isBadRequestServiceOrderError(usecase.ErrInvalidTransitionStatusToDiagnosis))
	assert.True(t, isBadRequestServiceOrderError(usecase.ErrInvalidTransitionStatusToEstimate))
	assert.True(t, isBadRequestServiceOrderError(usecase.ErrInvalidTransitionStatusToExecution))
	assert.True(t, isBadRequestServiceOrderError(usecase.ErrInvalidTransitionStatusToDelivery))
	assert.True(t, isBadRequestServiceOrderError(usecase.ErrInvalidStatus))
	assert.True(t, isBadRequestServiceOrderError(usecase.ErrInvalidFlow))
	assert.True(t, isBadRequestServiceOrderError(usecase.ErrInsufficientPartsSupply))
	assert.False(t, isBadRequestServiceOrderError(errors.New("x")))
}

func TestServiceOrderHandler_GetStatusError(t *testing.T) {
	assert.Equal(t, http.StatusNotFound, getStatusError(usecase.ErrServiceOrderNotFound))
	assert.Equal(t, http.StatusNotFound, getStatusError(usecase.ErrVehicleNotFound))
	assert.Equal(t, http.StatusNotFound, getStatusError(usecase.ErrCustomerNotFound))
	assert.Equal(t, http.StatusNotFound, getStatusError(usecase.ErrServiceNotFound))
	assert.Equal(t, http.StatusNotFound, getStatusError(usecase.ErrPartsSupplyNotFound))
	assert.Equal(t, http.StatusBadRequest, getStatusError(usecase.ErrInvalidID))
	assert.Equal(t, http.StatusBadRequest, getStatusError(usecase.ErrInvalidFlow))
	assert.Equal(t, http.StatusInternalServerError, getStatusError(errors.New("x")))
}
