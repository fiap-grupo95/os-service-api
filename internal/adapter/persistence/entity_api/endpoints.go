package entity_api

const (
	BASE_URL = "http://localhost:8083"
)

const (
	VEHICLES_ENDPOINT = BASE_URL + "/vehicles"
	VEHICLES_ID_ENDPOINT = BASE_URL + "/vehicles/%s"
	VEHICLES_PLATE_ENDPOINT = BASE_URL + "/vehicles/plate/%s"
	VEHICLES_CUSTOMER_ID_ENDPOINT = BASE_URL + "/vehicles/customer/%s"
)

const (
	CUSTOMERS_ENDPOINT = BASE_URL + "/customers"
	CUSTOMERS_ID_ENDPOINT = BASE_URL + "/customers/%s"
	CUSTOMERS_DOCUMENT_ENDPOINT = BASE_URL + "/customers/document/%s"
)