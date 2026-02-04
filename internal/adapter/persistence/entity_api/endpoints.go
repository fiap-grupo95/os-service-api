package entity_api

const (
	BASE_URL = "http://mockoon:8083"
)

const (
	VEHICLES_ENDPOINT = BASE_URL + "/v1/vehicles"
	VEHICLES_ID_ENDPOINT = BASE_URL + "/v1/vehicles/%d"
	VEHICLES_PLATE_ENDPOINT = BASE_URL + "/v1/vehicles/plate/%s"
	VEHICLES_CUSTOMER_ID_ENDPOINT = BASE_URL + "/v1/vehicles/customer/%d"
)

const (
	CUSTOMERS_ENDPOINT = BASE_URL + "/v1/customers"
	CUSTOMERS_ID_ENDPOINT = BASE_URL + "/v1/customers/%d"
	CUSTOMERS_DOCUMENT_ENDPOINT = BASE_URL + "/v1/customers/document/%s"
)