package entities

import (
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
	"time"
)

type ServiceOrder struct {
	ID                       uint                           `json:"id"`
	CustomerID               uint                           `json:"customer_id"`
	Customer                 *Customer                      `json:"customer"`
	VehicleID                uint                           `json:"vehicle_id"`
	Vehicle                  *Vehicle                       `json:"vehicle"`
	Status                   valueobject.ServiceOrderStatus `json:"service_order_status"`
	Estimate                 *Estimate                      `json:"estimate,omitempty"`
	Execution                *Execution                     `json:"execution,omitempty"`
	CreatedAt                *time.Time                     `json:"created_at,omitempty"`
	UpdatedAt                *time.Time                     `json:"updated_at,omitempty"`
	AdditionalRepairs        []AdditionalRepair             `json:"additional_repairs,omitempty"`
	PartsSupplies            []PartsSupply                  `json:"parts_supplies,omitempty"`
	Services                 []Service                      `json:"services,omitempty"`
}
