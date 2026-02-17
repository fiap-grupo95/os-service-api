package dto

import (
	"time"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
)

// N:N relationship between PartsSupply and ServiceOrder
type PartsSupplyServiceOrder struct {
	PartsSupplyID  string `gorm:"column:parts_supply_id;primaryKey"`
	ServiceOrderID string `gorm:"column:service_order_id;primaryKey"`
	Quantity       int  `gorm:"column:quantity"`
}

// N:N relationship between Service and ServiceOrder
type ServiceServiceOrder struct {
	ServiceID      string `gorm:"primaryKey"`
	ServiceOrderID string `gorm:"primaryKey"`
}

type ServiceOrderStatus struct {
	ID          string   `gorm:"primaryKey"`
	Description string `gorm:"size:50;not null"`
}

func (m *ServiceOrderStatus) TableName() string {
	return "service_order_status"
}

func (m *ServiceOrderStatus) ToDomain() valueobject.ServiceOrderStatus {
	return valueobject.ParseServiceOrderStatus(m.Description)
}

type ServiceModel struct {
	ID          string   `gorm:"primaryKey"`
	Name        string `gorm:"size:50;not null"`
	Description string `gorm:"size:255;not null"`
	Price       float64 `gorm:"type:decimal(10,2);not null"`
}

func (s ServiceModel) ToDomain() entities.Service{
	return entities.Service{
		ID: s.ID,
		Name: s.Name,
		Description: s.Description,
		Price: s.Price,
	}
}

type PartsSupplyModel struct {
	ID          string   `gorm:"primaryKey"`
	Price       float64 `gorm:"type:decimal(10,2);not null"`
	Quantity    int    `gorm:"column:quantity"`
}


func (s PartsSupplyModel) ToDomain() entities.PartsSupply{
	return entities.PartsSupply{
		ID: s.ID,
		Price: s.Price,
		Quantity: s.Quantity,
	}
}

type ServiceOrderModel struct {
	ID                       string               `gorm:"primaryKey"`
	CustomerID               string               `gorm:"not null"`
	VehicleID                string               `gorm:"not null"`
	ServiceOrderStatus       ServiceOrderStatus `gorm:"foreignKey:ServiceOrderStatusID"`
	Estimate                 float64            `gorm:"type:decimal(10,2)"`
	StartedExecutionDate     *time.Time
	FinalExecutionDate       *time.Time
	ExecutionDurationInHours float64
	CreatedAt                *time.Time              `gorm:"autoCreateTime"`
	UpdatedAt                *time.Time              `gorm:"autoUpdateTime"`
	AdditionalRepairs        []AdditionalRepairModel `gorm:"foreignKey:ServiceOrderID"`
	PartsSupplies            []PartsSupplyModel      `gorm:"many2many:parts_supply_service_orders;foreignKey:ID;joinForeignKey:ServiceOrderID;references:ID;joinReferences:PartsSupplyID"`
	Services                 []ServiceModel          `gorm:"many2many:service_service_orders;foreignKey:ID;joinForeignKey:ServiceOrderID;references:ID;joinReferences:ServiceID"`
}

func (m *ServiceOrderModel) TableName() string {
	return "service_order"
}

func (m *ServiceOrderModel) ToDomain() *entities.ServiceOrder {
	var additionalRepairs []entities.AdditionalRepair
	var partsSupplies []entities.PartsSupply
	var services []entities.Service

	// Convert AdditionalRepairs
	for _, ar := range m.AdditionalRepairs {
		additionalRepairs = append(additionalRepairs, ar.ToDomain())
	}

	// Convert PartsSupplies
	for _, ps := range m.PartsSupplies {
		partsSupplies = append(partsSupplies, ps.ToDomain())
	}

	// Convert Services
	for _, s := range m.Services {
		services = append(services, s.ToDomain())
	}

	return &entities.ServiceOrder{
		ID:                       m.ID,
		CustomerID:               m.CustomerID,
		VehicleID:                m.VehicleID,
		Status:                   m.ServiceOrderStatus.ToDomain(),
		CreatedAt:                m.CreatedAt,
		UpdatedAt:                m.UpdatedAt,
		AdditionalRepairs:        additionalRepairs,
		PartsSupplies:            partsSupplies,
		Services:                 services,
	}
}
