package routes

import (
	handlers "github.com/fiap-grupo95/os-service-api/internal/adapter/http/handlers"

	"github.com/gin-gonic/gin"
)

func addServiceOrderRoutes(rg *gin.Engine, serviceOrderHandler *handlers.ServiceOrderHandler) {
	// Search OS
	rg.GET(GetServiceOrder, serviceOrderHandler.GetServiceOrder)
	rg.GET(GetServiceOrderHistory, serviceOrderHandler.GetServiceOrderHistory)

	// Create OS
	rg.POST(PostServiceOrderCreate, serviceOrderHandler.CreateServiceOrder)

	// Cancel OS
	rg.POST(PostServiceOrderCancel, serviceOrderHandler.CancelServiceOrder)

	// Diagnosis OS
	rg.POST(PostServiceOrderDiagnosis, serviceOrderHandler.DiagnosisServiceOrder)
	rg.POST(PostServiceOrderDiagnosisFinish, serviceOrderHandler.FinishServiceOrderDiagnosis)

	// Estimate OS
	rg.POST(PostServiceOrderEstimateApprove, serviceOrderHandler.ApproveServiceOrderEstimate)
	rg.POST(PostServiceOrderEstimateReject, serviceOrderHandler.RejectServiceOrderEstimate)
	rg.POST(PostServiceOrderEstimateCancel, serviceOrderHandler.CancelServiceOrderEstimate)

	// Execution OS
	rg.POST(PostServiceOrderExecutionCreate, serviceOrderHandler.ExecutionServiceOrder)
	rg.POST(PostServiceOrderExecutionFinish, serviceOrderHandler.FinishServiceOrderExecution)

	// Delivery OS
	rg.POST(PostServiceOrderDelivery, serviceOrderHandler.DeliveryServiceOrder)
	rg.POST(PostServiceOrderPayment, serviceOrderHandler.PaymentServiceOrder)
}
