package usecase

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/logs"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/observability"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/adapter/operations"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/interfaces"

	"github.com/newrelic/go-agent/v3/newrelic"
)

// operation flow
const (
	DIAGNOSIS        = "diagnosis"
	ESTIMATE         = "estimate"
	EXECUTION        = "execution"
	DELIVERY         = "delivery"
	ESTIMATE_APPROVE = "estimate_approve"
	ESTIMATE_REJECT  = "estimate_reject"
	ESTIMATE_CANCEL  = "estimate_cancel"
)

var (
	ErrServiceOrderNotFound               = errors.New("service order not found")
	ErrInvalidTransitionStatusToDiagnosis = errors.New("invalid transition status to diagnosis")
	ErrInvalidTransitionStatusToExecution = errors.New("invalid transition status to execution")
	ErrInvalidTransitionStatusToDelivery  = errors.New("invalid transition status to delivery")
	ErrInvalidTransitionStatusToEstimate  = errors.New("invalid transition status to estimate")
	ErrInvalidStatus                      = errors.New("invalid service order status")
	ErrInsufficientPartsSupply            = errors.New("insufficient parts supply available")
	ErrInvalidFlow                        = errors.New("invalid flow")
	ErrInvalidID                          = errors.New("invalid ID provided")
	ErrCustomerNotFound                   = errors.New("customer not found")
	ErrVehicleNotFound                    = errors.New("vehicle not found")
	ErrServiceNotFound                    = errors.New("service not found")
	ErrPartsSupplyNotFound                = errors.New("parts supply not found")
	ErrInvalidEstimateValue               = errors.New("invalid estimate value")
)

type IServiceOrderUseCase interface {
	CreateServiceOrder(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error)
	DiagnosisServiceOrder(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error)
	EstimateServiceOrder(ctx context.Context, serviceOrderID uint, operation string) (*entities.ServiceOrder, error)
	// ExecutionServiceOrder(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error)
	// FinishServiceOrderExecution(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error)
	// PaymentServiceOrder(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error)
	// DeliveryServiceOrder(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error)
	UpdateServiceOrder(ctx context.Context, serviceOrder *entities.ServiceOrder, flow string) (*entities.ServiceOrder, error)
	GetServiceOrder(ctx context.Context, serviceOrder entities.ServiceOrder, isFullData bool) (*entities.ServiceOrder, error)
	ListServiceOrders(ctx context.Context) ([]*entities.ServiceOrder, error)
}

type ServiceOrderUseCase struct {
	repo               interfaces.IServiceOrderGateway
	vehicleRepo        interfaces.IVehicleGateway
	customerRepo       interfaces.ICustomerGateway
	serviceRepo        interfaces.IServiceGateway
	partsSupplyRepo    interfaces.IPartsSupplyGateway
	billingServiceRepo interfaces.IBillingServiceGateway
}

var _ IServiceOrderUseCase = (*ServiceOrderUseCase)(nil)

func NewServiceOrderUseCase(
	repo interfaces.IServiceOrderGateway,
	vehicleRepo interfaces.IVehicleGateway,
	customerRepo interfaces.ICustomerGateway,
	serviceRepo interfaces.IServiceGateway,
	partsSupplyRepo interfaces.IPartsSupplyGateway,
	billingServiceRepo interfaces.IBillingServiceGateway,
) *ServiceOrderUseCase {
	return &ServiceOrderUseCase{
		repo:               repo,
		vehicleRepo:        vehicleRepo,
		customerRepo:       customerRepo,
		serviceRepo:        serviceRepo,
		partsSupplyRepo:    partsSupplyRepo,
		billingServiceRepo: billingServiceRepo,
	}
}

