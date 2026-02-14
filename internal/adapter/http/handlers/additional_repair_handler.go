package handlers

import (
	"errors"
	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/request"
	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/logs"
	"github.com/fiap-grupo95/os-service-api/internal/usecase"
	"github.com/fiap-grupo95/os-service-api/pkg/utils/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	errInvalidAdditionalRepairID    = pkg.NewDomainErrorSimple("INVALID_ADDITIONAL_REPAIR_ID", "Invalid additional repair ID", http.StatusBadRequest)
	errInvalidAdditionalRepairInput = pkg.NewDomainErrorSimple("INVALID_ADDITIONAL_REPAIR_INPUT", "Invalid additional repair payload", http.StatusBadRequest)
)

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
// @Param id path int true "Additional Repair ID"
// @Success 200 {object} response.AdditionalRepairResponse
// @Failure 400 {object} pkg.AppError
// @Failure 404 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /additional-repairs/{id} [get]
func (h *AdditionalRepairHandler) GetAdditionalRepair(c *gin.Context) {
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "AdditionalRepair/Get")

	additionalRepairID, ok := parseAdditionalRepairID(c)
	if !ok {
		return
	}

	foundAdr, err := h.additionalRepairUseCase.GetAdditionalRepair(ctx, additionalRepairID)
	if err != nil {
		logger.Error().Err(err).Uint("ADR_ID", additionalRepairID).Msg("Failed to retrieve additional repair")
		appErr := mapAdditionalRepairError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, toAdditionalRepairResponse(foundAdr))
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
// @Router /additional-repairs [post]
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

// AddPartSupplyAndService godoc
// @Summary Add parts supply and service to additional repair
// @Description Add parts supply and service to an existing additional repair
// @Tags Additional Repairs
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Additional Repair ID"
// @Param repair body request.AdditionalRepairItemsRequest true "Parts Supply and Service Information"
// @Success 201 {object} response.OperationMessageResponse
// @Failure 400 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /additional-repairs/{id}/add [post]
func (h *AdditionalRepairHandler) AddPartSupplyAndService(c *gin.Context) {
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "AdditionalRepair/AddItems")
	additionalRepairID, ok := parseAdditionalRepairID(c)
	if !ok {
		return
	}

	var payload request.AdditionalRepairItemsRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		logger.Error().Err(err).Msg("Failed to bind JSON for add parts supply and service")
		c.JSON(errInvalidAdditionalRepairInput.HTTPStatus, errInvalidAdditionalRepairInput.ToHTTPError())
		return
	}

	request := toAdditionalRepairEntityFromItems(payload)
	if err := h.additionalRepairUseCase.AddPartSupplyAndService(ctx, additionalRepairID, request); err != nil {
		logger.Error().Err(err).Msg("Failed to add parts supply and service to additional repair")
		appErr := mapAdditionalRepairError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusCreated, response.OperationMessageResponse{Message: "Additional repair updated successfully"})
}

// RemovePartSupplyAndService godoc
// @Summary Remove parts supply and service from additional repair
// @Description Remove parts supply and service from an existing additional repair
// @Tags Additional Repairs
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Additional Repair ID"
// @Param repair body request.AdditionalRepairItemsRequest true "Parts Supply and Service Information"
// @Success 201 {object} response.OperationMessageResponse
// @Failure 400 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /additional-repairs/{id}/remove [delete]
func (h *AdditionalRepairHandler) RemovePartSupplyAndService(c *gin.Context) {
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "AdditionalRepair/RemoveItems")
	additionalRepairID, ok := parseAdditionalRepairID(c)
	if !ok {
		return
	}

	var payload request.AdditionalRepairItemsRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		logger.Error().Err(err).Uint("ADR_ID", additionalRepairID).Msg("Failed to bind JSON for remove parts supply and service")
		c.JSON(errInvalidAdditionalRepairInput.HTTPStatus, errInvalidAdditionalRepairInput.ToHTTPError())
		return
	}

	request := toAdditionalRepairEntityFromItems(payload)
	if err := h.additionalRepairUseCase.RemovePartSupplyAndService(ctx, additionalRepairID, request); err != nil {
		logger.Error().Err(err).Uint("ADR_ID", additionalRepairID).Msg("Failed to remove parts supply and service from additional repair")
		appErr := mapAdditionalRepairError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusCreated, response.OperationMessageResponse{Message: "Additional repair updated successfully"})
}

// CustomerApproval godoc
// @Summary Update customer approval status for additional repair
// @Description Update the customer approval status of an additional repair
// @Tags Additional Repairs
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Additional Repair ID"
// @Param status body request.AdditionalRepairApprovalRequest true "Approval Status Information"
// @Success 201 {object} response.OperationMessageResponse
// @Failure 400 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /additional-repairs/{id}/customer_approval [post]
func (h *AdditionalRepairHandler) CustomerApproval(c *gin.Context) {
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "AdditionalRepair/CustomerApproval")
	additionalRepairID, ok := parseAdditionalRepairID(c)
	if !ok {
		return
	}

	var payload request.AdditionalRepairApprovalRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		logger.Error().Err(err).Uint("ADR_ID", additionalRepairID).Msg("Failed to bind JSON for customer approval")
		c.JSON(errInvalidAdditionalRepairInput.HTTPStatus, errInvalidAdditionalRepairInput.ToHTTPError())
		return
	}

	status := entities.AdditionalRepairStatusDTO{ApprovalStatus: payload.ApprovalStatus}
	if err := h.additionalRepairUseCase.CustomerApprovalStatus(ctx, additionalRepairID, status); err != nil {
		logger.Error().Err(err).Uint("ADR_ID", additionalRepairID).Msg("Failed to update customer approval status for additional repair")
		appErr := mapAdditionalRepairError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusCreated, response.OperationMessageResponse{Message: "Additional repair updated successfully"})
}

func parseAdditionalRepairID(c *gin.Context) (uint, bool) {
	logger := logs.Logger()
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		logger.Error().Err(err).Str("id", c.Param("id")).Msg("Failed to parse additional repair ID")
		c.JSON(errInvalidAdditionalRepairID.HTTPStatus, errInvalidAdditionalRepairID.ToHTTPError())
		return 0, false
	}
	return uint(id), true
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
			ID: ps.ID,
			// QuantityReserve: ps.QuantityReserve,
		})
	}
	return result
}

func toAdditionalRepairResponse(entity entities.AdditionalRepair) response.AdditionalRepairResponse {
	res := response.AdditionalRepairResponse{
		ID:             entity.ID,
		ServiceOrderID: entity.ServiceOrderID,
		Description:    entity.Description,
		Status:         entity.ARStatus.String(),
		Estimate:       entity.Estimate,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
		Services:       make([]response.AdditionalRepairServiceResponse, 0, len(entity.Services)),
		PartsSupplies:  make([]response.AdditionalRepairPartsSupplyResponse, 0, len(entity.PartsSupplies)),
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
			ID:              ps.ID,
			Price:           ps.Price,
			// QuantityReserve: ps.QuantityReserve,
		})
	}

	return res
}