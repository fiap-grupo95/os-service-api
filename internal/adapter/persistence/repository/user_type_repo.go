package repository

import (
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/database/model"
)

type UserTypeDTO struct {
	ID    uint            `gorm:"primaryKey"`
	Type  string          `gorm:"size:50;not null"`
	Users []dto.UserModel `gorm:"foreignKey:UserTypeID;references:ID"`
}

func (utm *UserTypeDTO) ToDomain() valueobject.UserType {
	return valueobject.ParseUserType(utm.Type)
}
