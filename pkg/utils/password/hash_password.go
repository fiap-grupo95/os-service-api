package pkg

import (
	passwordpkg "github.com/fiap-grupo95/os-service-api/internal/domain/valueobject/password_utils"
)

func HashPassword(password string) (string, error) {
	return passwordpkg.HashPassword(password)
}
