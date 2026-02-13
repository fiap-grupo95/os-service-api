package entities

type Estimate struct {
	ID     string   `json:"id"`
	Value  *float64 `json:"value"`
	Status string  `json:"status"`
}