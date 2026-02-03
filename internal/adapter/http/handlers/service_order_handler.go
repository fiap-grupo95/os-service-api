package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	request "github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/request"
	response "github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/logs"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/observability"
	"github.com/fiap-grupo95/os-service-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

const (
	diagnosisFlow = "diagnosis"
	estimateFlow  = "estimate"
	executionFlow = "execution"
	deliveryFlow  = "delivery"
	getFlow       = "get"
	listFlow      = "list"
	createFlow    = "create"
)

const (
	metricServiceOrderCreate       = "service_order_create"
	metricServiceOrderStatusChange = "service_order_status_change"
	resultSuccess                  = "success"
	resultError                    = "error"
)

type IServiceOrderHandler interface {
	CreateServiceOrder(c *gin.Context)
	CancelServiceOrder(c *gin.Context)
	UpdateServiceOrderDiagnosis(c *gin.Context)
	FinishServiceOrderDiagnosis(c *gin.Context)
	ApproveServiceOrderEstimate(c *gin.Context)
	RejectServiceOrderEstimate(c *gin.Context)
	CancelServiceOrderEstimate(c *gin.Context)
	CreateServiceOrderExecution(c *gin.Context)
	FinishServiceOrderExecution(c *gin.Context)
	CreateServiceOrderPayment(c *gin.Context)
	DeliveryServiceOrder(c *gin.Context)
	GetServiceOrder(c *gin.Context)
	GetServiceOrderHistory(c *gin.Context)
	ListServiceOrders(c *gin.Context)
}

type ServiceOrderHandler struct {
	serviceOrderUseCase usecase.IServiceOrderUseCase
}

func NewServiceOrderHandler(useCase usecase.IServiceOrderUseCase) *ServiceOrderHandler {
	return &ServiceOrderHandler{
		serviceOrderUseCase: useCase,
	}
}

func (h *ServiceOrderHandler) CancelServiceOrder(c *gin.Context) {
	// TODO: Implement this method
}

func (h *ServiceOrderHandler) FinishServiceOrderDiagnosis(c *gin.Context) {
	// TODO: Implement this method
}

func (h *ServiceOrderHandler) ApproveServiceOrderEstimate(c *gin.Context) {
	// TODO: Implement this method
}

func (h *ServiceOrderHandler) RejectServiceOrderEstimate(c *gin.Context) {
	// TODO: Implement this method
}

func (h *ServiceOrderHandler) CancelServiceOrderEstimate(c *gin.Context) {
	// TODO: Implement this method
}

func (h *ServiceOrderHandler) CreateServiceOrderExecution(c *gin.Context) {
	// TODO: Implement this method
}

func (h *ServiceOrderHandler) FinishServiceOrderExecution(c *gin.Context) {
	// TODO: Implement this method
}

func (h *ServiceOrderHandler) CreateServiceOrderPayment(c *gin.Context) {
	// TODO: Implement this method
}

func (h *ServiceOrderHandler) DeliveryServiceOrder(c *gin.Context) {
	// TODO: Implement this method
}

func (h *ServiceOrderHandler) GetServiceOrderHistory(c *gin.Context) {
	// TODO: Implement this method
}

// CreateServiceOrder godoc
// @Summary Create a new service order
// @Description Create a new service order record
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param order body request.ServiceOrderCreateRequest true "Service Order Information"
// @Success 201 {object} response.ServiceOrderResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/service-orders/create [post]
func (h *ServiceOrderHandler) CreateServiceOrder(c *gin.Context) {
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "ServiceOrder/Create")

	var req request.ServiceOrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendToMetric(ctx, metricServiceOrderCreate, resultError, "", "", strconv.Itoa(http.StatusBadRequest))

		logger.Error().Err(err).Msg("Failed to bind JSON for create service order")

		writeBindingError(c, err)
		return
	}
	result, err := h.serviceOrderUseCase.CreateServiceOrder(ctx, req.ToEntity())
	if err != nil {
		sendToMetric(ctx, metricServiceOrderCreate, resultError, "", "", strconv.Itoa(getStatusError(err)))

		logger.Error().Err(err).Msg("Failed to create service order")

		writeServiceOrderError(c, err, "Failed to create service order")
		return
	}

	sendToMetric(ctx, metricServiceOrderCreate, resultSuccess, "", "", strconv.Itoa(http.StatusCreated))

	c.JSON(http.StatusCreated, response.NewServiceOrderResponse(result))
}

