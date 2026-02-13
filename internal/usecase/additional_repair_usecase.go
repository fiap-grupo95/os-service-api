package usecase

import (
	"context"
	"errors"

	"github.com/fiap-grupo95/os-service-api/internal/usecase/interfaces"

	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/logs"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
)

var (
	ErrAdditionalRepairNotFound = errors.New("additional repair not found")
	ErrStatusNotPermitted       = errors.New("additional repair status not permitted")
)

type IAdditionalRepairUseCase interface {
	CreateAdditionalRepair(ctx context.Context, adr entities.AdditionalRepair) (entities.AdditionalRepair, error)
	AddPartSupplyAndService(ctx context.Context, adrId uint, adr entities.AdditionalRepair) error
	RemovePartSupplyAndService(ctx context.Context, adrId uint, adr entities.AdditionalRepair) error
	GetAdditionalRepair(ctx context.Context, additionalRepairId uint) (entities.AdditionalRepair, error)
	CustomerApprovalStatus(ctx context.Context, additionalRepairId uint, status entities.AdditionalRepairStatusDTO) error
}

type AdditionalRepairUseCase struct {
	repo            interfaces.IAdditionalRepairRepository
	repoOS          interfaces.IServiceOrderGateway
	serviceRepo     interfaces.IServiceGateway
	partsSupplyRepo interfaces.IPartsSupplyGateway
}

var _ IAdditionalRepairUseCase = (*AdditionalRepairUseCase)(nil)

func NewSOAdditionalRepairUseCase(repo interfaces.IAdditionalRepairRepository, repoOS interfaces.IServiceOrderGateway, serviceRepo interfaces.IServiceGateway, partsSupplyRepo interfaces.IPartsSupplyGateway) *AdditionalRepairUseCase {
	return &AdditionalRepairUseCase{
		repo:            repo,
		repoOS:          repoOS,
		serviceRepo:     serviceRepo,
		partsSupplyRepo: partsSupplyRepo,
	}
}

func (u *AdditionalRepairUseCase) CreateAdditionalRepair(ctx context.Context, adr entities.AdditionalRepair) (entities.AdditionalRepair, error) {
	logger := logs.LoggerWithContext(ctx)

	_, err := u.repoOS.GetByID(ctx, adr.ServiceOrderID, false)
	if err != nil {
		logger.Error().Err(err).Any("service_order_id", adr.ServiceOrderID).Msg("error finding service order with id")
		return entities.AdditionalRepair{}, err
	}

	services, estimatedService, err := u.addServiceToAdditionalRepair(ctx, adr.Services)
	if err != nil {
		return entities.AdditionalRepair{}, err
	}

	partsSupplies, estimatedPartsSupply, err := u.addPartsSupplyToAdditionalRepair(ctx, adr.PartsSupplies)
	if err != nil {
		return entities.AdditionalRepair{}, err
	}

	newAdditionalRepair := entities.AdditionalRepair{
		ServiceOrderID: adr.ServiceOrderID,
		Description:    adr.Description,
		ARStatus:       valueobject.AdditionalRepairStatus("IN_ANALYSIS"),
		Estimate:       estimatedService + estimatedPartsSupply,
		Services:       services,
		PartsSupplies:  partsSupplies,
	}

	created, err := u.repo.Create(ctx, newAdditionalRepair)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating additional repair")
		return entities.AdditionalRepair{}, err
	}
	return created, nil
}

