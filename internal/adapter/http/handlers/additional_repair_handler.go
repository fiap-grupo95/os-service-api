package handlers

import (
	"errors"
	"net/http"

	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/request"
	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/logs"
	"github.com/fiap-grupo95/os-service-api/internal/usecase"
	pkg "github.com/fiap-grupo95/os-service-api/pkg/utils/errors"

	"github.com/gin-gonic/gin"
)

var (
	errInvalidAdditionalRepairID    = pkg.NewDomainErrorSimple("INVALID_ADDITIONAL_REPAIR_ID", "Invalid additional repair ID", http.StatusBadRequest)
	errInvalidAdditionalRepairInput = pkg.NewDomainErrorSimple("INVALID_ADDITIONAL_REPAIR_INPUT", "Invalid additional repair payload", http.StatusBadRequest)
)

const (
	APPROVED_FLOW = "APPROVED"
	REJECTED_FLOW = "REJECTED"
)

type IAdditionalRepairHandler interface {
	CreateAdditionalRepair(c *gin.Context)
	GetAdditionalRepair(c *gin.Context)
	GetAdditionalRepairBySO(c *gin.Context)
	ApproveAdditionalRepair(c *gin.Context)
	RejectAdditionalRepair(c *gin.Context)
	CancelAdditionalRepair(c *gin.Context)
}

// AdditionalRepairHandler handles HTTP requests for additional repairs
// @title Additional Repair API
// @version 1.0
// @description API for managing additional repairs in the workshop management system
type AdditionalRepairHandler struct {
	additionalRepairUseCase usecase.IAdditionalRepairUseCase
}

func NewAdditionalRepairHandler(useCase usecase.IAdditionalRepairUseCase) *AdditionalRepairHandler {
	return &AdditionalRepairHandler{
		additionalRepairUseCase: useCase,
	}
}

// GetAdditionalRepair godoc
// @Summary Get additional repair by ID
// @Description Retrieve an additional repair by its ID
// @Tags Additional Repairs
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "Additional Repair ID"
// @Success 200 {object} response.AdditionalRepairResponse
// @Failure 400 {object} pkg.AppError
// @Failure 404 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /v1/additional-repair/{id} [get]
func (h *AdditionalRepairHandler) GetAdditionalRepair(c *gin.Context) {
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "AdditionalRepair/Get")

	additionalRepairID := c.Param("id")
	if additionalRepairID == "" {
		logger.Error().Str("id", c.Param("id")).Msg("Failed to parse additional repair ID")
		c.JSON(errInvalidAdditionalRepairID.HTTPStatus, errInvalidAdditionalRepairID.ToHTTPError())
		return
	}

	foundAdr, err := h.additionalRepairUseCase.GetAdditionalRepair(ctx, additionalRepairID)
	if err != nil {
		logger.Error().Err(err).Str("ADR_ID", additionalRepairID).Msg("Failed to retrieve additional repair")
		appErr := mapAdditionalRepairError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, toAdditionalRepairResponse(foundAdr))
}

// GetAdditionalRepairBySO godoc
// @Summary Get additional repairs by service order ID
// @Description Retrieve additional repairs by service order ID
// @Tags Additional Repairs
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "Service Order ID"
// @Success 200 {array} response.AdditionalRepairResponse
// @Failure 400 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /v1/additional-repair/service-orders/{id} [get]
func (h *AdditionalRepairHandler) GetAdditionalRepairBySO(c *gin.Context) {
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "AdditionalRepair/GetAdditionalRepairBySO")

	serviceOrderID := c.Param("id")
	if serviceOrderID == "" {
		logger.Error().Str("id", c.Param("id")).Msg("Failed to parse service order ID")
		c.JSON(errInvalidAdditionalRepairID.HTTPStatus, errInvalidAdditionalRepairID.ToHTTPError())
		return
	}

	adrs, err := h.additionalRepairUseCase.GetAdditionalRepairBySO(ctx, serviceOrderID)
	if err != nil {
		logger.Error().Err(err).Str("SO_ID", serviceOrderID).Msg("Failed to retrieve additional repair")
		appErr := mapAdditionalRepairError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, toAdditionalRepairResponseList(adrs))
}

