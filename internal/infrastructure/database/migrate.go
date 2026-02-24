package database

import (
	"fmt"
	Model "github.com/fiap-grupo95/os-service-api/internal/infrastructure/database/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&Model.PartsSupplyModel{},
		&Model.ServiceModel{},
		&Model.VehicleModel{},
		&Model.ServiceOrderModel{},
		&Model.CustomerModel{},
		&Model.UserModel{},
		&Model.ServiceOrderStatus{},
		&Model.AdditionalRepairModel{},
		&Model.PartsSupplyServiceOrder{},
		&Model.PartsSupplyAdditionalRepair{},
		&Model.AdditionalRepairStatusModel{},
		&Model.UserTypeModel{},
		&Model.ServiceServiceOrder{},
		&Model.PaymentModel{},
	)
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	fmt.Println("Database migrated successfully")
}