// UpdateServiceOrderDiagnosis godoc
// @Summary Update service order diagnosis
// @Description Update the diagnosis information of a service order
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Order ID"
// @Param order body request.ServiceOrderDiagnosisUpdateRequest true "Service Order Diagnosis Information"
// @Success 200 {object} response.ServiceOrderResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /service-orders/{id}/diagnosis [patch]
func (h *ServiceOrderHandler) UpdateServiceOrderDiagnosis(c *gin.Context) {
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "ServiceOrder/UpdateDiagnosis")
	id, ok := parseServiceOrderIDParam(c)
	if !ok {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, diagnosisFlow, "", strconv.Itoa(http.StatusBadRequest))
		return
	}

	var req request.ServiceOrderDiagnosisUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, diagnosisFlow, req.ServiceOrderStatus, strconv.Itoa(http.StatusBadRequest))
		
		logger.Error().Err(err).Uint("OS_ID", id).Msg("Failed to bind JSON for update service order diagnosis")
		
		writeBindingError(c, err)
		return
	}

	result, err := h.serviceOrderUseCase.UpdateServiceOrder(ctx, req.ToEntity(id), diagnosisFlow)
	if err != nil {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, diagnosisFlow, req.ServiceOrderStatus, strconv.Itoa(getStatusError(err)))
		
		logger.Error().Err(err).Uint("OS_ID", id).Msg("Failed to update service order diagnosis")
		
		writeServiceOrderError(c, err, "Failed to update service order diagnosis")
		return
	}
	sendToMetric(ctx, metricServiceOrderStatusChange, resultSuccess, diagnosisFlow, req.ServiceOrderStatus, strconv.Itoa(http.StatusOK))
	c.JSON(http.StatusOK, response.NewServiceOrderResponse(result))
}

// UpdateServiceOrderEstimate godoc
// @Summary Update service order estimate
// @Description Update the estimate information of a service order
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Order ID"
// @Param order body request.ServiceOrderEstimateUpdateRequest true "Service Order Estimate Information"
// @Success 200 {object} response.ServiceOrderResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /service-orders/{id}/estimate [patch]
func (h *ServiceOrderHandler) UpdateServiceOrderEstimate(c *gin.Context) {
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "ServiceOrder/UpdateEstimate")
	id, ok := parseServiceOrderIDParam(c)
	if !ok {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, estimateFlow, "", strconv.Itoa(http.StatusBadRequest))
		return
	}

	var req request.ServiceOrderEstimateUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, estimateFlow, req.ServiceOrderStatus, strconv.Itoa(http.StatusBadRequest))
		
		logger.Error().Err(err).Uint("OS_ID", id).Msg("Failed to bind JSON for update service order estimate")
		
		writeBindingError(c, err)
		return
	}

	result, err := h.serviceOrderUseCase.UpdateServiceOrder(ctx, req.ToEntity(id), estimateFlow)
	if err != nil {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, estimateFlow, req.ServiceOrderStatus, strconv.Itoa(getStatusError(err)))
		
		logger.Error().Err(err).Uint("OS_ID", id).Msg("Failed to update service order estimate")
		
		writeServiceOrderError(c, err, "Failed to update service order estimate")
		return
	}
	sendToMetric(ctx, metricServiceOrderStatusChange, resultSuccess, estimateFlow, req.ServiceOrderStatus, strconv.Itoa(http.StatusOK))
	c.JSON(http.StatusOK, response.NewServiceOrderResponse(result))
}

