package gateway

import (
	"context"

	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/logs"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/interfaces"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type PartsSupplyGateway struct {
	repo interfaces.IPartsSupplyRepository
}

func NewPartsSupplyGateway(repo interfaces.IPartsSupplyRepository) *PartsSupplyGateway {
	return &PartsSupplyGateway{repo: repo}
}

func (s *PartsSupplyGateway) GetByID(ctx context.Context, id uint) (*entities.PartsSupply, error) {
	logger := logs.Logger()
	if txn := newrelic.FromContext(ctx); txn != nil {
		logger = logs.LoggerWithContext(ctx)
	}
	reponse, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Error().Err(err).Uint("ID", id).Msg("failed to find parts supply by id")
		return nil, err
	}
	return mapPartsSupplyResponseToDomain(ctx, *reponse), nil
}

func (s *PartsSupplyGateway) GetByServiceOrderID(ctx context.Context, serviceOrderID uint) ([]entities.PartsSupply, error){
	logger := logs.Logger()
	partsSupplies := make([]entities.PartsSupply, 0)
	response, err := s.repo.GetByServiceOrderID(ctx, serviceOrderID)
	if err != nil {
		logger.Error().Err(err).Uint("OS_ID", serviceOrderID).Msg("failed to find parts supply by service order id")
		return nil, err
	}

	for _, r := range response {
		partsSupplies = append(partsSupplies, *mapPartsSupplyResponseToDomain(ctx, r))
	}
	return partsSupplies, nil
}

func (s *PartsSupplyGateway) Reserve(ctx context.Context, partsSupply []entities.PartsSupply) error{
	partsSupplyRequest := mapPartsSupplyDomainToRequest(partsSupply)
	return s.repo.Reserve(ctx, partsSupplyRequest)
}

func (s *PartsSupplyGateway) Release(ctx context.Context, partsSupply []entities.PartsSupply) error{
	partsSupplyRequest := mapPartsSupplyDomainToRequest(partsSupply)
	return s.repo.Release(ctx, partsSupplyRequest)
}

func (s *PartsSupplyGateway) AuthorizeReserve(ctx context.Context, partsSupply []entities.PartsSupply) error{
	partsSupplyRequest := mapPartsSupplyDomainToRequest(partsSupply)
	return s.repo.AuthorizeReserve(ctx, partsSupplyRequest)
}

func (s *PartsSupplyGateway) WriteOff(ctx context.Context, partsSupply []entities.PartsSupply) error{
	partsSupplyRequest := mapPartsSupplyDomainToRequest(partsSupply)
	return s.repo.WriteOff(ctx, partsSupplyRequest)
}

func mapPartsSupplyResponseToDomain(ctx context.Context, response response.PartsSupplyResponse) *entities.PartsSupply{
	return &entities.PartsSupply{
		ID: response.ID,
		Price: response.Price,
		Quantity: response.QuantityTotal,
	}
}