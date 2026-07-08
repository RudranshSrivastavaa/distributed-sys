package webhook

import (
	"crypto/hmac"
	"encoding/json"

	"github.com/rudransh/distributed-commerce/payment/internal/dto"
	"github.com/rudransh/distributed-commerce/payment/internal/provider"
)

type Verifier struct {
	secret string
}

func NewVerifier(secret string) *Verifier {
	return &Verifier{
		secret: secret,
	}
}

func (v *Verifier) Verify(
	body []byte,
	signature string,
) bool {

	expected := provider.GenerateSignature(
		body,
		v.secret,
	)

	return hmac.Equal(
		[]byte(signature),
		[]byte(expected),
	)

}

func PayloadForVerification(
	request dto.WebhookRequest,
) ([]byte, error) {

	payload := struct {
		EventID string `json:"eventId"`

		Provider string `json:"provider"`

		ProviderReference string `json:"providerReference"`

		Status string `json:"status"`
	}{
		EventID: request.EventID,
		Provider: request.Provider,
		ProviderReference: request.ProviderReference,
		Status: request.Status,
	}

	return json.Marshal(payload)
}