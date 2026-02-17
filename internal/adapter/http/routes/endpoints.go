package routes

const (
	PathHealthCheck = "/v1/ping"
	PostLogin       = "/v1/login"
)

const (
	PostServiceOrderCreate          = "/v1/service-orders/create"
	PostServiceOrderCancel          = "/v1/service-orders/:id/cancel"
	PostServiceOrderDiagnosis       = "/v1/service-orders/:id/diagnosis"
	PostServiceOrderEstimateApprove = "/v1/service-orders/:id/estimate/approve"
	PostServiceOrderEstimateReject  = "/v1/service-orders/:id/estimate/reject"
	PostServiceOrderEstimateCancel  = "/v1/service-orders/:id/estimate/cancel"
	PostServiceOrderExecutionCreate = "/v1/service-orders/:id/execution/create"
	PostServiceOrderExecutionFinish = "/v1/service-orders/:id/execution/finish"
	PostServiceOrderPayment         = "/v1/service-orders/:id/payment"
	PostServiceOrderDelivery        = "/v1/service-orders/:id/delivery"

	GetServiceOrder     = "/v1/service-orders/:id"
	GetServiceOrderList = "/v1/service-orders"
)

const (
	PostAdditionalRepair        = "/v1/additional-repair"
	GetAdditionalRepair         = "/v1/additional-repair/:id"
	GetAdditionalRepairBySO     = "/v1/additional-repair/service-orders/:id"
	PostAdditionalRepairApprove = "/v1/additional-repair/:id/approve"
	PostAdditionalRepairReject  = "/v1/additional-repair/:id/reject"
	PostAdditionalRepairCancel  = "/v1/additional-repair/:id/cancel"
)
