package entities

type PartsSupply struct {
	ID                uint               `json:"id"`
	Name              string             `json:"name"`
	Description       string             `json:"description"`
	Price             float64            `json:"price"`
	QuantityTotal     int                `json:"quantity_total"`
	QuantityReserve   int                `json:"quantity_reserve"`
}
