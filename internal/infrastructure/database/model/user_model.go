package dto

import (
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
	"time"

	"gorm.io/gorm"
)

type UserModel struct {
	ID        string                 `gorm:"primaryKey"`
	Email     string               `gorm:"size:100;not null;unique" json:"email"`
	Password  string               `gorm:"size:255;not null" json:"password"`
	UserType  valueobject.UserType `gorm:"not null"`
	CreatedAt time.Time            `gorm:"autoCreateTime"`
	UpdatedAt time.Time            `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt       `gorm:"index"`
}

func (m *UserModel) TableName() string {
	return "user"
}

func (m *UserModel) ToDomain() *entities.User {
	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}
	return &entities.User{
		ID:        m.ID,
		Email:     m.Email,
		Password:  m.Password,
		UserType:  m.UserType,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func FromDomainUser(user *entities.User) *UserModel {
	var deletedAt gorm.DeletedAt
	if user.DeletedAt != nil {
		deletedAt = gorm.DeletedAt{Time: *user.DeletedAt, Valid: true}
	}
	return &UserModel{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		UserType:  user.UserType,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: deletedAt,
	}
}
