package entity_api

const (
	BASE_URL = "http://localhost:8083"
)

const (
	VEHICLES_ENDPOINT = BASE_URL + "/vehicles"
	VEHICLES_ID_ENDPOINT = BASE_URL + "/vehicles/%d"
	VEHICLES_PLATE_ENDPOINT = BASE_URL + "/vehicles/plate/%s"
	VEHICLES_CUSTOMER_ID_ENDPOINT = BASE_URL + "/vehicles/customer/%d"
)

const (
	CUSTOMERS_ENDPOINT = BASE_URL + "/customers"
	CUSTOMERS_ID_ENDPOINT = BASE_URL + "/customers/%d"
	CUSTOMERS_DOCUMENT_ENDPOINT = BASE_URL + "/customers/document/%s"
)