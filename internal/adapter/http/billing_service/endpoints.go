package billing_service

const (
	BASE_URL = "http://mockoon:8083"
)

// Billing Service Endpoints
const (
	ESTIMATE_CREATE_ENDPOINT        = BASE_URL + "/v1/estimates"
	ESTIMATE_REJECT_ENDPOINT        = BASE_URL + "/v1/estimates/%s/reject"
	ESTIMATE_APPROVE_ENDPOINT       = BASE_URL + "/v1/estimates/%s/approve"
	ESTIMATE_CANCEL_ENDPOINT        = BASE_URL + "/v1/estimates/%s/cancel"
	PAYMENT_CREATE_ENDPOINT         = BASE_URL + "/v1/payments/%s"
	PAYMENT_BY_ESTIMATE_ID_ENDPOINT = BASE_URL + "/v1/payments/estimate/%s"
)
