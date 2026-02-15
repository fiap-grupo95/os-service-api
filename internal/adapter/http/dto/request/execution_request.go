package request

type ExecutionRequest struct {
	ID             string   `json:"id"`
	ServiceOrderID uint   `json:"service_order_id"`
}