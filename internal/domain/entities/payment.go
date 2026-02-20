package entities

import (
	"time"
)

type Payer struct {
	Email string `json:"email"`
}

type Payment struct {
	ID              string    `json:"id"`
	EstimateID      string    `json:"estimate_id"`
	PaymentDate     time.Time `json:"payment_date"`
	Amount          float64   `json:"amount"`
	PaymentMethodID string    `json:"payment_method_id"`
	Payer           Payer     `json:"payer"`
}