// CreateServiceOrder creates a new service order after validating the vehicle and customer.
// It sets the initial status of the service order to "Recebida".
// If the vehicle or customer validation fails, it logs the error and returns it.
func (u *ServiceOrderUseCase) CreateServiceOrder(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error) {
	logger := logs.Logger()
	txn := newrelic.FromContext(ctx)
	if txn != nil {
		logger = logs.LoggerWithContext(ctx)
	}

	// Validation segment - track validation operations separately
	var seg *newrelic.Segment
	if txn != nil {
		seg = txn.StartSegment("ServiceOrderUseCase.CreateServiceOrder.Validation")
		defer seg.End()
	}

	if serviceOrder == nil {
		return nil, ErrInvalidID
	}

	err := u.validateVehicle(ctx, serviceOrder)
	if err != nil {
		logger.Error().Err(err).Uint("vehicle_id", serviceOrder.VehicleID).Msg("Error validating vehicle")
		return nil, err
	}

	err = u.validateCustomer(ctx, serviceOrder)
	if err != nil {
		logger.Error().Err(err).Uint("customer_id", serviceOrder.CustomerID).Msg("Error validating customer")
		return nil, err
	}

	newServiceOrder := &entities.ServiceOrder{
		CustomerID: serviceOrder.CustomerID,
		VehicleID:  serviceOrder.VehicleID,
		Status:     valueobject.StatusRecebida,
	}

	// Repository operation segment - track database operations separately
	if txn != nil {
		repoSeg := txn.StartSegment("ServiceOrderUseCase.CreateServiceOrder.Repository")
		defer repoSeg.End()
	}

	register, err := u.repo.Create(ctx, newServiceOrder)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating service order")
		return nil, err
	}

	return register, nil
}

// UpdateServiceOrder updates an existing service order.
// TODO: Deprecate this method
func (u *ServiceOrderUseCase) UpdateServiceOrder(ctx context.Context, request *entities.ServiceOrder, flow string) (*entities.ServiceOrder, error) {
	logger := logs.Logger()
	if txn := newrelic.FromContext(ctx); txn != nil {
		logger = logs.LoggerWithContext(ctx)
	}

	var update = &entities.ServiceOrder{}

	if request.ID == 0 {
		return nil, ErrInvalidID
	}

	update.ID = request.ID

	serviceOrderRecord, err := u.repo.GetByID(ctx, request.ID, false)
	if err != nil {
		logger.Error().Err(err).Any("OS_ID", update.ID).Msg("Error finding service order with id")
		return nil, err
	}
	if serviceOrderRecord == nil {
		logger.Error().Any("OS_ID", update.ID).Msg("Service order with id not found")
		return nil, ErrServiceOrderNotFound
	}

	switch flow {
	case EXECUTION:
		update, err = u.validateExecution(ctx, request, serviceOrderRecord, update)
		if err != nil {
			logger.Error().Err(err).Msg("Error validating execution")
			return nil, err
		}
	case DELIVERY:
		update, err = u.validateDelivery(ctx, request, serviceOrderRecord, update)
		if err != nil {
			logger.Error().Err(err).Msg("Error validating delivery")
			return nil, err
		}
	default:
		logger.Error().Str("flow", flow).Msg("Invalid flow")
		return nil, ErrInvalidFlow
	}

	err = u.repo.Update(ctx, update)
	if err != nil {
		logger.Error().Err(err).Msg("Error updating service order")
		return nil, err
	}

	updatedSO, err := u.repo.GetByID(ctx, request.ID, false)
	if err != nil || updatedSO == nil {
		return nil, err
	}

	if updatedSO.PaymentID != nil {
		updatedSO.PaymentID = nil
	}

	return updatedSO, nil
}

func (u *ServiceOrderUseCase) DiagnosisServiceOrder(ctx context.Context, request *entities.ServiceOrder) (*entities.ServiceOrder, error) {
	logger := logs.Logger()

	if txn := newrelic.FromContext(ctx); txn != nil {
		logger = logs.LoggerWithContext(ctx)
		startSegment := txn.StartSegment("ServiceOrderUseCase.DiagnosisServiceOrder")
		defer startSegment.End()
	}

	result, err := u.checkIfServiceOrderExists(ctx, request.ID)
	if err != nil {
		logger.Error().Err(err).Msg("Error finding service order with id")
		return nil, err
	}

	actualStatus := result.Status
	request.Status = actualStatus

	serviceOrder, err := u.validateDiagnosis(ctx, request)
	if err != nil {
		logger.Error().Err(err).Msg("Error validating diagnosis")
		return nil, err
	}

	if serviceOrder.Status.IsEmDiagnostico() {
		return serviceOrder, nil
	}

	var estimate *entities.Estimate
	estimate, err = u.billingServiceRepo.CreateEstimate(ctx, serviceOrder)
	if err != nil {
		logger.Error().Err(err).Msg("Error calculating estimate")
		return nil, err
	}

	if estimate == nil {
		logger.Error().Msg("Error calculating estimate: the value from estimate creation is nil")
		return nil, ErrInvalidEstimateValue
	}

	serviceOrder.Status = valueobject.StatusAguardandoAprovacao

	err = u.repo.Update(ctx, serviceOrder)
	if err != nil {
		logger.Error().Err(err).Msg("Error updating service order")
		return nil, err
	}

	return serviceOrder, nil
}