// UpdateServiceOrderExecution godoc
// @Summary Update service order execution
// @Description Update the execution information of a service order
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Order ID"
// @Param order body request.ServiceOrderExecutionUpdateRequest true "Service Order Execution Information"
// @Success 200 {object} response.ServiceOrderResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /service-orders/{id}/execution [patch]
func (h *ServiceOrderHandler) UpdateServiceOrderExecution(c *gin.Context) {
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "ServiceOrder/UpdateExecution")
	id, ok := parseServiceOrderIDParam(c)
	if !ok {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, executionFlow, "", strconv.Itoa(http.StatusBadRequest))
		return
	}

	var req request.ServiceOrderExecutionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, executionFlow, req.ServiceOrderStatus, strconv.Itoa(http.StatusBadRequest))
		
		logger.Error().Err(err).Uint("OS_ID", id).Msg("Failed to bind JSON for update service order execution")
		
		writeBindingError(c, err)
		return
	}

	result, err := h.serviceOrderUseCase.UpdateServiceOrder(ctx, req.ToEntity(id), executionFlow)
	if err != nil {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, executionFlow, req.ServiceOrderStatus, strconv.Itoa(getStatusError(err)))
		logger.Error().Err(err).Uint("OS_ID", id).Msg("Failed to update service order execution")
		writeServiceOrderError(c, err, "Failed to update service order execution")
		return
	}
	sendToMetric(ctx, metricServiceOrderStatusChange, resultSuccess, executionFlow, req.ServiceOrderStatus, strconv.Itoa(http.StatusOK))
	c.JSON(http.StatusOK, response.NewServiceOrderResponse(result))
}

// UpdateServiceOrderDelivery godoc
// @Summary Update service order delivery
// @Description Update the delivery information of a service order
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Order ID"
// @Param order body request.ServiceOrderDeliveryUpdateRequest true "Service Order Delivery Information"
// @Success 200 {object} response.ServiceOrderResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /service-orders/{id}/delivery [patch]
func (h *ServiceOrderHandler) UpdateServiceOrderDelivery(c *gin.Context) {
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "ServiceOrder/UpdateDelivery")
	id, ok := parseServiceOrderIDParam(c)
	if !ok {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, deliveryFlow, "", strconv.Itoa(http.StatusBadRequest))
		return
	}

	var req request.ServiceOrderDeliveryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, deliveryFlow, req.ServiceOrderStatus, strconv.Itoa(http.StatusBadRequest))
		
		logger.Error().Err(err).Uint("OS_ID", id).Msg("Failed to bind JSON for update service order delivery")
		
		writeBindingError(c, err)
		return
	}

	result, err := h.serviceOrderUseCase.UpdateServiceOrder(ctx, req.ToEntity(id), deliveryFlow)
	if err != nil {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, deliveryFlow, req.ServiceOrderStatus, strconv.Itoa(getStatusError(err)))
		logger.Error().Err(err).Uint("OS_ID", id).Msg("Failed to update service order delivery")
		writeServiceOrderError(c, err, "Failed to update service order delivery")
		return
	}
	sendToMetric(ctx, metricServiceOrderStatusChange, resultSuccess, deliveryFlow, req.ServiceOrderStatus, strconv.Itoa(http.StatusOK))
	c.JSON(http.StatusOK, response.NewServiceOrderResponse(result))
}

// GetServiceOrder godoc
// @Summary Get service order by ID
// @Description Retrieve a service order by its ID
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Order ID"
// @Success 200 {object} response.ServiceOrderResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /service-orders/{id} [get]
func (h *ServiceOrderHandler) GetServiceOrder(c *gin.Context) {
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "ServiceOrder/Get")
	id, ok := parseServiceOrderIDParam(c)
	if !ok {
		return
	}

	serviceOrder := entities.ServiceOrder{ID: id}
	result, err := h.serviceOrderUseCase.GetServiceOrder(ctx, serviceOrder)
	if err != nil {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, getFlow, "", strconv.Itoa(getStatusError(err)))
		
		logger.Error().Err(err).Uint("OS_ID", id).Msg("Failed to retrieve service order")

		writeServiceOrderError(c, err, "Failed to retrieve service order")
		return
	}
	sendToMetric(ctx, metricServiceOrderStatusChange, resultSuccess, getFlow, "", strconv.Itoa(http.StatusOK))
	c.JSON(http.StatusOK, response.NewServiceOrderResponse(result))
}

