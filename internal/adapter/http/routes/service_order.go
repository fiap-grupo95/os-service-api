package routes

import (
	handlers "github.com/fiap-grupo95/os-service-api/internal/adapter/http/handlers"

	"github.com/gin-gonic/gin"
)

func addServiceOrderRoutes(rg *gin.Engine, serviceOrderHandler *handlers.ServiceOrderHandler) {
	// Get endpoints
	rg.GET(GetServiceOrder, serviceOrderHandler.GetServiceOrder)
	rg.GET(GetServiceOrderHistory, serviceOrderHandler.GetServiceOrderHistory)

	// Patch endpoints
	rg.PATCH(PatchServiceOrderDiagnosis, serviceOrderHandler.UpdateServiceOrderDiagnosis)

	// Post endpoints
	rg.POST(PostServiceOrderCreate, serviceOrderHandler.CreateServiceOrder)
	rg.POST(PostServiceOrderCancel, serviceOrderHandler.CancelServiceOrder)
	rg.POST(PostServiceOrderDiagnosis, serviceOrderHandler.FinishServiceOrderDiagnosis)
	rg.POST(PostServiceOrderEstimateApprove, serviceOrderHandler.ApproveServiceOrderEstimate)
	rg.POST(PostServiceOrderEstimateReject, serviceOrderHandler.RejectServiceOrderEstimate)
	rg.POST(PostServiceOrderEstimateCancel, serviceOrderHandler.CancelServiceOrderEstimate)
	rg.POST(PostServiceOrderExecutionCreate, serviceOrderHandler.CreateServiceOrderExecution)
	rg.POST(PostServiceOrderExecutionFinish, serviceOrderHandler.FinishServiceOrderExecution)
	rg.POST(PostServiceOrderPayment, serviceOrderHandler.CreateServiceOrderPayment)
	rg.POST(PostServiceOrderDelivery, serviceOrderHandler.DeliveryServiceOrder)
}
