package dto

import (
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"time"

	"gorm.io/gorm"
)

// N:N relationship between PartsSupply and ServiceOrder
// 1:N relationship between PartsSupply and AdditionalRepair
type PartsSupplyModel struct {
	ID                uint                    `gorm:"primaryKey"`
	Name              string                  `gorm:"size:100;not null"`
	Description       string                  `gorm:"type:text"`
	Price             float64                 `gorm:"type:decimal(10,2);not null"`
	QuantityTotal     int                     `gorm:"not null;default:0"`
	QuantityReserve   int                     `gorm:"not null;default:0"`
	CreatedAt         time.Time               `gorm:"autoCreateTime"`
	UpdatedAt         time.Time               `gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt          `gorm:"index"`
	AdditionalRepairs []AdditionalRepairModel `gorm:"many2many:parts_supply_additional_repairs;joinForeignKey:parts_supply_id;joinReferences:additional_repair_id"`
	ServiceOrders     []ServiceOrderModel     `gorm:"many2many:parts_supply_service_orders;joinForeignKey:parts_supply_id;joinReferences:service_order_id"`
}

func (m *PartsSupplyModel) TableName() string {
	return "parts_supply"
}

func (m *PartsSupplyModel) ToDomain() entities.PartsSupply {
	return entities.PartsSupply{
		ID:              m.ID,
		Price:           m.Price,
	}
}
