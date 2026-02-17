package usecase

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/logs"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/observability"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/adapter/operations"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/constants"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/interfaces"

	"github.com/newrelic/go-agent/v3/newrelic"
)

var (
	ErrServiceOrderNotFound               = errors.New("service order not found")
	ErrInvalidTransitionStatusToDiagnosis = errors.New("invalid transition status to diagnosis")
	ErrInvalidTransitionStatusToExecution = errors.New("invalid transition status to execution")
	ErrInvalidTransitionStatusToDelivery  = errors.New("invalid transition status to delivery")
	ErrInvalidTransitionStatusToEstimate  = errors.New("invalid transition status to estimate")
	ErrInvalidTransitionStatusToCancel    = errors.New("invalid transition status to cancel")
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
	EstimateServiceOrder(ctx context.Context, serviceOrderID string, operation string) (*entities.ServiceOrder, error)
	ExecutionServiceOrder(ctx context.Context, serviceOrderID string, operation string) (*entities.ServiceOrder, error)
	PaymentServiceOrder(ctx context.Context, serviceOrderID string) (*entities.ServiceOrder, error)
	DeliveryServiceOrder(ctx context.Context, serviceOrderID string) (*entities.ServiceOrder, error)
	GetServiceOrder(ctx context.Context, serviceOrder entities.ServiceOrder, isFullData bool) (*entities.ServiceOrder, error)
	ListServiceOrders(ctx context.Context) ([]*entities.ServiceOrder, error)
	CancelServiceOrder(ctx context.Context, serviceOrderID string) (*entities.ServiceOrder, error)
}

type ServiceOrderUseCase struct {
	repo               interfaces.IServiceOrderGateway
	vehicleRepo        interfaces.IVehicleGateway
	customerRepo       interfaces.ICustomerGateway
	serviceRepo        interfaces.IServiceGateway
	partsSupplyRepo    interfaces.IPartsSupplyGateway
	billingServiceRepo interfaces.IBillingServiceGateway
	executionRepo      interfaces.IExecutionGateway
}

var _ IServiceOrderUseCase = (*ServiceOrderUseCase)(nil)