// ValidateDiagnosis checks if the service order status is valid for diagnosis.
// If the status is "Recebida" or "Em Diagnostico" and the request status is "Em Diagnostico",
// it updates the service order status to "Em Diagnostico".
func (u *ServiceOrderUseCase) validateDiagnosis(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error) {
	logger := logs.LoggerWithContext(ctx)
	if txn := newrelic.FromContext(ctx); txn != nil {
		logger = logs.LoggerWithContext(ctx)
		startSegment := txn.StartSegment("ServiceOrderUseCase.validateDiagnosis")
		defer startSegment.End()
	}
	if !serviceOrder.Status.IsValid() {
		return nil, ErrInvalidStatus
	}

	if !(serviceOrder.Status.IsRecebida() || serviceOrder.Status.IsEmDiagnostico()) {
		return nil, ErrInvalidTransitionStatusToDiagnosis
	}

	serviceOrder.Status = valueobject.StatusEmDiagnostico

	if len(serviceOrder.Services) == 0 && len(serviceOrder.PartsSupplies) == 0 {
		observability.IncrementCounter(ctx, "service_order_diagnosis_pending_services_parts_supplies", nil)
		return serviceOrder, nil
	}

	if len(serviceOrder.Services) == 0 || len(serviceOrder.PartsSupplies) == 0 {
		observability.IncrementCounter(ctx,
			"service_order_diagnosis_no_services_or_no_parts_supplies",
			map[string]string{
				"services":       fmt.Sprint(len(serviceOrder.Services)),
				"parts_supplies": fmt.Sprint(len(serviceOrder.PartsSupplies)),
			},
		)
	}

	err := u.validateIfServicesExists(ctx, serviceOrder.Services)
	if err != nil {
		logs.Logger().Error().Err(err).Any("services", serviceOrder.Services).Msg("Error validating services")
		return nil, err
	}

	err = u.partsSupplyRepo.AuthorizeReserve(ctx, serviceOrder.PartsSupplies)
	if err != nil {
		logs.Logger().Error().Err(err).Msg("Error authorizing reserve parts supply")
		return nil, err
	}

	err = u.partsSupplyRepo.Reserve(ctx, serviceOrder.PartsSupplies)
	if err != nil {
		logger.Error().Err(err).Any("parts_supply", serviceOrder.PartsSupplies).Msg("Error reserving parts supply")
		return nil, err
	}
	return serviceOrder, nil
}

func (u *ServiceOrderUseCase) EstimateServiceOrder(ctx context.Context, serviceOrderID uint, operation string) (*entities.ServiceOrder, error) {
	logger := logs.LoggerWithContext(ctx)

	serviceOrder, err := u.checkIfServiceOrderExists(ctx, serviceOrderID)
	if err != nil {
		logger.Error().Err(err).Msg("Error finding service order with id")
		return nil, err
	}

	if !serviceOrder.Status.IsValid() {
		return nil, ErrInvalidStatus
	}

	if serviceOrder.Status.IsAguardandoAprovacao() {
		return nil, ErrInvalidTransitionStatusToEstimate
	}

	factory := operations.NewEstimateStrategyFactory()
	strategy, err := factory.GetStrategy(operation)
	if err != nil {
		return nil, err
	}

	return strategy.Execute(ctx, serviceOrder)
}

func (u *ServiceOrderUseCase) validateExecution(ctx context.Context, request *entities.ServiceOrder, current *entities.ServiceOrder, update *entities.ServiceOrder) (*entities.ServiceOrder, error) {
	oldStatus := current.Status
	currentTime := time.Now()

	if !request.Status.IsValid() {
		return nil, ErrInvalidStatus
	}
	if oldStatus.IsAprovada() && request.Status.IsEmExecucao() {
		update.Status = valueobject.StatusEmExecucao
		update.StartedExecutionDate = &currentTime
		return update, nil
	}
	if oldStatus.IsEmExecucao() && request.Status.IsFinalizada() {
		executionDuration, err := CalculateExecutionDurationInHours(current.StartedExecutionDate, &currentTime)
		if err != nil {
			return nil, err
		}
		update.Status = valueobject.StatusFinalizada
		update.FinalExecutionDate = &currentTime
		update.ExecutionDurationInHours = executionDuration
		return update, nil
	}
	return nil, ErrInvalidTransitionStatusToExecution
}

