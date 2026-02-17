package entity_api

const (
	BASE_URL = "http://mockoon:8083"
)

// Vehicle Endpoints
const (
	VEHICLES_ENDPOINT             = BASE_URL + "/v1/vehicles"
	VEHICLES_ID_ENDPOINT          = BASE_URL + "/v1/vehicles/%s"
	VEHICLES_PLATE_ENDPOINT       = BASE_URL + "/v1/vehicles/plate/%s"
	VEHICLES_CUSTOMER_ID_ENDPOINT = BASE_URL + "/v1/vehicles/customer/%s"
)

// Customer Endpoints
const (
	CUSTOMERS_ENDPOINT          = BASE_URL + "/v1/customers"
	CUSTOMERS_ID_ENDPOINT       = BASE_URL + "/v1/customers/%s"
	CUSTOMERS_DOCUMENT_ENDPOINT = BASE_URL + "/v1/customers/document/%s"
)

// Parts Supply Endpoints
const (
	PARTS_SUPPLY_ID_ENDPOINT                = BASE_URL + "/v1/parts-supply/%s"
	PARTS_SUPPLY_SERVICE_ORDER_ID_ENDPOINT  = BASE_URL + "/v1/parts-supply/service-order/%s"
	PARTS_SUPPLY_RESERVE_ENDPOINT           = BASE_URL + "/v1/parts-supply/reserve"
	PARTS_SUPPLY_RELEASE_ENDPOINT           = BASE_URL + "/v1/parts-supply/release"
	PARTS_SUPPLY_WRITEOFF_ENDPOINT          = BASE_URL + "/v1/parts-supply/writeoff"
	PARTS_SUPPLY_AUTHORIZE_RESERVE_ENDPOINT = BASE_URL + "/v1/parts-supply/authorize-reserve"
)

// Service Endpoints
const (
	SERVICE_ENDPOINT    = BASE_URL + "/v1/service"
	SERVICE_ID_ENDPOINT = BASE_URL + "/v1/service/%s"
)
