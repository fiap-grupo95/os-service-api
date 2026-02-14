package entities

type PartsSupply struct {
	ID                uint               `json:"id"`
	Price             float64            `json:"price"`
	Quantity          int                `json:"quantity"`
}
