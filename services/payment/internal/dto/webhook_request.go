package dto

type WebhookRequest struct {
	EventID string `json:"eventId"`

	Provider string `json:"provider"`

	ProviderReference string `json:"providerReference"`

	Status string `json:"status"`

	Signature string `json:"signature"`
}