package interfaces

import (
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
)

type IUserRepository interface {
	GetByID(id string) (*entities.User, error)
	GetByEmail(email string) (*entities.User, error)
	Create(User *entities.User) error
	Update(User *entities.User) error
	Delete(id string) error
	List() ([]entities.User, error)
}
