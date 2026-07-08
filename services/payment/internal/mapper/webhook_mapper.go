package mapper

import (
	"errors"

	"github.com/rudransh/distributed-commerce/payment/internal/dto"
	"github.com/rudransh/distributed-commerce/payment/internal/model"
	"github.com/rudransh/distributed-commerce/payment/internal/provider"
)

func ToWebhookEvent(
	request dto.WebhookRequest,
) (provider.WebhookEvent, error) {

	event := provider.WebhookEvent{
		EventID: request.EventID,

		Provider: model.PaymentProvider(request.Provider),

		ProviderReference: request.ProviderReference,

		Status: model.PaymentStatus(request.Status),

		Signature: request.Signature,
	}

	if !event.IsValid() {
		return provider.WebhookEvent{},
			errors.New("invalid webhook event")
	}

	return event, nil

}