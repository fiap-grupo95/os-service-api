package usecase

import (
	"context"
	"errors"

	"github.com/fiap-grupo95/os-service-api/internal/usecase/interfaces"
	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/logs"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
)

var (
	ErrAdditionalRepairNotFound = errors.New("additional repair not found")
	ErrStatusNotPermitted       = errors.New("additional repair status not permitted")
)

const (
	APPROVED_FLOW = "APPROVED"
	REJECTED_FLOW = "REJECTED"
)

type IAdditionalRepairUseCase interface {
	CreateAdditionalRepair(ctx context.Context, adr entities.AdditionalRepair) (*entities.AdditionalRepair, error)
	GetAdditionalRepair(ctx context.Context, additionalRepairId string) (*entities.AdditionalRepair, error)
	CustomerApprovalStatus(ctx context.Context, additionalRepairId string, flow string) (*entities.AdditionalRepair, error)
	CancelAdditionalRepair(ctx context.Context, additionalRepairId string) (*entities.AdditionalRepair, error)
	Rollback(ctx context.Context, additionalRepair *entities.AdditionalRepair, hasReservedPartsSupply bool) error
}

type AdditionalRepairUseCase struct {
	repo            interfaces.IAdditionalRepairGateway
	repoOS          interfaces.IServiceOrderGateway
	serviceRepo     interfaces.IServiceGateway
	partsSupplyRepo interfaces.IPartsSupplyGateway
	billRepo        interfaces.IBillingServiceGateway
}

var _ IAdditionalRepairUseCase = (*AdditionalRepairUseCase)(nil)

func NewAdditionalRepairUseCase(repo interfaces.IAdditionalRepairGateway, repoOS interfaces.IServiceOrderGateway, serviceRepo interfaces.IServiceGateway, partsSupplyRepo interfaces.IPartsSupplyGateway, billRepo interfaces.IBillingServiceGateway) *AdditionalRepairUseCase {
	return &AdditionalRepairUseCase{
		repo:            repo,
		repoOS:          repoOS,
		serviceRepo:     serviceRepo,
		partsSupplyRepo: partsSupplyRepo,
		billRepo:        billRepo,
	}
}

func (u *AdditionalRepairUseCase) CreateAdditionalRepair(ctx context.Context, adr entities.AdditionalRepair) (*entities.AdditionalRepair, error) {
	logger := logs.LoggerWithContext(ctx)

	_, err := u.repoOS.GetByID(ctx, adr.ServiceOrderID, false)
	if err != nil {
		logger.Error().Err(err).Any("service_order_id", adr.ServiceOrderID).Msg("error finding service order with id")
		return nil, err
	}

	services, err := u.addServiceToAdditionalRepair(ctx, adr.Services)
	if err != nil {
		return nil, err
	}

	partsSupplies, err := u.addPartsSupplyToAdditionalRepair(ctx, adr.PartsSupplies)
	if err != nil {
		return nil, err
	}

	newAdditionalRepair := &entities.AdditionalRepair{
		ServiceOrderID: adr.ServiceOrderID,
		Description:    adr.Description,
		Status:         valueobject.StatusARAberta,
		Services:       services,
		PartsSupplies:  partsSupplies,
	}

	if err = u.partsSupplyRepo.AuthorizeReserve(ctx, partsSupplies); err != nil {
		logger.Error().Err(err).Msg("Error authorizing parts supply reserve")
		return nil, err
	}

	created, err := u.repo.CreateAdditionalRepair(ctx, newAdditionalRepair)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating additional repair")
		return nil, err
	}

	if err = u.partsSupplyRepo.Reserve(ctx, created.PartsSupplies); err != nil {
		logger.Error().Err(err).Msg("Error reserving parts supply")
		if err := u.Rollback(ctx, created, false); err != nil {
			logger.Error().Err(err).Msg("Error rolling back additional repair")
		}
		return nil, err
	}

	estimate, err := u.billRepo.CreateEstimate(ctx, nil, created)
	if err != nil {
		logger.Error().Err(err).Any("service_order_id", created.ServiceOrderID).Msg("error creating estimate with id")
		if err := u.Rollback(ctx, created, false); err != nil {
			logger.Error().Err(err).Msg("Error rolling back additional repair")
		}
		return nil, err
	}
	if estimate == nil || estimate.ID == "" {
		logger.Error().Msg("estimate not found")
		if err := u.Rollback(ctx, created, false); err != nil {
			logger.Error().Err(err).Msg("Error rolling back additional repair")
		}
		return nil, errors.New("invalid estimate, failed to create estimate")
	}

	created.Estimate = estimate
	created.Status = valueobject.StatusARAguardandoAprovacao

	_, err = u.repo.UpdateAdditionalRepair(ctx, created)
	if err != nil {
		logger.Error().Err(err).Msg("Error updating additional repair")
		if err := u.Rollback(ctx, created, false); err != nil {
			logger.Error().Err(err).Msg("Error rolling back additional repair")
		}
		return nil, err
	}

	return created, nil
}

