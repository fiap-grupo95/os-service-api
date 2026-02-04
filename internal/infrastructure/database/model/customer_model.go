package dto

import (
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
)

// 1:1 relationship between Customer and User
// 1:N relationship between Customer and Vehicle
// 1:N relationship between Customer and ServiceOrder
type CustomerModel struct {
	ID            uint                `gorm:"primaryKey"`
	UserID        uint                `gorm:"unique;not null"`
	User          *UserModel          `gorm:"foreignKey:UserID;references:ID"`
	CpfCnpj       string              `gorm:"size:20;not null"`
	PhoneNumber   string              `gorm:"size:20;not null"`
	FullName      string              `gorm:"column:fullname;size:100;not null"`
}

func (cm *CustomerModel) TableName() string {
	return "customer"
}

func (cm *CustomerModel) ToDomain() *entities.Customer {
	var user *entities.User
	if cm.User != nil {
		user = cm.User.ToDomain()
	}
	var email string
	if user != nil {
		email = user.Email
	}
	return &entities.Customer{
		ID:            cm.ID,
		UserID:        cm.UserID,
		Email:         email,
		CpfCnpj:       valueobject.CpfCnpj(cm.CpfCnpj),
		PhoneNumber:   cm.PhoneNumber,
		FullName:      cm.FullName,
	}
}

func FromDomainCustomer(c *entities.Customer) *CustomerModel {
	var userModel *UserModel
	if c.User != nil {
		userModel = &UserModel{
			ID:       c.User.ID,
			Email:    c.User.Email,
			Password: c.User.Password,
			UserType: c.User.UserType,
		}
	}

	return &CustomerModel{
		ID:          c.ID,
		UserID:      c.UserID,
		User:        userModel,
		CpfCnpj:     string(c.CpfCnpj),
		PhoneNumber: c.PhoneNumber,
		FullName:    c.FullName,
	}
}
