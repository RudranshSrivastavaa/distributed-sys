package provider

import "github.com/rudransh/distributed-commerce/payment/internal/model"

type WebhookEvent struct {

	EventID string

	Provider model.PaymentProvider

	ProviderReference string

	Status model.PaymentStatus

	Signature string

}

type webhookPayload struct {
	EventID           string `json:"eventId"`
	Provider          string `json:"provider"`
	ProviderReference string `json:"providerReference"`
	Status            string `json:"status"`
}


func (e WebhookEvent) IsSuccessful() bool {

	return e.Status == model.StatusSuccess

}

func (e WebhookEvent) IsFailed() bool {

	return e.Status == model.StatusFailed

}

func (e WebhookEvent) IsValid() bool {

	if e.EventID == "" {
		return false
	}

	if !e.Provider.IsSupported() {
		return false
	}

	if e.ProviderReference == "" {
		return false
	}

	switch e.Status {

	case model.StatusSuccess,
		model.StatusFailed:

	default:
		return false
	}

	return true

}