package entities

import "time"

type Execution struct {
	ID             uint       `json:"id"`
	ServiceOrderID uint       `json:"service_order_id"`
	Status         string     `json:"status"`
	StartedAt      *time.Time `json:"started_at"`
	FinishedAt     *time.Time `json:"finished_at"`
}