// CreateAdditionalRepair godoc
// @Summary Create a new additional repair
// @Description Create a new additional repair record
// @Tags Additional Repairs
// @Security Bearer
// @Accept json
// @Produce json
// @Param repair body request.AdditionalRepairCreateRequest true "Additional Repair Information"
// @Success 201 {object} response.AdditionalRepairResponse
// @Failure 400 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /v1/additional-repair [post]
func (h *AdditionalRepairHandler) CreateAdditionalRepair(c *gin.Context) {
	var payload request.AdditionalRepairCreateRequest
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "AdditionalRepair/Create")

	if err := c.ShouldBindJSON(&payload); err != nil {
		logger.Error().Err(err).Msg("Failed to bind JSON for create additional repair")
		c.JSON(errInvalidAdditionalRepairInput.HTTPStatus, errInvalidAdditionalRepairInput.ToHTTPError())
		return
	}

	request := toAdditionalRepairEntityFromCreate(payload)
	created, err := h.additionalRepairUseCase.CreateAdditionalRepair(ctx, request)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create additional repair")
		appErr := mapAdditionalRepairError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusCreated, toAdditionalRepairResponse(created))
}

// ApproveAdditionalRepair godoc
// @Summary Approve an additional repair
// @Description Approve an additional repair by its ID
// @Tags Additional Repairs
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "Additional Repair ID"
// @Success 200 {object} response.AdditionalRepairResponse
// @Failure 400 {object} pkg.AppError
// @Failure 404 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /v1/additional-repair/{id}/approve [post]
func (h *AdditionalRepairHandler) ApproveAdditionalRepair(c *gin.Context) {
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "AdditionalRepair/ApproveAdditionalRepair")
	additionalRepairID := c.Param("id")
	if additionalRepairID == "" {
		logger.Error().Str("id", c.Param("id")).Msg("Failed to parse additional repair ID")
		c.JSON(errInvalidAdditionalRepairID.HTTPStatus, errInvalidAdditionalRepairID.ToHTTPError())
		return
	}

	ar, err := h.additionalRepairUseCase.CustomerApprovalStatus(ctx, additionalRepairID, APPROVED_FLOW)
	if err != nil {
		logger.Error().Err(err).Str("ADR_ID", additionalRepairID).Msg("Failed to approve additional repair")
		appErr := mapAdditionalRepairError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, toAdditionalRepairResponse(ar))
}

// RejectAdditionalRepair godoc
// @Summary Reject an additional repair
// @Description Reject an additional repair by its ID
// @Tags Additional Repairs
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "Additional Repair ID"
// @Success 200 {object} response.AdditionalRepairResponse
// @Failure 400 {object} pkg.AppError
// @Failure 404 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /v1/additional-repair/{id}/reject [post]
func (h *AdditionalRepairHandler) RejectAdditionalRepair(c *gin.Context) {
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "AdditionalRepair/RejectAdditionalRepair")
	additionalRepairID := c.Param("id")
	if additionalRepairID == "" {
		logger.Error().Str("id", c.Param("id")).Msg("Failed to parse additional repair ID")
		c.JSON(errInvalidAdditionalRepairID.HTTPStatus, errInvalidAdditionalRepairID.ToHTTPError())
		return
	}

	ar, err := h.additionalRepairUseCase.CustomerApprovalStatus(ctx, additionalRepairID, REJECTED_FLOW)
	if err != nil {
		logger.Error().Err(err).Str("ADR_ID", additionalRepairID).Msg("Failed to reject additional repair")
		appErr := mapAdditionalRepairError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, toAdditionalRepairResponse(ar))
}

// CancelAdditionalRepair godoc
// @Summary Cancel an additional repair
// @Description Cancel an additional repair by its ID
// @Tags Additional Repairs
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "Additional Repair ID"
// @Success 200 {object} response.AdditionalRepairResponse
// @Failure 400 {object} pkg.AppError
// @Failure 404 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /v1/additional-repair/{id}/cancel [post]
func (h *AdditionalRepairHandler) CancelAdditionalRepair(c *gin.Context) {
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "AdditionalRepair/CancelAdditionalRepair")
	additionalRepairID := c.Param("id")
	if additionalRepairID == "" {
		logger.Error().Str("id", c.Param("id")).Msg("Failed to parse additional repair ID")
		c.JSON(errInvalidAdditionalRepairID.HTTPStatus, errInvalidAdditionalRepairID.ToHTTPError())
		return
	}

	ar, err := h.additionalRepairUseCase.CancelAdditionalRepair(ctx, additionalRepairID)
	if err != nil {
		logger.Error().Err(err).Str("ADR_ID", additionalRepairID).Msg("Failed to cancel additional repair")
		appErr := mapAdditionalRepairError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, toAdditionalRepairResponse(ar))
}