func (u *AdditionalRepairUseCase) AddPartSupplyAndService(ctx context.Context, additionalRepairId uint, adr entities.AdditionalRepair) error {
	logger := logs.LoggerWithContext(ctx)

	additionalRepair, err := u.repo.GetByID(ctx, additionalRepairId)
	if err != nil {
		logger.Error().Err(err).Any("additional_repair_id", additionalRepairId).Msg("error finding additional repair with id")
		return err
	}
	if additionalRepair.ID == 0 {
		return ErrAdditionalRepairNotFound
	}
	if err := u.ValidateAdditionalRepairStatus(additionalRepair.ARStatus.String()); err != nil {
		return err
	}

	services, estimatedService, err := u.addServiceToAdditionalRepair(ctx, adr.Services)
	if err != nil {
		return err
	}

	partsSupplies, estimatedPartsSupply, err := u.addPartsSupplyToAdditionalRepair(ctx, adr.PartsSupplies)
	if err != nil {
		return err
	}

	newEstimate := additionalRepair.Estimate + estimatedService + estimatedPartsSupply
	if err := u.repo.AddPartSupplyAndService(ctx, additionalRepairId, services, partsSupplies, newEstimate); err != nil {
		logger.Error().Err(err).Msg("Error adding part supply and services for additional repair")
		return err
	}
	return nil
}

func (u *AdditionalRepairUseCase) RemovePartSupplyAndService(ctx context.Context, additionalRepairId uint, adr entities.AdditionalRepair) error {
	logger := logs.LoggerWithContext(ctx)

	additionalRepair, err := u.repo.GetByID(ctx, additionalRepairId)
	if err != nil {
		logger.Error().Err(err).Any("additional_repair_id", additionalRepairId).Msg("error finding additional repair with id")
		return err
	}
	if additionalRepair.ID == 0 {
		return ErrAdditionalRepairNotFound
	}
	if err := u.ValidateAdditionalRepairStatus(additionalRepair.ARStatus.String()); err != nil {
		return err
	}

	services, estimatedService, err := u.addServiceToAdditionalRepair(ctx, adr.Services)
	if err != nil {
		return err
	}
	partsSupplies, estimatedPartsSupply, err := u.addPartsSupplyToAdditionalRepair(ctx, adr.PartsSupplies)
	if err != nil {
		return err
	}

	newEstimate := estimatedService + estimatedPartsSupply
	if err := u.repo.ReplacePartSupplyAndService(ctx, additionalRepairId, services, partsSupplies, newEstimate); err != nil {
		logger.Error().Err(err).Msg("Error updating parts supply and services for additional repair")
		return err
	}
	return nil
}

func (u *AdditionalRepairUseCase) GetAdditionalRepair(ctx context.Context, additionalRepairId uint) (entities.AdditionalRepair, error) {
	logger := logs.LoggerWithContext(ctx)

	additionalRepair, err := u.repo.GetByID(ctx, additionalRepairId)
	if err != nil {
		logger.Error().Err(err).Any("additional_repair_id", additionalRepairId).Msg("error finding additional repair with id")
		return entities.AdditionalRepair{}, err
	}
	if additionalRepair.ID == 0 {
		logger.Error().Msg("additional repair not found")
		return entities.AdditionalRepair{}, ErrAdditionalRepairNotFound
	}
	return additionalRepair, nil
}

func (u *AdditionalRepairUseCase) CustomerApprovalStatus(ctx context.Context, additionalRepairId uint, status entities.AdditionalRepairStatusDTO) error {
	logger := logs.LoggerWithContext(ctx)

	additionalRepair, err := u.repo.GetByID(ctx, additionalRepairId)
	if err != nil {
		logger.Error().Err(err).Any("additional_repair_id", additionalRepairId).Msg("error finding additional repair with id")
		return err
	}
	if additionalRepair.ID == 0 {
		return ErrAdditionalRepairNotFound
	}
	if err := u.ValidateAdditionalRepairStatus(additionalRepair.ARStatus.String()); err != nil {
		return err
	}

	if status.ApprovalStatus == "DENIED" {
		for _, ps := range additionalRepair.PartsSupplies {
			if err := unreservePartsSupply(ctx, ps, u.partsSupplyRepo); err != nil {
				logger.Error().Err(err).Msg("Error unreserving parts supply")
			}
		}
		logger.Info().Any("additional_repair_id", additionalRepairId).Msg("Customer rejected additional repair")
	} else {
		for _, ps := range additionalRepair.PartsSupplies {
			quantity, err := u.repo.GetPartsSupplyQuantity(ctx, ps.ID, additionalRepair.ID)
			if err != nil {
				logger.Error().Err(err).Msg("Error getting parts supply additional repair relation")
				return err
			}

			entity := entities.PartsSupply{
				ID:              ps.ID,
				QuantityReserve: quantity,
				QuantityTotal:   quantity,
			}
			if err := releaseReservedPartsSupply(ctx, entity, u.partsSupplyRepo); err != nil {
				logger.Error().Err(err).Msg("Error releasing reserved parts supply")
				return err
			}
		}
	}

	if err := u.repo.UpdateStatus(ctx, additionalRepairId, status); err != nil {
		logger.Error().Err(err).Any("additional_repair_id", additionalRepairId).Msg("error updating customer approval with id")
		return err
	}

	if err := u.repoOS.UpdateEstimate(ctx, additionalRepair.ServiceOrderID, additionalRepair.Estimate); err != nil {
		logger.Error().Err(err).Any("service_order_id", additionalRepair.ServiceOrderID).Msg("error updating service order estimate with id")
		return err
	}

	return nil
}

