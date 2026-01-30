package entities

import (
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
	"time"
)

type User struct {
	ID        uint
	Email     string
	Password  string
	UserType  valueobject.UserType
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Customer  *Customer
}
