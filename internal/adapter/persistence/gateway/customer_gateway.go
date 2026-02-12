package gateway

import (
	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/interfaces"
)

type CustomerGateway struct {
	repo interfaces.ICustomerRepository
}

func NewCustomerGateway(repo interfaces.ICustomerRepository) interfaces.ICustomerGateway {
	return &CustomerGateway{repo: repo}
}

func (g *CustomerGateway) GetByID(id uint) (*entities.Customer, error) {
	response, err := g.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	customer := mapSearchReponseToDomain(response)
	return customer, nil
}

func (g *CustomerGateway) GetByDocument(CpfCnpj string) (*entities.Customer, error) {
	response, err := g.repo.GetByDocument(CpfCnpj)
	if err != nil {
		return nil, err
	}

	customer := mapSearchReponseToDomain(response)
	return customer, nil
}

func mapSearchReponseToDomain(customer *response.CustomerResponse) *entities.Customer {
	return &entities.Customer{
		ID:            customer.ID,
		UserID:        customer.UserID,
		FullName:      customer.FullName,
		Email:         customer.Email,
		PhoneNumber:   customer.PhoneNumber,
		CpfCnpj:       valueobject.CpfCnpj(customer.Document),
		User:          nil,
		Vehicles:      nil,
		ServiceOrders: nil,
	}
}
