package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"

	request "github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/request"
	response "github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/logs"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/observability"
	"github.com/fiap-grupo95/os-service-api/internal/usecase"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/constants"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

const (
	getFlow    = "get"
	listFlow   = "list"
	createFlow = "create"
)

const (
	metricServiceOrderCreate       = "service_order_create"
	metricServiceOrderStatusChange = "service_order_status_change"
	resultSuccess                  = "success"
	resultError                    = "error"
)

type IServiceOrderHandler interface {
	GetServiceOrder(c *gin.Context)
	ListServiceOrders(c *gin.Context)
	CreateServiceOrder(c *gin.Context)
	CancelServiceOrder(c *gin.Context)
	DiagnosisServiceOrder(c *gin.Context)
	ApproveServiceOrderEstimate(c *gin.Context)
	RejectServiceOrderEstimate(c *gin.Context)
	CancelServiceOrderEstimate(c *gin.Context)
	ExecutionServiceOrder(c *gin.Context)
	FinishServiceOrderExecution(c *gin.Context)
	PaymentServiceOrder(c *gin.Context)
	DeliveryServiceOrder(c *gin.Context)
}

type ServiceOrderHandler struct {
	serviceOrderUseCase usecase.IServiceOrderUseCase
}

func NewServiceOrderHandler(useCase usecase.IServiceOrderUseCase) *ServiceOrderHandler {
	return &ServiceOrderHandler{
		serviceOrderUseCase: useCase,
	}
}

// GetServiceOrder godoc
// @Summary Retrieve service order details
// @Description Fetches a specific service order by its unique identifier. Use the 'full_data' query parameter to include complete related entity information (customer, vehicle, services, parts, additional repairs).
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Order ID" minimum(1)
// @Param full_data query boolean false "Include complete related entity data" default(false)
// @Success 200 {object} response.ServiceOrderResponse "Service order retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid service order ID format"
// @Failure 404 {object} map[string]string "Service order not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/service-orders/{id} [get]
func (h *ServiceOrderHandler) GetServiceOrder(c *gin.Context) {
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "ServiceOrder/Get")
	id := parseServiceOrderIDParam(c)
	if id == "" {
		return
	}

	isFullData := parseIsFullDataParam(c)

	serviceOrder := entities.ServiceOrder{ID: id}
	result, err := h.serviceOrderUseCase.GetServiceOrder(ctx, serviceOrder, isFullData)
	if err != nil {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, getFlow, "", strconv.Itoa(getStatusError(err)))

		logger.Error().Err(err).Str("OS_ID", id).Msg("Failed to retrieve service order")

		writeServiceOrderError(c, err, "Failed to retrieve service order")
		return
	}
	sendToMetric(ctx, metricServiceOrderStatusChange, resultSuccess, getFlow, "", strconv.Itoa(http.StatusOK))
	c.JSON(http.StatusOK, response.NewServiceOrderResponse(result))
}

// ListServiceOrders godoc
// @Summary List all service orders
// @Description Retrieves a complete list of all service orders in the system. Returns basic information for each service order without detailed related entities.
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Success 200 {array} response.ServiceOrderResponse "List of service orders retrieved successfully"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/service-orders [get]
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

// CreateServiceOrder godoc
// @Summary Create a new service order
// @Description Creates a new service order with initial status 'PENDING'. Requires valid customer ID, vehicle ID, and at least one service. The service order will be initialized with creation timestamp and default status.
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param order body request.ServiceOrderCreateRequest true "Service order creation payload with customer, vehicle, and service details"
// @Success 201 {object} response.ServiceOrderResponse "Service order created successfully"
// @Failure 400 {object} map[string]string "Invalid request payload or validation error"
// @Failure 404 {object} map[string]string "Referenced customer, vehicle, or service not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/service-orders/create [post]
func (h *ServiceOrderHandler) CreateServiceOrder(c *gin.Context) {
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "ServiceOrder/Create")

	var req request.ServiceOrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendToMetric(ctx, metricServiceOrderCreate, resultError, createFlow, "", strconv.Itoa(http.StatusBadRequest))

		logger.Error().Err(err).Msg("Failed to bind JSON for create service order")

		writeBindingError(c, err)
		return
	}
	result, err := h.serviceOrderUseCase.CreateServiceOrder(ctx, req.ToEntity())
	if err != nil {
		sendToMetric(ctx, metricServiceOrderCreate, resultError, createFlow, "", strconv.Itoa(getStatusError(err)))

		logger.Error().Err(err).Msg("Failed to create service order")

		writeServiceOrderError(c, err, "Failed to create service order")
		return
	}

	sendToMetric(ctx, metricServiceOrderCreate, resultSuccess, createFlow, "", strconv.Itoa(http.StatusCreated))

	c.JSON(http.StatusCreated, response.NewServiceOrderResponse(result))
}

