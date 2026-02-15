package entity_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/request"
	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/logs"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/interfaces"
)

// CustomerRepository implements ICustomerRepository interface
type CustomerRepository struct {
	http *http.Client
}

func NewCustomerRepository() interfaces.ICustomerRepository {
	return &CustomerRepository{http: http.DefaultClient}
}

func (r *CustomerRepository) Create(customer *request.CustomerCreateRequest) error {
	logger := logs.Logger()
	path := fmt.Sprintf(CUSTOMERS_ENDPOINT)

	payload, err := json.Marshal(customer)
	if err != nil {
		logger.Error().Err(err).Msg("failed to marshal customer")
		return err
	}

	body := io.NopCloser(bytes.NewReader(payload))
	resp, err := r.http.Post(path, "application/json", body)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create customer")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		logger.Error().Err(err).Msg("failed to create customer")
		return fmt.Errorf("failed to create customer: status %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var createdCustomer response.CustomerResponse
	if err := json.NewDecoder(resp.Body).Decode(&createdCustomer); err != nil {
		logger.Error().Err(err).Msg("failed to decode customer")
		return err
	}

	return nil
}

func (r *CustomerRepository) GetByID(id uint) (*response.CustomerResponse, error) {
	logger := logs.Logger()
	path := fmt.Sprintf(CUSTOMERS_ID_ENDPOINT, id)
	resp, err := r.http.Get(path)
	if err != nil {
		logger.Error().Err(err).Msg("failed to find customer by ID")
		return nil, err
	}
	defer resp.Body.Close()

	var customer response.CustomerResponse
	if err := json.NewDecoder(resp.Body).Decode(&customer); err != nil {
		logger.Error().Err(err).Msg("failed to decode customer")
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepository) GetByDocument(CpfCnpj string) (*response.CustomerResponse, error) {
	logger := logs.Logger()
	path := fmt.Sprintf(CUSTOMERS_DOCUMENT_ENDPOINT, CpfCnpj)
	resp, err := r.http.Get(path)
	if err != nil {
		logger.Error().Err(err).Msg("failed to find customer by document")
		return nil, err
	}
	defer resp.Body.Close()

	var customer response.CustomerResponse
	if err := json.NewDecoder(resp.Body).Decode(&customer); err != nil {
		logger.Error().Err(err).Msg("failed to decode customer")
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepository) Update(customer *request.CustomerUpdateRequest, id uint) error {
logger := logs.Logger()
	path := fmt.Sprintf(CUSTOMERS_ID_ENDPOINT, id)

	payload, err := json.Marshal(customer)
	if err != nil {
		logger.Error().Err(err).Msg("failed to marshal customer")
		return err	
	}

	req, err := http.NewRequest(http.MethodPatch, path, bytes.NewReader(payload))
	if err != nil {
		logger.Error().Err(err).Msg("failed to create request")
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := r.http.Do(req)
	if err != nil {
		logger.Error().Err(err).Msg("failed to send request")
		return err
	}	
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		logger.Error().Err(err).Msg("failed to update customer")
		return fmt.Errorf("failed to update customer: status %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func (r *CustomerRepository) List() ([]response.CustomerResponse, error) {
	logger := logs.Logger()
	resp, err := http.Get(CUSTOMERS_ENDPOINT)
	if err != nil {
		logger.Error().Err(err).Msg("failed to find all customers")
		return nil, err
	}
	defer resp.Body.Close()

	var customers []response.CustomerResponse
	if err := json.NewDecoder(resp.Body).Decode(&customers); err != nil {
		logger.Error().Err(err).Msg("failed to decode customers")
		return nil, err
	}
	return customers, nil
}
