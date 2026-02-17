package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	request "github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/request"
	response "github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/response"
	handler "github.com/fiap-grupo95/os-service-api/internal/adapter/http/handlers"
	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/handlers/mocks"
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupADRRouter(h *handler.AdditionalRepairHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/additional-repair", h.CreateAdditionalRepair)
	r.GET("/additional-repair/:id", h.GetAdditionalRepair)
	r.POST("/additional-repair/approve/:id", h.ApproveAdditionalRepair)
	r.POST("/additional-repair/reject/:id", h.RejectAdditionalRepair)

	return r
}

func TestCreateAdditionalRepair_Success(t *testing.T) {
	mockUC := mocks.NewMockIAdditionalRepairUseCase()
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	payload := request.AdditionalRepairCreateRequest{ServiceOrderID: "1", Description: "desc"}
	expected := entities.AdditionalRepair{ServiceOrderID: "1", Description: "desc"}
	created := &entities.AdditionalRepair{
		ID:             "10",
		ServiceOrderID: expected.ServiceOrderID,
		Description:    expected.Description,
		Status:         valueobject.StatusARAberta,
	}
	mockUC.On("CreateAdditionalRepair", mock.Anything, expected).Return(created, nil)

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/additional-repair", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp response.AdditionalRepairResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, resp.ID)
	assert.Equal(t, created.ServiceOrderID, resp.ServiceOrderID)
	assert.Equal(t, created.Description, resp.Description)
	assert.Equal(t, created.Status.String(), resp.Status)
	mockUC.AssertExpectations(t)
}

func TestCreateAdditionalRepair_Error(t *testing.T) {
	mockUC := mocks.NewMockIAdditionalRepairUseCase()
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	payload := request.AdditionalRepairCreateRequest{ServiceOrderID: "1", Description: "desc"}
	expected := entities.AdditionalRepair{ServiceOrderID: "1", Description: "desc"}
	mockUC.On("CreateAdditionalRepair", mock.Anything, expected).Return((*entities.AdditionalRepair)(nil), errors.New("fail"))

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/additional-repair", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetAdditionalRepair_Success(t *testing.T) {
	mockUC := mocks.NewMockIAdditionalRepairUseCase()
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	now := time.Now()
	expected := &entities.AdditionalRepair{
		ID:             "1",
		ServiceOrderID: "10",
		Description:    "desc",
		Status:         valueobject.StatusAAprovada,
		Estimate:       &entities.Estimate{ID: "e1", ServiceOrderID: "10", Value: 123.45, Status: "APPROVED"},
		CreatedAt:      now,
		UpdatedAt:      now,
		Services: []entities.Service{
			{ID: "2", Name: "service", Price: 80},
		},
		PartsSupplies: []entities.PartsSupply{
			{ID: "3", Price: 43.45, Quantity: 2},
		},
	}
	mockUC.On("GetAdditionalRepair", mock.Anything, "1").Return(expected, nil)

	req, _ := http.NewRequest("GET", "/additional-repair/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp response.AdditionalRepairResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, resp.ID)
	assert.Equal(t, expected.ServiceOrderID, resp.ServiceOrderID)
	assert.Equal(t, expected.Description, resp.Description)
	assert.Equal(t, expected.Status.String(), resp.Status)
	assert.Equal(t, expected.Estimate.Value, resp.Estimate)
	assert.Len(t, resp.Services, 1)
	assert.Equal(t, expected.Services[0].ID, resp.Services[0].ID)
	assert.Len(t, resp.PartsSupplies, 1)
	assert.Equal(t, expected.PartsSupplies[0].ID, resp.PartsSupplies[0].ID)
	mockUC.AssertExpectations(t)
}

func TestGetAdditionalRepair_Error(t *testing.T) {
	mockUC := mocks.NewMockIAdditionalRepairUseCase()
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	// Return empty struct and error
	mockUC.On("GetAdditionalRepair", mock.Anything, "999").Return((*entities.AdditionalRepair)(nil), errors.New("not found"))

	req, _ := http.NewRequest("GET", "/additional-repair/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}

func TestApproveAdditionalRepair_Success(t *testing.T) {
	mockUC := mocks.NewMockIAdditionalRepairUseCase()
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	ar := &entities.AdditionalRepair{ID: "1", ServiceOrderID: "10", Description: "desc", Status: valueobject.StatusAAprovada}
	mockUC.On("CustomerApprovalStatus", mock.Anything, "1", "APPROVED").Return(ar, nil)

	req, _ := http.NewRequest("POST", "/additional-repair/approve/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp response.AdditionalRepairResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, ar.ID, resp.ID)
	assert.Equal(t, ar.Status.String(), resp.Status)
	mockUC.AssertExpectations(t)
}

func TestApproveAdditionalRepair_Error(t *testing.T) {
	mockUC := mocks.NewMockIAdditionalRepairUseCase()
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	mockUC.On("CustomerApprovalStatus", mock.Anything, "1", "APPROVED").Return((*entities.AdditionalRepair)(nil), errors.New("fail"))

	req, _ := http.NewRequest("POST", "/additional-repair/approve/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}

func TestRejectAdditionalRepair_Success(t *testing.T) {
	mockUC := mocks.NewMockIAdditionalRepairUseCase()
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	ar := &entities.AdditionalRepair{ID: "1", ServiceOrderID: "10", Description: "desc", Status: valueobject.StatusARRejeitada}
	mockUC.On("CustomerApprovalStatus", mock.Anything, "1", "REJECTED").Return(ar, nil)

	req, _ := http.NewRequest("POST", "/additional-repair/reject/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp response.AdditionalRepairResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, ar.ID, resp.ID)
	assert.Equal(t, ar.Status.String(), resp.Status)
	mockUC.AssertExpectations(t)
}

func TestRejectAdditionalRepair_Error(t *testing.T) {
	mockUC := mocks.NewMockIAdditionalRepairUseCase()
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	mockUC.On("CustomerApprovalStatus", mock.Anything, "1", "REJECTED").Return((*entities.AdditionalRepair)(nil), errors.New("fail"))

	req, _ := http.NewRequest("POST", "/additional-repair/reject/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}