func (u *AdditionalRepairUseCase) CancelAdditionalRepair(ctx context.Context, additionalRepairId string) (*entities.AdditionalRepair, error) {
	logger := logs.LoggerWithContext(ctx)

	additionalRepair, err := u.repo.GetByID(ctx, additionalRepairId)
	if err != nil {
		logger.Error().Err(err).Msg("Error getting additional repair")
		return nil, err
	}

	if additionalRepair.Status.IsCancelada() {
		logger.Error().Msg("additional repair already canceled")
		return nil, errors.New("additional repair already canceled")
	}

	additionalRepair.Status = valueobject.StatusARCancelada

	_, err = u.repo.UpdateAdditionalRepair(ctx, additionalRepair)
	if err != nil {
		logger.Error().Err(err).Msg("Error updating additional repair")
		return nil, err
	}

	return additionalRepair, nil
}

func (u *AdditionalRepairUseCase) Rollback(ctx context.Context, additionalRepair *entities.AdditionalRepair, hasReservedPartsSupply bool) error {
	logger := logs.LoggerWithContext(ctx)
	txn := newrelic.FromContext(ctx)
	if txn != nil {
		seg := txn.StartSegment("AdditionalRepairUseCase.Rollback")
		defer seg.End()
	}

	if additionalRepair == nil {
		logger.Error().Msg("additional repair not found")
		return ErrAdditionalRepairNotFound
	}

	if additionalRepair.Status.IsCancelada() || additionalRepair.Status.IsRejeitada() {
		logger.Error().Msg("additional repair already canceled or rejected")
		return ErrAdditionalRepairNotFound
	}

	if additionalRepair.Status.IsAberta() || additionalRepair.Status.IsAguardandoAprovacao() || additionalRepair.Status.IsAprovada() {
		if hasReservedPartsSupply {
			err := u.partsSupplyRepo.Release(ctx, additionalRepair.PartsSupplies)
			if err != nil {
				logger.Error().Err(err).Msg("Error releasing parts supply")
				return err
			}
		}

		if additionalRepair.Estimate != nil {
			estimate, err := u.billRepo.CancelEstimate(ctx, nil, additionalRepair)
			if err != nil {
				logger.Error().Err(err).Msg("Error deleting estimate")
				return err
			}
			additionalRepair.Estimate = estimate
			_, err = u.repo.UpdateAdditionalRepair(ctx, additionalRepair)
			if err != nil {
				logger.Error().Err(err).Msg("Error updating additional repair")
				return err
			}
		}
	}

	_, err := u.CancelAdditionalRepair(ctx, additionalRepair.ID)
	if err != nil {
		logger.Error().Err(err).Msg("Error canceling additional repair")
		return err
	}

	return nil
}