func mapAdditionalRepairError(err error) *pkg.AppError {
	switch {
	case errors.Is(err, usecase.ErrAdditionalRepairNotFound):
		return pkg.NewDomainErrorSimple("ADDITIONAL_REPAIR_NOT_FOUND", "Additional repair not found", http.StatusNotFound)
	case errors.Is(err, usecase.ErrStatusNotPermitted):
		return pkg.NewDomainErrorSimple("STATUS_NOT_PERMITTED", "Additional repair status not permitted", http.StatusBadRequest)
	case errors.Is(err, usecase.ErrServiceNotFound):
		return pkg.NewDomainErrorSimple("SERVICE_NOT_FOUND", "Service not found", http.StatusNotFound)
	case errors.Is(err, usecase.ErrPartsSupplyNotFound):
		return pkg.NewDomainErrorSimple("PARTS_SUPPLY_NOT_FOUND", "Parts supply not found", http.StatusNotFound)
	default:
		return pkg.NewDomainError("INTERNAL_ERROR", "An internal error occurred", err, http.StatusInternalServerError)
	}
}

func toAdditionalRepairEntityFromCreate(payload request.AdditionalRepairCreateRequest) entities.AdditionalRepair {
	return entities.AdditionalRepair{
		ServiceOrderID: payload.ServiceOrderID,
		Description:    payload.Description,
		Services:       mapAdditionalRepairServices(payload.Services),
		PartsSupplies:  mapAdditionalRepairPartsSupplies(payload.PartsSupplies),
	}
}

func toAdditionalRepairEntityFromItems(payload request.AdditionalRepairItemsRequest) entities.AdditionalRepair {
	return entities.AdditionalRepair{
		ServiceOrderID: payload.ServiceOrderID,
		Description:    payload.Description,
		Services:       mapAdditionalRepairServices(payload.Services),
		PartsSupplies:  mapAdditionalRepairPartsSupplies(payload.PartsSupplies),
	}
}

func mapAdditionalRepairServices(services []request.AdditionalRepairServiceItem) []entities.Service {
	if len(services) == 0 {
		return nil
	}

	result := make([]entities.Service, 0, len(services))
	for _, service := range services {
		result = append(result, entities.Service{ID: service.ID})
	}
	return result
}

func mapAdditionalRepairPartsSupplies(partsSupplies []request.AdditionalRepairPartsSupplyItem) []entities.PartsSupply {
	if len(partsSupplies) == 0 {
		return nil
	}

	result := make([]entities.PartsSupply, 0, len(partsSupplies))
	for _, ps := range partsSupplies {
		result = append(result, entities.PartsSupply{
			ID:       ps.ID,
			Quantity: ps.Quantity,
		})
	}
	return result
}

func toAdditionalRepairResponse(entity *entities.AdditionalRepair) response.AdditionalRepairResponse {
	if entity == nil {
		return response.AdditionalRepairResponse{}
	}
	res := response.AdditionalRepairResponse{
		ID:             entity.ID,
		ServiceOrderID: entity.ServiceOrderID,
		Description:    entity.Description,
		Status:         entity.Status.String(),
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
		Services:       make([]response.AdditionalRepairServiceResponse, 0, len(entity.Services)),
		PartsSupplies:  make([]response.AdditionalRepairPartsSupplyResponse, 0, len(entity.PartsSupplies)),
	}

	if entity.Estimate != nil {
		res.Estimate = &response.EstimateResponse{
			ID:                 entity.Estimate.ID,
			Value:              entity.Estimate.Value,
			Status:             entity.Estimate.Status,
			ServiceOrderID:     entity.Estimate.ServiceOrderID,
			AdditionalRepairID: entity.Estimate.AdditionalRepairID,
		}
	}

	for _, service := range entity.Services {
		res.Services = append(res.Services, response.AdditionalRepairServiceResponse{
			ID:    service.ID,
			Name:  service.Name,
			Price: service.Price,
		})
	}

	for _, ps := range entity.PartsSupplies {
		res.PartsSupplies = append(res.PartsSupplies, response.AdditionalRepairPartsSupplyResponse{
			ID:       ps.ID,
			Price:    ps.Price,
			Quantity: ps.Quantity,
		})
	}

	return res
}

func toAdditionalRepairResponseList(entities []entities.AdditionalRepair) []response.AdditionalRepairResponse {
	list := make([]response.AdditionalRepairResponse, 0, len(entities))
	for _, entity := range entities {
		list = append(list, toAdditionalRepairResponse(&entity))
	}
	return list
}