// ListServiceOrders godoc
// @Summary List all service orders
// @Description Get a list of all service orders
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Success 200 {array} response.ServiceOrderResponse
// @Failure 500 {object} map[string]string
// @Router /service-orders [get]
func (h *ServiceOrderHandler) ListServiceOrders(c *gin.Context) {
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "ServiceOrder/List")
	serviceOrders, err := h.serviceOrderUseCase.ListServiceOrders(ctx)
	if err != nil {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, listFlow, "", strconv.Itoa(getStatusError(err)))
		
		logger.Error().Err(err).Msg("Failed to retrieve service orders")

		writeServiceOrderError(c, err, "Failed to retrieve service orders")
		return
	}
	sendToMetric(ctx, metricServiceOrderStatusChange, resultSuccess, listFlow, "", strconv.Itoa(http.StatusOK))
	c.JSON(http.StatusOK, response.NewServiceOrderListResponse(serviceOrders))
}

func parseServiceOrderIDParam(c *gin.Context) (uint, bool) {
	logger := logs.Logger()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		logger.Error().Err(err).Str("id", c.Param("id")).Msg("Invalid service order ID")
		
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service order ID"})
		return 0, false
	}
	return uint(id), true
}

func writeBindingError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error":   "Invalid input",
		"details": err.Error(),
	})
}

func writeServiceOrderError(c *gin.Context, err error, fallbackMessage string) {
	switch {
	case errors.Is(err, usecase.ErrServiceOrderNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "Service order not found"})
	case errors.Is(err, usecase.ErrVehicleNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
	case errors.Is(err, usecase.ErrCustomerNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
	case errors.Is(err, usecase.ErrServiceNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
	case errors.Is(err, usecase.ErrPartsSupplyNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "Parts supply not found"})
	case errors.Is(err, usecase.ErrInvalidCustomerID),
		errors.Is(err, usecase.ErrInvalidID):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case isBadRequestServiceOrderError(err):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		logs.Logger().Error().Err(err).Msg(fallbackMessage)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   fallbackMessage,
			"details": err.Error(),
		})
	}
}

func isBadRequestServiceOrderError(err error) bool {
	switch {
	case errors.Is(err, usecase.ErrInvalidTransitionStatusToDiagnosis),
		errors.Is(err, usecase.ErrInvalidTransitionStatusToEstimate),
		errors.Is(err, usecase.ErrInvalidTransitionStatusToExecution),
		errors.Is(err, usecase.ErrInvalidTransitionStatusToDelivery),
		errors.Is(err, usecase.ErrInvalidStatus),
		errors.Is(err, usecase.ErrInvalidFlow),
		errors.Is(err, usecase.ErrInsufficientPartsSupply):
		return true
	default:
		return false
	}
}

func retrieveTransactioAndContext(c *gin.Context, name string) context.Context {
	// 1) Recupera a transaction criada pelo middleware nrgin
	txn := nrgin.Transaction(c)
	if txn != nil {
		// Opcional: dar um nome mais semântico para o trace
		txn.SetName(name)
		// Opcional: adicionar atributos úteis ao trace
		txn.AddAttribute("http.method", c.Request.Method)
		txn.AddAttribute("http.route", c.FullPath())
	}
	// 2) Propaga a transaction para o context que vai para o usecase
	ctx := c.Request.Context()
	if txn != nil {
		ctx = newrelic.NewContext(ctx, txn)
	}
	return ctx
}

func getStatusError(err error) int {
	switch {
	case errors.Is(err, usecase.ErrServiceOrderNotFound):
		return http.StatusNotFound
	case errors.Is(err, usecase.ErrVehicleNotFound):
		return http.StatusNotFound
	case errors.Is(err, usecase.ErrCustomerNotFound):
		return http.StatusNotFound
	case errors.Is(err, usecase.ErrServiceNotFound):
		return http.StatusNotFound
	case errors.Is(err, usecase.ErrPartsSupplyNotFound):
		return http.StatusNotFound
	case errors.Is(err, usecase.ErrInvalidCustomerID),
		errors.Is(err, usecase.ErrInvalidID):
		return http.StatusBadRequest
	case isBadRequestServiceOrderError(err):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func sendToMetric(ctx context.Context, metricName string, result string, flow string, status string, code string) {
	var tags map[string]string
	
	if result != "" {
		tags["result"] = result
	}
	if flow != "" {
		tags["flow"] = flow
	}
	if status != "" {
		tags["status"] = status
	}
	if code != "" {
		tags["code"] = code
	}
	
	observability.IncrementCounter(ctx, metricName, tags)
}