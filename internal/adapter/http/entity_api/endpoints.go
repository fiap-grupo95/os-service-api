package entity_api

const (
	BASE_URL = "http://mockoon:8083"
)

// Vehicle Endpoints
const (
	VEHICLES_ENDPOINT             = BASE_URL + "/v1/vehicles"
	VEHICLES_ID_ENDPOINT          = BASE_URL + "/v1/vehicles/%d"
	VEHICLES_PLATE_ENDPOINT       = BASE_URL + "/v1/vehicles/plate/%s"
	VEHICLES_CUSTOMER_ID_ENDPOINT = BASE_URL + "/v1/vehicles/customer/%d"
)

// Customer Endpoints
const (
	CUSTOMERS_ENDPOINT          = BASE_URL + "/v1/customers"
	CUSTOMERS_ID_ENDPOINT       = BASE_URL + "/v1/customers/%d"
	CUSTOMERS_DOCUMENT_ENDPOINT = BASE_URL + "/v1/customers/document/%s"
)

// Parts Supply Endpoints
const (
	PARTS_SUPPLY_ID_ENDPOINT               = BASE_URL + "/v1/parts-supply/%d"
	PARTS_SUPPLY_SERVICE_ORDER_ID_ENDPOINT = BASE_URL + "/v1/parts-supply/service-order/%d"
)

// Service Endpoints
const (
	SERVICE_ENDPOINT    = BASE_URL + "/v1/service"
	SERVICE_ID_ENDPOINT = BASE_URL + "/v1/service/%d"
)