// CancelServiceOrder godoc
// @Summary Cancel a service order
// @Description Cancels a service order, transitioning status to 'CANCELLED'. This action can be performed on service orders that are not yet finalized, delivered, or already cancelled. Once cancelled, the service order cannot be resumed.
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Order ID" minimum(1)
// @Success 200 {object} response.ServiceOrderResponse "Service order cancelled successfully"
// @Failure 400 {object} map[string]string "Invalid ID or invalid status transition"
// @Failure 404 {object} map[string]string "Service order not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/service-orders/{id}/cancel [post]
func (h *ServiceOrderHandler) CancelServiceOrder(c *gin.Context) {
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "ServiceOrder/Cancel")
	id := parseServiceOrderIDParam(c)
	if id == "" {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, constants.CANCEL, "", strconv.Itoa(http.StatusBadRequest))
		return
	}
	result, err := h.serviceOrderUseCase.CancelServiceOrder(ctx, id)
	if err != nil {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, constants.CANCEL, "", strconv.Itoa(getStatusError(err)))

		logger.Error().Err(err).Str("OS_ID", id).Msg("Failed to cancel service order")

		writeServiceOrderError(c, err, "Failed to cancel service order")
		return
	}
	sendToMetric(ctx, metricServiceOrderStatusChange, resultSuccess, constants.CANCEL, result.Status.String(), strconv.Itoa(http.StatusOK))
	c.JSON(http.StatusOK, response.NewServiceOrderResponse(result))
}

// DiagnosisServiceOrder godoc
// @Summary Submit service order diagnosis
// @Description Updates the service order with diagnosis information and transitions status to 'DIAGNOSIS'. Requires the service order to be in a valid state for diagnosis. Includes diagnosis description, estimated cost, and required parts.
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Order ID" minimum(1)
// @Param order body request.ServiceOrderDiagnosisUpdateRequest true "Diagnosis details including description, estimated cost, and parts list"
// @Success 200 {object} response.ServiceOrderResponse "Diagnosis submitted successfully"
// @Failure 400 {object} map[string]string "Invalid ID, payload, or invalid status transition"
// @Failure 404 {object} map[string]string "Service order or referenced parts not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/service-orders/{id}/diagnosis [post]
func (h *ServiceOrderHandler) DiagnosisServiceOrder(c *gin.Context) {
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "ServiceOrder/Diagnosis")
	id := parseServiceOrderIDParam(c)
	if id == "" {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, constants.DIAGNOSIS, "", strconv.Itoa(http.StatusBadRequest))
		return
	}

	var req request.ServiceOrderDiagnosisUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if errors.Is(err, io.EOF) {
			req = request.ServiceOrderDiagnosisUpdateRequest{}
		} else {
			sendToMetric(ctx, metricServiceOrderStatusChange, resultError, constants.DIAGNOSIS, "", strconv.Itoa(http.StatusBadRequest))
			logger.Error().Err(err).Str("OS_ID", id).Msg("Failed to bind JSON for update service order diagnosis")
			writeBindingError(c, err)
			return
		}
	}

	result, err := h.serviceOrderUseCase.DiagnosisServiceOrder(ctx, req.ToEntity(id))
	if err != nil {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, constants.DIAGNOSIS, "", strconv.Itoa(getStatusError(err)))

		logger.Error().Err(err).Str("OS_ID", id).Msg("Failed to update service order diagnosis")

		writeServiceOrderError(c, err, "Failed to update service order diagnosis")
		return
	}
	sendToMetric(ctx, metricServiceOrderStatusChange, resultSuccess, constants.DIAGNOSIS, "", strconv.Itoa(http.StatusOK))
	c.JSON(http.StatusOK, response.NewServiceOrderResponse(result))
}

