package usecase

import (
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/interfaces"
	"github.com/fiap-grupo95/os-service-api/pkg/utils/errors"
	"net/http"
)

const (
	ErrCodeInvalidCredential = "INVALID_CREDENTIALS"
	ErrMsgInvalidCredential  = "Invalid email or password"
	ErrCodeTokenGeneration   = "TOKEN_GENERATION_ERROR"
	ErrMsgTokenGeneration    = "Failed to generate token"
)

type AuthInterface interface {
	Login(Email, Password string) (string, *pkg.AppError)
}

type authUseCase struct {
	tokenService interfaces.AuthTokenInterface
	userRepo     interfaces.IUserRepository
}

func NewAuthUseCase(tokenService interfaces.AuthTokenInterface, userRepo interfaces.IUserRepository) *authUseCase {
	return &authUseCase{
		tokenService: tokenService,
		userRepo:     userRepo,
	}
}

// Login handles user login and returns a JWT token
func (a *authUseCase) Login(Email, Password string) (string, *pkg.AppError) {
	userFromDB, err := a.userRepo.GetByEmail(Email)
	if err != nil {
		return "", pkg.NewDomainErrorSimple(ErrCodeInvalidCredential, ErrMsgInvalidCredential, http.StatusUnauthorized)
	}

	hashedPass := valueobject.Password(userFromDB.Password)

	if !hashedPass.Verify(Password) {
		return "", pkg.NewDomainErrorSimple(ErrCodeInvalidCredential, ErrMsgInvalidCredential, http.StatusUnauthorized)
	}

	token, err := a.tokenService.GenerateToken(userFromDB.Email)
	if err != nil {
		return "", pkg.NewInfraError(ErrCodeTokenGeneration, ErrMsgTokenGeneration, err, http.StatusInternalServerError)
	}

	return token, nil
}
