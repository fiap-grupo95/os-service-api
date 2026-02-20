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

func TestServiceOrderHandler_Diagnosis_EmptyBodyEOF(t *testing.T) {
	mockUC, h, r := setupServiceOrderHandlerTest(t)
	gin.SetMode(gin.TestMode)

	r.POST("/os/:id/diagnosis", h.DiagnosisServiceOrder)

	mockUC.On("DiagnosisServiceOrder", mock.Anything, mock.Anything).Return(&entities.ServiceOrder{ID: "1"}, nil).Once()

	req, _ := http.NewRequest(http.MethodPost, "/os/1/diagnosis", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}