// ApproveServiceOrderEstimate godoc
// @Summary Approve service order estimate
// @Description Customer approves the service order estimate, transitioning the status to 'APPROVED'. This allows the service order to proceed to execution phase. Must be in 'DIAGNOSIS' status.
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Order ID" minimum(1)
// @Success 200 {object} response.ServiceOrderResponse "Estimate approved successfully"
// @Failure 400 {object} map[string]string "Invalid ID or invalid status transition"
// @Failure 404 {object} map[string]string "Service order not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/service-orders/{id}/estimate/approve [post]
func (h *ServiceOrderHandler) ApproveServiceOrderEstimate(c *gin.Context) {
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "ServiceOrder/Estimate")
	id := parseServiceOrderIDParam(c)
	if id == "" {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, constants.ESTIMATE_APPROVE, "", strconv.Itoa(http.StatusBadRequest))
		return
	}
	result, err := h.serviceOrderUseCase.EstimateServiceOrder(ctx, id, constants.ESTIMATE_APPROVE)
	if err != nil {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, constants.ESTIMATE_APPROVE, "", strconv.Itoa(getStatusError(err)))

		logger.Error().Err(err).Str("OS_ID", id).Msg("Failed to approve service order estimate")

		writeServiceOrderError(c, err, "Failed to approve service order estimate")
		return
	}
	sendToMetric(ctx, metricServiceOrderStatusChange, resultSuccess, constants.ESTIMATE_APPROVE, result.Status.String(), strconv.Itoa(http.StatusOK))
	c.JSON(http.StatusOK, response.NewServiceOrderResponse(result))
}

// RejectServiceOrderEstimate godoc
// @Summary Reject service order estimate
// @Description Customer rejects the service order estimate, transitioning the status to 'REJECTED'. The service order will not proceed to execution. Must be in 'DIAGNOSIS' status.
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Order ID" minimum(1)
// @Success 200 {object} response.ServiceOrderResponse "Estimate rejected successfully"
// @Failure 400 {object} map[string]string "Invalid ID or invalid status transition"
// @Failure 404 {object} map[string]string "Service order not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/service-orders/{id}/estimate/reject [post]
func (h *ServiceOrderHandler) RejectServiceOrderEstimate(c *gin.Context) {
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "ServiceOrder/Estimate")
	id := parseServiceOrderIDParam(c)
	if id == "" {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, constants.ESTIMATE_REJECT, "", strconv.Itoa(http.StatusBadRequest))
		return
	}
	result, err := h.serviceOrderUseCase.EstimateServiceOrder(ctx, id, constants.ESTIMATE_REJECT)
	if err != nil {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, constants.ESTIMATE_REJECT, "", strconv.Itoa(getStatusError(err)))

		logger.Error().Err(err).Str("OS_ID", id).Msg("Failed to reject service order estimate")

		writeServiceOrderError(c, err, "Failed to reject service order estimate")
		return
	}
	sendToMetric(ctx, metricServiceOrderStatusChange, resultSuccess, constants.ESTIMATE_REJECT, result.Status.String(), strconv.Itoa(http.StatusOK))
	c.JSON(http.StatusOK, response.NewServiceOrderResponse(result))
}