func (u *AdditionalRepairUseCase) addPartsSupplyToAdditionalRepair(ctx context.Context, partsSupplyId []entities.PartsSupply) ([]entities.PartsSupply, float64, error) {
	var listPartsSupply []entities.PartsSupply
	var estimatedPrice float64

	for _, ps := range partsSupplyId {
		err := reservePartsSupply(ctx, ps, u.partsSupplyRepo)
		if err != nil {
			logs.Logger().Error().Err(err).Msg("Error reserving parts supply")
			return nil, 0, err
		}
		psDto, err := u.partsSupplyRepo.GetByID(ctx, ps.ID)
		if err != nil {
			logs.Logger().Error().Err(err).Msg("error finding parts supply with id")
		}
		if psDto.ID == 0 {
			logs.Logger().Error().Msg("parts supply with id not found")
			return listPartsSupply, estimatedPrice, ErrServiceNotFound
		}
		psDto.QuantityReserve = ps.QuantityReserve
		if ps.QuantityTotal > 0 {
			psDto.QuantityTotal = ps.QuantityTotal
		}
		estimatedPrice += psDto.Price * float64(psDto.QuantityReserve)
		listPartsSupply = append(listPartsSupply, *psDto)
	}
	return listPartsSupply, estimatedPrice, nil
}

func (u *AdditionalRepairUseCase) addServiceToAdditionalRepair(ctx context.Context, services []entities.Service) ([]entities.Service, float64, error) {
	var listServices []entities.Service
	var estimatedPrice float64

	for _, s := range services {
		serviceDto, err := u.serviceRepo.GetByID(ctx, s.ID)
		if err != nil {
			logs.Logger().Error().Err(err).Msg("error finding service with id")
		}
		if serviceDto.ID == 0 {
			logs.Logger().Error().Msg("service with id not found")
			return listServices, estimatedPrice, ErrServiceNotFound
		}
		estimatedPrice += serviceDto.Price
		listServices = append(listServices, *serviceDto)
	}
	return listServices, estimatedPrice, nil
}

func (u *AdditionalRepairUseCase) ValidateAdditionalRepairStatus(status string) error {
	if status != "IN_ANALYSIS" {
		logs.Logger().Error().Msg("invalid additional repair status")
		return ErrStatusNotPermitted
	}
	return nil
}

func reservePartsSupply(ctx context.Context, partsSupply entities.PartsSupply, repo interfaces.IPartsSupplyGateway) error {
	return repo.Reserve(ctx, []entities.PartsSupply{partsSupply})
}

func unreservePartsSupply(ctx context.Context, partsSupply entities.PartsSupply, repo interfaces.IPartsSupplyGateway) error {
	return repo.Release(ctx, []entities.PartsSupply{partsSupply})
}

func releaseReservedPartsSupply(ctx context.Context, partsSupply entities.PartsSupply, repo interfaces.IPartsSupplyGateway) error {
	return repo.Release(ctx, []entities.PartsSupply{partsSupply})
}