func NewServiceOrderUseCase(
	repo interfaces.IServiceOrderGateway,
	vehicleRepo interfaces.IVehicleGateway,
	customerRepo interfaces.ICustomerGateway,
	serviceRepo interfaces.IServiceGateway,
	partsSupplyRepo interfaces.IPartsSupplyGateway,
	billingServiceRepo interfaces.IBillingServiceGateway,
	executionRepo interfaces.IExecutionGateway,
) *ServiceOrderUseCase {
	return &ServiceOrderUseCase{
		repo:               repo,
		vehicleRepo:        vehicleRepo,
		customerRepo:       customerRepo,
		serviceRepo:        serviceRepo,
		partsSupplyRepo:    partsSupplyRepo,
		billingServiceRepo: billingServiceRepo,
		executionRepo:      executionRepo,
	}
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

// CreateServiceOrder creates a new service order after validating the vehicle and customer.
// It sets the initial status of the service order to "Recebida".
// If the vehicle or customer validation fails, it logs the error and returns it.
func (u *ServiceOrderUseCase) CreateServiceOrder(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error) {
	logger := logs.Logger()
	txn := newrelic.FromContext(ctx)
	if txn != nil {
		seg := txn.StartSegment("ServiceOrderUseCase.CreateServiceOrder.Validation")
		defer seg.End()
		logger = logs.LoggerWithContext(ctx)
	}

	if serviceOrder == nil {
		return nil, ErrInvalidID
	}

	err := u.validateVehicle(ctx, serviceOrder)
	if err != nil {
		logger.Error().Err(err).Str("vehicle_id", serviceOrder.VehicleID).Msg("Error validating vehicle")
		return nil, err
	}

	err = u.validateCustomer(ctx, serviceOrder)
	if err != nil {
		logger.Error().Err(err).Str("customer_id", serviceOrder.CustomerID).Msg("Error validating customer")
		return nil, err
	}

	newServiceOrder := &entities.ServiceOrder{
		CustomerID: serviceOrder.CustomerID,
		VehicleID:  serviceOrder.VehicleID,
		Status:     valueobject.StatusRecebida,
	}

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

func (u *ServiceOrderUseCase) DiagnosisServiceOrder(ctx context.Context, request *entities.ServiceOrder) (*entities.ServiceOrder, error) {
	logger := logs.Logger()

	if txn := newrelic.FromContext(ctx); txn != nil {
		logger = logs.LoggerWithContext(ctx)
		startSegment := txn.StartSegment("ServiceOrderUseCase.DiagnosisServiceOrder")
		defer startSegment.End()
	}

	result, err := u.validateServiceOrderExists(ctx, request.ID)
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

	err = u.partsSupplyRepo.Reserve(ctx, serviceOrder.PartsSupplies)
	if err != nil {
		logger.Error().Err(err).Any("parts_supply", serviceOrder.PartsSupplies).Msg("Error reserving parts supply")
		return nil, err
	}

	if !serviceOrder.IsDiagnosisPending() {
		var estimate *entities.Estimate
		estimate, err = u.billingServiceRepo.CreateEstimate(ctx, serviceOrder, nil)
		if err != nil {
			logger.Error().Err(err).Msg("Error calculating estimate")
			return nil, err
		}

		if estimate == nil {
			logger.Error().Msg("Error calculating estimate: the value from estimate creation is nil")
			return nil, ErrInvalidEstimateValue
		}

		serviceOrder.Status = valueobject.StatusAguardandoAprovacao
		serviceOrder.Estimate = estimate
	} else {
		serviceOrder.Status = valueobject.StatusEmDiagnostico
	}

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
	logger := logs.Logger()
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

	if serviceOrder.IsDiagnosisPending() {
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

	err := u.validateServicesExists(ctx, serviceOrder.Services)
	if err != nil {
		logger.Error().Err(err).Any("services", serviceOrder.Services).Msg("Error validating services")
		return nil, err
	}

	err = u.partsSupplyRepo.AuthorizeReserve(ctx, serviceOrder.PartsSupplies)
	if err != nil {
		logger.Error().Err(err).Msg("Error authorizing reserve parts supply")
		return nil, err
	}
	return serviceOrder, nil
}

func (u *ServiceOrderUseCase) EstimateServiceOrder(ctx context.Context, serviceOrderID string, operation string) (*entities.ServiceOrder, error) {
	logger := logs.LoggerWithContext(ctx)

	serviceOrder, err := u.validateServiceOrderExists(ctx, serviceOrderID)
	if err != nil {
		logger.Error().Err(err).Msg("Error finding service order with id")
		return nil, err
	}

	if !serviceOrder.Status.IsValid() {
		return nil, ErrInvalidStatus
	}

	if !serviceOrder.Status.IsAguardandoAprovacao() {
		return nil, ErrInvalidTransitionStatusToEstimate
	}

	factory := operations.NewEstimateStrategyFactory(u.partsSupplyRepo, u.billingServiceRepo)
	strategy, err := factory.GetStrategy(operation)
	if err != nil {
		return nil, err
	}

	result, err := strategy.Execute(ctx, serviceOrder)
	if err != nil {
		logger.Error().Err(err).Msg("Error approving estimate")
		return nil, err
	}

	err = u.repo.Update(ctx, result)
	if err != nil {
		logger.Error().Err(err).Msg("Error updating service order")
		return nil, err
	}

	return result, nil
}

func (u *ServiceOrderUseCase) ExecutionServiceOrder(ctx context.Context, serviceOrderID string, operation string) (*entities.ServiceOrder, error) {
	logger := logs.LoggerWithContext(ctx)
	if txn := newrelic.FromContext(ctx); txn != nil {
		startSegment := txn.StartSegment("ServiceOrderUseCase.validateExecution")
		defer startSegment.End()
	}

	serviceOrder, err := u.validateServiceOrderExists(ctx, serviceOrderID)
	if err != nil {
		logger.Error().Err(err).Msg("Error finding service order with id")
		return nil, err
	}

	if !serviceOrder.Status.IsValid() {
		return nil, ErrInvalidStatus
	}

	if !(serviceOrder.Status.IsAprovada() || serviceOrder.Status.IsEmExecucao()) {
		return nil, ErrInvalidTransitionStatusToExecution
	}

	if operation == constants.EXECUTION_START {
		execution, err := u.executionRepo.CreateExecution(ctx, serviceOrder)
		if err != nil {
			return nil, err
		}
		serviceOrder.Execution = execution
		serviceOrder.Status = valueobject.StatusEmExecucao

	} else if operation == constants.EXECUTION_FINISH {
		execution, err := u.executionRepo.FinishExecution(ctx, serviceOrder)
		if err != nil {
			return nil, err
		}
		serviceOrder.Execution = execution
		serviceOrder.Status = valueobject.StatusFinalizada
	}

	err = u.repo.Update(ctx, serviceOrder)
	if err != nil {
		logger.Error().Err(err).Msg("Error updating service order")
		return nil, err
	}

	return serviceOrder, nil
}

func (u *ServiceOrderUseCase) DeliveryServiceOrder(ctx context.Context, serviceOrderID string) (*entities.ServiceOrder, error) {
	logger := logs.LoggerWithContext(ctx)
	if txn := newrelic.FromContext(ctx); txn != nil {
		startSegment := txn.StartSegment("ServiceOrderUseCase.validateExecution")
		defer startSegment.End()
	}

	serviceOrder, err := u.validateServiceOrderExists(ctx, serviceOrderID)
	if err != nil {
		logger.Error().Err(err).Msg("Error finding service order with id")
		return nil, err
	}

	if !serviceOrder.Status.IsFinalizada() {
		return nil, ErrInvalidTransitionStatusToDelivery
	}

	if serviceOrder.Estimate == nil || serviceOrder.Estimate.ID == "" {
		return nil, errors.New("estimate is required for delivery")
	}

	payment, err := u.billingServiceRepo.GetPaymentByEstimateID(ctx, serviceOrder.Estimate.ID)
	if err != nil {
		return nil, err
	}

	if payment == nil || payment.ID == "" {
		return nil, errors.New("payment information is required for delivery")
	}

	serviceOrder.Status = valueobject.StatusEntregue

	err = u.repo.Update(ctx, serviceOrder)
	if err != nil {
		logger.Error().Err(err).Msg("Error updating service order")
		return nil, err
	}

	return serviceOrder, nil
}

func (u *ServiceOrderUseCase) PaymentServiceOrder(ctx context.Context, serviceOrderID string) (*entities.ServiceOrder, error) {
	logger := logs.LoggerWithContext(ctx)
	if txn := newrelic.FromContext(ctx); txn != nil {
		startSegment := txn.StartSegment("ServiceOrderUseCase.validateExecution")
		defer startSegment.End()
	}

	serviceOrder, err := u.validateServiceOrderExists(ctx, serviceOrderID)
	if err != nil {
		logger.Error().Err(err).Msg("Error finding service order with id")
		return nil, err
	}

	if !serviceOrder.Status.IsFinalizada() {
		return nil, ErrInvalidTransitionStatusToDelivery
	}

	if serviceOrder.Estimate == nil || serviceOrder.Estimate.ID == "" {
		return nil, errors.New("estimate is required for delivery")
	}

	if _, err := u.billingServiceRepo.CreatePayment(ctx, serviceOrder.Estimate.ID); err != nil {
		return nil, err
	}

	return serviceOrder, nil
}

func (u *ServiceOrderUseCase) CancelServiceOrder(ctx context.Context, serviceOrderID string) (*entities.ServiceOrder, error) {
	logger := logs.LoggerWithContext(ctx)
	if txn := newrelic.FromContext(ctx); txn != nil {
		startSegment := txn.StartSegment("ServiceOrderUseCase.validateExecution")
		defer startSegment.End()
	}

	serviceOrder, err := u.validateServiceOrderExists(ctx, serviceOrderID)
	if err != nil {
		logger.Error().Err(err).Msg("Error finding service order with id")
		return nil, err
	}

	if serviceOrder.Status.IsFinalizada() || serviceOrder.Status.IsEntregue() || serviceOrder.Status.IsCancelada() {
		return nil, ErrInvalidTransitionStatusToCancel
	}

	serviceOrder.Status = valueobject.StatusCancelada

	err = u.repo.Update(ctx, serviceOrder)
	if err != nil {
		logger.Error().Err(err).Msg("Error updating service order")
		return nil, err
	}
	return serviceOrder, nil
}

func (u *ServiceOrderUseCase) validateVehicle(ctx context.Context, serviceOrder *entities.ServiceOrder) error {
	logger := logs.LoggerWithContext(ctx)
	if serviceOrder.VehicleID != "" {
		vehicle, err := u.vehicleRepo.FindByID(serviceOrder.VehicleID)
		if err != nil {
			logger.Error().Err(err).Any("vehicle_id", serviceOrder.VehicleID).Msg("error finding vehicle with id")
			return err
		}
		if vehicle == nil || vehicle.ID == "" {
			return ErrVehicleNotFound
		}
		return nil
	}
	return ErrInvalidID
}

func (u *ServiceOrderUseCase) validateCustomer(ctx context.Context, serviceOrder *entities.ServiceOrder) error {
	logger := logs.LoggerWithContext(ctx)
	if serviceOrder.CustomerID != "" {
		customer, err := u.customerRepo.GetByID(serviceOrder.CustomerID)
		if err != nil {
			logger.Error().Err(err).Any("customer_id", serviceOrder.CustomerID).Msg("error finding customer with id")
			return err
		}
		if customer == nil || customer.ID == "" {
			return ErrCustomerNotFound
		}
		return nil
	}
	return ErrInvalidID
}

func (u *ServiceOrderUseCase) validateServicesExists(ctx context.Context, services []entities.Service) error {
	for _, s := range services {
		_, err := u.serviceRepo.GetByID(ctx, s.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *ServiceOrderUseCase) validateServiceOrderExists(ctx context.Context, id string) (*entities.ServiceOrder, error) {
	logger := logs.LoggerWithContext(ctx)
	serviceOrder, err := u.repo.GetByID(ctx, id, false)
	if err != nil {
		logger.Error().Err(err).Any("OS_ID", id).Msg("Error finding service order with id")
		return nil, err
	}

	if serviceOrder == nil || serviceOrder.ID == "" {
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