func (u *ServiceOrderUseCase) validateDelivery(ctx context.Context, request *entities.ServiceOrder, current *entities.ServiceOrder, update *entities.ServiceOrder) (*entities.ServiceOrder, error) {
	oldStatus := current.Status

	if !request.Status.IsValid() {
		return nil, ErrInvalidStatus
	}

	if oldStatus.IsFinalizada() && request.Status.IsEntregue() {
		if current.PaymentID == nil {
			return nil, errors.New("payment information is required for delivery")
		}
		update.Status = valueobject.StatusEntregue
		return update, nil
	}
	return nil, ErrInvalidTransitionStatusToDelivery
}

func (u *ServiceOrderUseCase) validateVehicle(ctx context.Context, serviceOrder *entities.ServiceOrder) error {
	logger := logs.LoggerWithContext(ctx)

	if serviceOrder.VehicleID != 0 {
		vehicle, err := u.vehicleRepo.FindByID(serviceOrder.VehicleID)
		if err != nil {
			logger.Error().Err(err).Any("vehicle_id", serviceOrder.VehicleID).Msg("error finding vehicle with id")
			return err
		}
		if vehicle == nil || vehicle.ID == 0 {
			return ErrVehicleNotFound
		}
		return nil
	}
	return ErrInvalidID
}

func (u *ServiceOrderUseCase) validateCustomer(ctx context.Context, serviceOrder *entities.ServiceOrder) error {
	logger := logs.LoggerWithContext(ctx)

	if serviceOrder.CustomerID != 0 {
		customer, err := u.customerRepo.GetByID(serviceOrder.CustomerID)
		if err != nil {
			logger.Error().Err(err).Any("customer_id", serviceOrder.CustomerID).Msg("error finding customer with id")
			return err
		}
		if customer == nil {
			return ErrCustomerNotFound
		}
		return nil
	}
	return ErrInvalidID
}

func (u *ServiceOrderUseCase) getSeviceById(ctx context.Context, s entities.Service) (*entities.Service, error) {
	logger := logs.LoggerWithContext(ctx)

	if s.ID == 0 {
		return nil, ErrInvalidID
	}
	result, err := u.serviceRepo.GetByID(ctx, s.ID)
	if err != nil {
		logger.Error().Err(err).Any("service_id", s.ID).Msg("error finding service with id")
		return nil, err
	}
	if result == nil {
		logger.Error().Any("service_id", s.ID).Msg("service with id not found")
		return nil, ErrServiceNotFound
	}
	return result, nil
}

func (u *ServiceOrderUseCase) getPartsSupplyByID(ctx context.Context, id uint) (*entities.PartsSupply, error) {
	result, err := u.partsSupplyRepo.GetByID(ctx, id)
	if err != nil {
		logs.Logger().Error().Err(err).Any("parts_supply_id", id).Msg("error finding parts supply with id")
		return nil, err
	}
	if result == nil {
		logs.Logger().Error().Any("parts_supply_id", id).Msg("parts supply with id not found")
		return nil, ErrPartsSupplyNotFound
	}
	return result, nil
}

func (u *ServiceOrderUseCase) validateIfServicesExists(ctx context.Context, services []entities.Service) error {
	for _, s := range services {
		_, err := u.getSeviceById(ctx, s)
		if err != nil {
			return err
		}
	}
	return nil
}

// releaseReservedPartsSupply is when a service order is approved - Baixa de estoque
func (u *ServiceOrderUseCase) releaseReservedPartsSupply(ctx context.Context, partsSupply []entities.PartsSupply) error {
	err := u.partsSupplyRepo.Release(ctx, partsSupply)
	if err != nil {
		logs.Logger().Error().Err(err).Msg("error releasing reserved parts supply")
		return err
	}
	return nil
}

func (u *ServiceOrderUseCase) GetServiceOrder(ctx context.Context, serviceOrder entities.ServiceOrder, isFullData bool) (*entities.ServiceOrder, error) {
	logger := logs.LoggerWithContext(ctx)

	serviceOrderRecord, err := u.repo.GetByID(ctx, serviceOrder.ID, isFullData)
	if err != nil {
		logger.Error().Err(err).Any("service_order_id", serviceOrder.ID).Msg("error finding service order")
		return nil, err
	}
	if serviceOrderRecord == nil {
		logger.Error().Any("service_order_id", serviceOrder.ID).Msg("service order not found")
		return nil, ErrServiceOrderNotFound
	}
	return serviceOrderRecord, nil
}