func (u *AdditionalRepairUseCase) GetAdditionalRepair(ctx context.Context, additionalRepairId string) (*entities.AdditionalRepair, error) {
	logger := logs.LoggerWithContext(ctx)

	additionalRepair, err := u.repo.GetByID(ctx, additionalRepairId)
	if err != nil {
		logger.Error().Err(err).Any("additional_repair_id", additionalRepairId).Msg("error finding additional repair with id")
		return nil, err
	}
	if additionalRepair.ID == "" {
		logger.Error().Msg("additional repair not found")
		return nil, ErrAdditionalRepairNotFound
	}
	return additionalRepair, nil
}

func (u *AdditionalRepairUseCase) CustomerApprovalStatus(ctx context.Context, additionalRepairId string, flow string) (*entities.AdditionalRepair, error) {
	logger := logs.LoggerWithContext(ctx)

	additionalRepair, err := u.repo.GetByID(ctx, additionalRepairId)
	if err != nil {
		logger.Error().Err(err).Any("additional_repair_id", additionalRepairId).Msg("error finding additional repair with id")
		return nil, err
	}
	if additionalRepair.ID == "" {
		return nil, ErrAdditionalRepairNotFound
	}

	if !additionalRepair.Status.IsAguardandoAprovacao() {
		logs.Logger().Error().Msg("invalid additional repair status")
		return nil, ErrStatusNotPermitted
	}

	if flow == REJECTED_FLOW {
		if err := u.partsSupplyRepo.Release(ctx, additionalRepair.PartsSupplies); err != nil {
			logger.Error().Err(err).Msg("Error releasing parts supply")
			return nil, err
		}

		estimate, err := u.billRepo.RejectEstimate(ctx, nil, additionalRepair)
		if err != nil {
			logger.Error().Err(err).Msg("Error rejecting estimate")
			return nil, err
		}
		additionalRepair.Estimate = estimate
		additionalRepair.Status = valueobject.StatusARRejeitada
	} else if flow == APPROVED_FLOW {
		if err := u.partsSupplyRepo.WriteOff(ctx, additionalRepair.PartsSupplies); err != nil {
			logger.Error().Err(err).Msg("Error writing off parts supply")
			return nil, err
		}
		estimate, err := u.billRepo.ApproveEstimate(ctx, nil, additionalRepair)
		if err != nil {
			logger.Error().Err(err).Msg("Error approving estimate")
			return nil, err
		}
		additionalRepair.Estimate = estimate
		additionalRepair.Status = valueobject.StatusAAprovada
	}

	_, err = u.repo.UpdateAdditionalRepair(ctx, additionalRepair)
	if err != nil {
		logger.Error().Err(err).Msg("Error updating additional repair")
		return nil, err
	}
	return additionalRepair, nil
}

func (u *AdditionalRepairUseCase) addPartsSupplyToAdditionalRepair(ctx context.Context, partsSupplyId []entities.PartsSupply) ([]entities.PartsSupply, error) {
	var listPartsSupply []entities.PartsSupply
	for _, ps := range partsSupplyId {
		psDto, err := u.partsSupplyRepo.GetByID(ctx, ps.ID)
		if err != nil {
			logs.Logger().Error().Err(err).Msg("error finding parts supply with id")
		}
		if psDto.ID == "" {
			logs.Logger().Error().Msg("parts supply with id not found")
			return listPartsSupply, ErrServiceNotFound
		}
		listPartsSupply = append(listPartsSupply, *psDto)
	}
	return listPartsSupply, nil
}

func (u *AdditionalRepairUseCase) addServiceToAdditionalRepair(ctx context.Context, services []entities.Service) ([]entities.Service, error) {
	var listServices []entities.Service
	for _, s := range services {
		serviceDto, err := u.serviceRepo.GetByID(ctx, s.ID)
		if err != nil {
			logs.Logger().Error().Err(err).Msg("error finding service with id")
		}
		if serviceDto.ID == "" {
			logs.Logger().Error().Msg("service with id not found")
			return listServices, ErrServiceNotFound
		}
		listServices = append(listServices, *serviceDto)
	}
	return listServices, nil
}