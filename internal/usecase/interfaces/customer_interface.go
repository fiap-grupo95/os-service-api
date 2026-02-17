package interfaces

import (
	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
)

type ICustomerGateway interface {
	GetByID(id string) (*entities.Customer, error)
	GetByDocument(CpfCnpj string) (*entities.Customer, error)
}

type ICustomerRepository interface {
	GetByID(id string) (*response.CustomerResponse, error)
	GetByDocument(CpfCnpj string) (*response.CustomerResponse, error)
}