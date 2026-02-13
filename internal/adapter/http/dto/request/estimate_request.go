package request

type EstimateRequest struct {
	ID             string               `json:"id"`
	ServiceOrderID uint                 `json:"service_order_id"`
	Services       []ServiceRequest     `json:"services"`
	PartsSupplies  []PartsSupplyRequest `json:"parts_supplies"`
}

type PartsSupplyRequest struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required"`
}

type ServiceRequest struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
}
