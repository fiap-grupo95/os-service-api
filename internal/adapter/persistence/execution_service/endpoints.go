package execution_service

const (
	BASE_URL = "http://mockoon:8083"
)

// Execution Service Endpoints
const (
	EXECUTION_SERVICE_ENDPOINT    = BASE_URL + "/v1/execution"
	EXECUTION_SERVICE_FINISH_ENDPOINT = BASE_URL + "/v1/execution/finish/%d"
	EXECUTION_SERVICE_ID_ENDPOINT = BASE_URL + "/v1/execution/%d"
)