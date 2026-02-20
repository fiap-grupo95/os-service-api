package request

import (
	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
)

// PaymentCreateRequest represents the payload to create a payment.
type PaymentCreateRequest struct {
	MpPayload PaymentMpPayloadRequest `json:"mp_payload" binding:"required"`
}

type PaymentMpPayloadRequest struct {
	PaymentMethodID string              `json:"payment_method_id" binding:"required"`
	Payer           PaymentPayerRequest `json:"payer" binding:"required"`
}

type PaymentPayerRequest struct {
	Email string `json:"email" binding:"required"`
}

// ToEntity converts the request into a Payment entity.
func (r PaymentCreateRequest) ToEntity() (entities.Payment, error) {
	return entities.Payment{
		PaymentMethodID: r.MpPayload.PaymentMethodID,
		Payer: entities.Payer{
			Email: r.MpPayload.Payer.Email,
		},
	}, nil
}
