package interfaces

import (
	"context"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
)

type IPaymentRepo interface {
	Create(ctx context.Context, payment entities.Payment) (entities.Payment, error)
	GetByID(ctx context.Context, id string) (entities.Payment, error)
	GetByServiceOrderID(ctx context.Context, serviceOrderID string) (entities.Payment, error)
}