func (u *ServiceOrderUseCase) ListServiceOrders(ctx context.Context) ([]*entities.ServiceOrder, error) {
	logger := logs.LoggerWithContext(ctx)

	serviceOrders, err := u.repo.List(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("error listing service orders")
		return nil, err
	}

	filtered := make([]*entities.ServiceOrder, 0, len(serviceOrders))
	for _, so := range serviceOrders {
		if so == nil {
			continue
		}

		status := so.Status
		if status.IsFinalizada() || status.IsEntregue() {
			// Skip finalized or delivered orders from the listing.
			continue
		}

		filtered = append(filtered, so)
	}

	sort.SliceStable(filtered, func(i, j int) bool {
		si := filtered[i]
		sj := filtered[j]

		pi := u.serviceOrderPriority(si.Status)
		pj := u.serviceOrderPriority(sj.Status)
		if pi != pj {
			return pi < pj
		}

		ti := si.CreatedAt
		tj := sj.CreatedAt

		if ti == nil && tj == nil {
			return si.ID < sj.ID
		}
		if ti == nil {
			return false
		}
		if tj == nil {
			return true
		}
		if !ti.Equal(*tj) {
			return ti.Before(*tj)
		}
		return si.ID < sj.ID
	})

	return filtered, nil
}

func (u *ServiceOrderUseCase) getServicesByIDs(ctx context.Context, services []entities.Service) ([]entities.Service, error) {
	if len(services) == 0 {
		return nil, errors.New("no services provided")
	}

	var serviceDb []entities.Service
	for _, s := range services {
		item, err := u.serviceRepo.GetByID(ctx, s.ID)
		if err != nil {
			logs.Logger().Error().Err(err).Any("service_id", s.ID).Msg("error getting service by ID")
			return nil, err
		}
		serviceDb = append(serviceDb, *item)
	}

	// Assuming we have a service repository to get the services by IDs
	return serviceDb, nil
}

func (u *ServiceOrderUseCase) getPartsSupplyByIDs(ctx context.Context, partsSupplies []entities.PartsSupply) ([]entities.PartsSupply, error) {
	if len(partsSupplies) == 0 {
		return nil, errors.New("no services provided")
	}

	var psDb []entities.PartsSupply
	for _, ps := range partsSupplies {
		item, err := u.partsSupplyRepo.GetByID(ctx, ps.ID)
		if err != nil {
			logs.Logger().Error().Err(err).Any("parts_supply_id", ps.ID).Msg("error getting Parts Supplies by ID")
			return nil, err
		}
		psDb = append(psDb, *item)
	}

	return psDb, nil
}

func (u *ServiceOrderUseCase) getPartsSuppliesByServiceOrderID(ctx context.Context, serviceOrderID uint) ([]entities.PartsSupply, error) {
	if serviceOrderID == 0 {
		return nil, ErrInvalidID
	}

	partsSupplies, err := u.partsSupplyRepo.GetByServiceOrderID(ctx, serviceOrderID)
	if err != nil {
		logs.Logger().Error().Err(err).Any("service_order_id", serviceOrderID).Msg("error getting Parts Supplies by Service Order ID")
		return nil, err
	}

	return partsSupplies, nil
}

func (u *ServiceOrderUseCase) checkIfServiceOrderExists(ctx context.Context, id uint) (*entities.ServiceOrder, error) {
	logger := logs.LoggerWithContext(ctx)
	serviceOrder, err := u.repo.GetByID(ctx, id, false)
	if err != nil {
		logger.Error().Err(err).Any("OS_ID", id).Msg("Error finding service order with id")
		return nil, err
	}

	if serviceOrder == nil {
		logger.Error().Any("OS_ID", id).Msg("Service order with id not found")
		return nil, ErrServiceOrderNotFound
	}
	return serviceOrder, nil
}

func (u *ServiceOrderUseCase) serviceOrderPriority(status valueobject.ServiceOrderStatus) int {
	switch {
	case status.IsEmExecucao():
		return 0
	case status.IsAguardandoAprovacao():
		return 1
	case status.IsEmDiagnostico():
		return 2
	case status.IsRecebida():
		return 3
	default:
		return 4
	}
}