// CancelServiceOrderEstimate godoc
// @Summary Cancel service order estimate
// @Description Cancels the service order estimate process, transitioning the status to 'CANCELLED'. This action is typically initiated by the customer or service provider. Must be in 'DIAGNOSIS' status.
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Order ID" minimum(1)
// @Success 200 {object} response.ServiceOrderResponse "Estimate cancelled successfully"
// @Failure 400 {object} map[string]string "Invalid ID or invalid status transition"
// @Failure 404 {object} map[string]string "Service order not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/service-orders/{id}/estimate/cancel [post]
func (h *ServiceOrderHandler) CancelServiceOrderEstimate(c *gin.Context) {
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "ServiceOrder/Estimate")
	id := parseServiceOrderIDParam(c)
	if id == "" {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, constants.ESTIMATE_CANCEL, "", strconv.Itoa(http.StatusBadRequest))
		return
	}
	result, err := h.serviceOrderUseCase.EstimateServiceOrder(ctx, id, constants.ESTIMATE_CANCEL)
	if err != nil {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, constants.ESTIMATE_CANCEL, "", strconv.Itoa(getStatusError(err)))

		logger.Error().Err(err).Str("OS_ID", id).Msg("Failed to cancel service order estimate")

		writeServiceOrderError(c, err, "Failed to cancel service order estimate")
		return
	}
	sendToMetric(ctx, metricServiceOrderStatusChange, resultSuccess, constants.ESTIMATE_CANCEL, result.Status.String(), strconv.Itoa(http.StatusOK))
	c.JSON(http.StatusOK, response.NewServiceOrderResponse(result))
}

// ExecutionServiceOrder godoc
// @Summary Start service order execution
// @Description Initiates the execution phase of the service order, transitioning status to 'IN_EXECUTION'. This indicates that work has begun on the vehicle. Must be in 'APPROVED' status.
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Order ID" minimum(1)
// @Success 200 {object} response.ServiceOrderResponse "Execution started successfully"
// @Failure 400 {object} map[string]string "Invalid ID or invalid status transition"
// @Failure 404 {object} map[string]string "Service order not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/service-orders/{id}/execution [post]
func (h *ServiceOrderHandler) ExecutionServiceOrder(c *gin.Context) {
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "ServiceOrder/Execution")
	id := parseServiceOrderIDParam(c)
	if id == "" {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, constants.EXECUTION_START, "", strconv.Itoa(http.StatusBadRequest))
		return
	}

	result, err := h.serviceOrderUseCase.ExecutionServiceOrder(ctx, id, constants.EXECUTION_START)
	if err != nil {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, constants.EXECUTION_START, "", strconv.Itoa(getStatusError(err)))
		logger.Error().Err(err).Str("OS_ID", id).Msg("Failed to update service order execution")
		writeServiceOrderError(c, err, "Failed to update service order execution")
		return
	}
	sendToMetric(ctx, metricServiceOrderStatusChange, resultSuccess, constants.EXECUTION_START, result.Status.String(), strconv.Itoa(http.StatusOK))
	c.JSON(http.StatusOK, response.NewServiceOrderResponse(result))
}

// FinishServiceOrderExecution godoc
// @Summary Complete service order execution
// @Description Marks the execution phase as complete, transitioning status to 'EXECUTED'. This indicates that all work on the vehicle has been finished and the service order is ready for payment. Must be in 'IN_EXECUTION' status.
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Order ID" minimum(1)
// @Success 200 {object} response.ServiceOrderResponse "Execution completed successfully"
// @Failure 400 {object} map[string]string "Invalid ID or invalid status transition"
// @Failure 404 {object} map[string]string "Service order not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/service-orders/{id}/execution/finish [post]
func (h *ServiceOrderHandler) FinishServiceOrderExecution(c *gin.Context) {
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "ServiceOrder/Execution")
	id := parseServiceOrderIDParam(c)
	if id == "" {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, constants.EXECUTION_FINISH, "", strconv.Itoa(http.StatusBadRequest))
		return
	}

	result, err := h.serviceOrderUseCase.ExecutionServiceOrder(ctx, id, constants.EXECUTION_FINISH)
	if err != nil {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, constants.EXECUTION_FINISH, "", strconv.Itoa(getStatusError(err)))
		logger.Error().Err(err).Str("OS_ID", id).Msg("Failed to update service order execution")
		writeServiceOrderError(c, err, "Failed to update service order execution")
		return
	}
	sendToMetric(ctx, metricServiceOrderStatusChange, resultSuccess, constants.EXECUTION_FINISH, result.Status.String(), strconv.Itoa(http.StatusOK))
	c.JSON(http.StatusOK, response.NewServiceOrderResponse(result))
}

