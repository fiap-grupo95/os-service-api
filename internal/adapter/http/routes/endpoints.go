package routes

const (
	PathHealthCheck      = "/ping"
	PathServiceOrders    = "/service-orders"
	PathAdditionalRepair = "/additional-repair"
)

const (
	PostServiceOrderCreate          = "/v1/service-orders/create"
	PostServiceOrderCancel          = "/v1/service-orders/:id/cancel"
	PostServiceOrderDiagnosis       = "/v1/service-orders/:id/diagnosis/finish"
	PatchServiceOrderDiagnosis      = "/v1/service-orders/:id/diagnosis"
	PostServiceOrderEstimateApprove = "/v1/service-orders/:id/estimate/approve"
	PostServiceOrderEstimateReject  = "/v1/service-orders/:id/estimate/reject"
	PostServiceOrderEstimateCancel  = "/v1/service-orders/:id/estimate/cancel"
	PostServiceOrderExecutionCreate = "/v1/service-orders/:id/execution/create"
	PostServiceOrderExecutionFinish = "/v1/service-orders/:id/execution/finish"
	PostServiceOrderPayment         = "/v1/service-orders/:id/payment"
	PostServiceOrderDelivery        = "/v1/service-orders/:id/delivery"

	GetServiceOrder        = "/v1/service-orders/:id"
	GetServiceOrderHistory = "/v1/service-orders/:id/history"
)
