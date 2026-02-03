package usecase

import (
	"errors"
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/interfaces"
	"github.com/fiap-grupo95/os-service-api/pkg/utils/password"
)

var (
	ErrGeneric               = errors.New("unknown error")
	ErrCustomerNotFound      = errors.New("customer not found")
	ErrInvalidDocumentFormat = errors.New("invalid document format")
	ErrCustomerAlreadyExists = errors.New("customer already exists")
	ErrInvalidCustomerID     = errors.New("invalid customer ID")
)

// ICustomerUseCase defines the interface for customers use cases
type ICustomerUseCase interface {
	GetById(id uint) (*entities.Customer, error)
	GetByDocument(CpfCnpj string) (*entities.Customer, error)
	CreateCustomer(customer *entities.Customer) (*entities.Customer, error)
	UpdateCustomer(id uint, customer *entities.Customer) error
	ListCustomer() ([]entities.Customer, error)
}
type CustomerUseCase struct {
	customerRepo interfaces.ICustomerGateway
	userRepo     interfaces.IUserRepository
}

func NewCustomerUseCase(customerRepo interfaces.ICustomerGateway, userRepo interfaces.IUserRepository) ICustomerUseCase {
	return &CustomerUseCase{customerRepo: customerRepo, userRepo: userRepo}
}

func (uc *CustomerUseCase) GetById(id uint) (*entities.Customer, error) {
	customerDTO, err := uc.customerRepo.GetByID(id)

	if err != nil {
		return nil, ErrGeneric
	}
	if customerDTO == nil {
		return nil, ErrCustomerNotFound
	}

	return customerDTO, nil
}

func (uc *CustomerUseCase) GetByDocument(CpfCnpj string) (*entities.Customer, error) {
	customerDTO, err := uc.customerRepo.GetByDocument(CpfCnpj)

	if err != nil {
		return nil, ErrGeneric
	}
	if customerDTO == nil {
		return nil, ErrCustomerNotFound
	}

	return customerDTO, nil
}

func (uc *CustomerUseCase) CreateCustomer(customer *entities.Customer) (*entities.Customer, error) {
	if e := customer.CpfCnpj.IsValid(); e != nil {
		return nil, ErrInvalidDocumentFormat
	}
	password, _ := pkg.HashPassword(customer.CpfCnpj.String())

	customer.User = &entities.User{
		Email:    customer.Email,
		UserType: valueobject.ParseUserType("customer"),
		Password: password,
	}
	if err := uc.customerRepo.Create(customer); err != nil {
		return nil, err
	}
	return customer, nil
}

func (uc *CustomerUseCase) UpdateCustomer(id uint, customer *entities.Customer) error {
	existingDTO, err := uc.customerRepo.GetByID(id)
	if err != nil {
		return err
	}
	if customer.FullName != "" {
		existingDTO.FullName = customer.FullName
	}

	if customer.PhoneNumber != "" {
		existingDTO.PhoneNumber = customer.PhoneNumber
	}

	return uc.customerRepo.Update(existingDTO)
}

func (uc *CustomerUseCase) ListCustomer() ([]entities.Customer, error) {
	customers, err := uc.customerRepo.List()
	if err != nil {
		return nil, err
	}
	return customers, nil
}
