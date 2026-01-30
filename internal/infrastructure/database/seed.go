package database

import (
	"fmt"
	dto "mecanica_xpto/internal/infrastructure/database/model"

	"gorm.io/gorm"
)

const (
	defaultPassword = "jKHN1SmGKuGKyiXhbnaOZg==.0/rdilUJyR5raIXVdOCaX8szZCEUzIpIhYTQIMaLwc8="
	defaultUserType = "admin"
)

func Seed(db *gorm.DB) {
	// Seed user_dtos
	var countUsers int64
	db.Model(&dto.UserModel{}).Count(&countUsers)
	if countUsers == 0 {
		users := []dto.UserModel{
			{
				Email:    "admin@xpto.com",
				Password: defaultPassword,
				UserType: defaultUserType,
			},
			{
				Email:    "joao@xpto.com",
				Password: defaultPassword,
				UserType: defaultUserType,
			},
			{
				Email:    "joana@xpto.com",
				Password: defaultPassword,
				UserType: defaultUserType,
			},
		}

		if err := db.Create(&users).Error; err != nil {
			fmt.Println("Erro ao criar usuários:", err)
			return
		}
		fmt.Println("Seeded users successfully")
	} else {
		fmt.Println("Users already seeded")
	}

	// Seed Service Order Status
	var countSOStatus int64
	db.Model(&dto.ServiceOrderStatus{}).Count(&countSOStatus)
	if countSOStatus == 0 {
		serviceOrderStatus := []dto.ServiceOrderStatus{
			{
				Description: "RECEBIDA",
			},
			{
				Description: "EM DIAGNÓSTICO",
			},
			{
				Description: "AGUARDANDO APROVAÇÃO",
			},
			{
				Description: "APROVADA",
			},
			{
				Description: "REJEITADA",
			},
			{
				Description: "EM EXECUÇÃO",
			},
			{
				Description: "FINALIZADA",
			},
			{
				Description: "ENTREGUE",
			},
			{
				Description: "CANCELADA",
			},
		}

		if err := db.Create(&serviceOrderStatus).Error; err != nil {
			fmt.Println("Error creating Service Order Status:", err)
			return
		}
		fmt.Println("Seeded Service Order Status successfully")
	} else {
		fmt.Println("Service Order Status already seeded")
	}
	// Seed status_Additional_Repair_dtos
	var countARStatus int64
	db.Model(&dto.AdditionalRepairStatusModel{}).Count(&countARStatus)
	if countARStatus == 0 {
		additionalRepairStatus := []dto.AdditionalRepairStatusModel{
			{
				Description: "ABERTA",
			},
			{
				Description: "IN_ANALYSIS",
			},
			{
				Description: "APPROVED",
			},
			{
				Description: "DENIED",
			},
		}

		if err := db.Create(&additionalRepairStatus).Error; err != nil {
			fmt.Println("Error creating Additional Repair Status:", err)
			return
		}
		fmt.Println("Seeded Additional Repair Status successfully")
	} else {
		fmt.Println("Additional Repair Status already seeded")
	}
}
