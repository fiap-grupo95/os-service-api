package interfaces

import (
	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/request"
	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
)

type ICustomerGateway interface {
	GetByID(id uint) (*entities.Customer, error)
	GetByDocument(CpfCnpj string) (*entities.Customer, error)
	Create(customer *entities.Customer) error
	Update(customer *entities.Customer) error
	List() ([]entities.Customer, error)
}

type ICustomerRepository interface {
	GetByID(id uint) (*response.CustomerResponse, error)
	GetByDocument(CpfCnpj string) (*response.CustomerResponse, error)
	Create(customer *request.CustomerCreateRequest) error
	Update(customer *request.CustomerUpdateRequest, id uint) error
	List() ([]response.CustomerResponse, error)
}