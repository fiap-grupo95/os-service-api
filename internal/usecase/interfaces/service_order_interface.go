package interfaces

import (
	"context"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
)

type IServiceOrderGateway interface {
	Create(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error)
	GetByID(ctx context.Context, id string, isFullData bool) (*entities.ServiceOrder, error)
	Update(ctx context.Context, serviceOrder *entities.ServiceOrder) error
	List(ctx context.Context) ([]*entities.ServiceOrder, error)
}

type IServiceOrderRepository interface {
	Create(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error)
	GetByID(ctx context.Context, id string) (*entities.ServiceOrder, error)
	Update(ctx context.Context, serviceOrder *entities.ServiceOrder) error
	List(ctx context.Context) ([]*entities.ServiceOrder, error)
}
