package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/interfaces"
)

var (
	ErrorPaymentNotFound         = errors.New("payment not found")
	ErrPaymentAlreadyExists      = errors.New("payment already exists")
	ErrPaymentAmountDoesNotMatch = errors.New("payment amount does not match service order estimate")
)

type IPaymentUseCase interface {
	CreatePayment(ctx context.Context, payment *entities.Payment) (*entities.Payment, error)
	GetPaymentByID(ctx context.Context, id uint) (*entities.Payment, error)
}

type PaymentUseCase struct {
	repo             interfaces.IPaymentRepo
	serviceOrderRepo interfaces.IServiceOrderGateway
}

var _ IPaymentUseCase = (*PaymentUseCase)(nil)

func NewPaymentUseCase(repo interfaces.IPaymentRepo, serviceOrderRepo interfaces.IServiceOrderGateway) *PaymentUseCase {
	return &PaymentUseCase{
		repo:             repo,
		serviceOrderRepo: serviceOrderRepo,
	}
}

func (p *PaymentUseCase) CreatePayment(ctx context.Context, payment *entities.Payment) (*entities.Payment, error) {
	serviceOrder, err := p.serviceOrderRepo.GetByID(ctx, payment.ServiceOrderID, false)
	if err != nil {
		return nil, err
	}
	if serviceOrder.Estimate == nil {
		return nil, ErrPaymentAmountDoesNotMatch
	}
	if serviceOrder.Estimate.Value != payment.Amount {
		return nil, ErrPaymentAmountDoesNotMatch
	}
	existingPayment, err := p.repo.GetByServiceOrderID(ctx, payment.ServiceOrderID)
	if err != nil {
		return nil, err
	}
	if existingPayment.ID != 0 {
		return nil, ErrPaymentAlreadyExists
	}

	if payment.PaymentDate.IsZero() {
		payment.PaymentDate = time.Now()
	}

	created, err := p.repo.Create(ctx, *payment)
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func (p *PaymentUseCase) GetPaymentByID(ctx context.Context, id uint) (*entities.Payment, error) {
	payment, err := p.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if payment.ID == 0 {
		return nil, ErrorPaymentNotFound
	}
	return &payment, nil
}