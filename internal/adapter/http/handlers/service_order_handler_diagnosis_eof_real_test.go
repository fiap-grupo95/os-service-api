package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceOrderHandler_Diagnosis_EmptyStringBodyEOF(t *testing.T) {
	mockUC, h, r := setupServiceOrderHandlerTest(t)
	gin.SetMode(gin.TestMode)

	r.POST("/os/:id/diagnosis", h.DiagnosisServiceOrder)

	mockUC.On("DiagnosisServiceOrder", mock.Anything, mock.Anything).Return(&entities.ServiceOrder{ID: "1"}, nil).Once()

	req, _ := http.NewRequest(http.MethodPost, "/os/1/diagnosis", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// If Gin decodes as EOF, handler should accept empty request and still call usecase.
	// If Gin decodes as invalid json, it will be 400. We assert it's not 500.
	assert.NotEqual(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}
