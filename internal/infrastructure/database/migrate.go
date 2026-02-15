package database

import (
	"fmt"
	Model "github.com/fiap-grupo95/os-service-api/internal/infrastructure/database/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&Model.ServiceOrderModel{},
		&Model.UserModel{},
		&Model.ServiceOrderStatus{},
		&Model.AdditionalRepairModel{},
		&Model.AdditionalRepairStatusModel{},
		&Model.UserTypeModel{},
	)
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	fmt.Println("Database migrated successfully")
}
