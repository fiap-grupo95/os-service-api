package response

type EstimateResponse struct {
	ID             uint     `json:"id"`
	ServiceOrderID uint     `json:"service_order_id"`
	Value          *float64 `json:"value"`
	Status         string   `json:"status"`
}
