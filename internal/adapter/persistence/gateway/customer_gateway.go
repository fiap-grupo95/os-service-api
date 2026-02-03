package gateway

import (
	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/request"
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
func (g *CustomerGateway) Create(customer *entities.Customer) error {
	request := &request.CustomerCreateRequest{
		FullName:    customer.FullName,
		Email:       customer.Email,
		PhoneNumber: customer.PhoneNumber,
		Document:    customer.CpfCnpj.String(),
	}
	err := g.repo.Create(request)
	if err != nil {
		return err
	}
	return nil
}

func (g *CustomerGateway) Update(customer *entities.Customer) error {
	request := &request.CustomerUpdateRequest{
		FullName:    customer.FullName,
		Email:       customer.Email,
		PhoneNumber: customer.PhoneNumber,
	}
	err := g.repo.Update(request, customer.ID)
	if err != nil {
		return err
	}
	return nil
}

func (g *CustomerGateway) List() ([]entities.Customer, error) {
	response, err := g.repo.List()
	if err != nil {
		return nil, err
	}

	customers := make([]entities.Customer, len(response))
	for _, customer := range response {
		c := mapSearchReponseToDomain(&customer)
		customers = append(customers, *c)
	}
	return customers, nil
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
