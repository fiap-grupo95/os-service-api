package interfaces

import (
	"context"

	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/request"
	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
)

type IBillingServiceGateway interface {
	CreateEstimate(ctx context.Context, serviceOrder *entities.ServiceOrder, additionalRepair *entities.AdditionalRepair) (*entities.Estimate, error)
	ApproveEstimate(ctx context.Context, serviceOrder *entities.ServiceOrder, additionalRepair *entities.AdditionalRepair) (*entities.Estimate, error)
	RejectEstimate(ctx context.Context, serviceOrder *entities.ServiceOrder, additionalRepair *entities.AdditionalRepair) (*entities.Estimate, error)
	CancelEstimate(ctx context.Context, serviceOrder *entities.ServiceOrder, additionalRepair *entities.AdditionalRepair) (*entities.Estimate, error)
	GetPaymentByEstimateID(ctx context.Context, estimateID string) (*entities.Payment, error)
	CreatePayment(ctx context.Context, estimateID string) (*entities.Payment, error)
}

type IBillingServiceRepository interface {
	CreateEstimate(ctx context.Context, request *request.EstimateRequest) (*response.EstimateResponse, error)
	ApproveEstimate(ctx context.Context, request *request.EstimateRequest) (*response.EstimateResponse, error)
	RejectEstimate(ctx context.Context, request *request.EstimateRequest) (*response.EstimateResponse, error)
	CancelEstimate(ctx context.Context, request *request.EstimateRequest) (*response.EstimateResponse, error)
	GetPaymentByEstimateID(ctx context.Context, estimateID string) (*response.PaymentResponse, error)
	CreatePayment(ctx context.Context, estimateID string) (*response.PaymentResponse, error)
}