// PaymentServiceOrder godoc
// @Summary Process service order payment
// @Description Processes payment for the service order, transitioning status to 'PAID'. This triggers billing service integration and marks the service order as financially settled. Must be in 'EXECUTED' status.
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Order ID" minimum(1)
// @Success 200 {object} response.ServiceOrderResponse "Payment processed successfully"
// @Failure 400 {object} map[string]string "Invalid ID or invalid status transition"
// @Failure 404 {object} map[string]string "Service order not found"
// @Failure 500 {object} map[string]string "Internal server error or billing service failure"
// @Router /v1/service-orders/{id}/payment [post]
func (h *ServiceOrderHandler) PaymentServiceOrder(c *gin.Context) {
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "ServiceOrder/Payment")
	id := parseServiceOrderIDParam(c)
	if id == "" {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, constants.PAYMENT, "", strconv.Itoa(http.StatusBadRequest))
		return
	}

	result, err := h.serviceOrderUseCase.PaymentServiceOrder(ctx, id)
	if err != nil {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, constants.PAYMENT, "", strconv.Itoa(getStatusError(err)))
		logger.Error().Err(err).Str("OS_ID", id).Msg("Failed to create service order payment")
		writeServiceOrderError(c, err, "Failed to create service order payment")
		return
	}
	sendToMetric(ctx, metricServiceOrderStatusChange, resultSuccess, constants.PAYMENT, "", strconv.Itoa(http.StatusOK))
	c.JSON(http.StatusOK, result)
}

// DeliveryServiceOrder godoc
// @Summary Complete service order delivery
// @Description Marks the service order as delivered to the customer, transitioning status to 'DELIVERED'. This is the final step in the service order lifecycle. Must be in 'PAID' status.
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Order ID" minimum(1)
// @Success 200 {object} response.ServiceOrderResponse "Service order delivered successfully"
// @Failure 400 {object} map[string]string "Invalid ID or invalid status transition"
// @Failure 404 {object} map[string]string "Service order not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/service-orders/{id}/delivery [post]
func (h *ServiceOrderHandler) DeliveryServiceOrder(c *gin.Context) {
	logger := logs.Logger()
	ctx := retrieveTransactioAndContext(c, "ServiceOrder/Delivery")
	id := parseServiceOrderIDParam(c)
	if id == "" {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, constants.DELIVERY, "", strconv.Itoa(http.StatusBadRequest))
		return
	}

	result, err := h.serviceOrderUseCase.DeliveryServiceOrder(ctx, id)
	if err != nil {
		sendToMetric(ctx, metricServiceOrderStatusChange, resultError, constants.DELIVERY, "", strconv.Itoa(getStatusError(err)))
		logger.Error().Err(err).Str("OS_ID", id).Msg("Failed to update service order delivery")
		writeServiceOrderError(c, err, "Failed to update service order delivery")
		return
	}
	sendToMetric(ctx, metricServiceOrderStatusChange, resultSuccess, constants.DELIVERY, "", strconv.Itoa(http.StatusOK))
	c.JSON(http.StatusOK, response.NewServiceOrderResponse(result))
}

func parseServiceOrderIDParam(c *gin.Context) string {
	return c.Param("id")
}

func parseIsFullDataParam(c *gin.Context) bool {
	return c.Query("full_data") == "true"
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
	case errors.Is(err, usecase.ErrInvalidID):
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
	case errors.Is(err, usecase.ErrInvalidID):
		return http.StatusBadRequest
	case isBadRequestServiceOrderError(err):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func sendToMetric(ctx context.Context, metricName string, result string, flow string, status string, code string) {
	tags := make(map[string]string)

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
