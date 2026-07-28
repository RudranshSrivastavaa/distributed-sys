package provider

import (
	"bytes"
	"encoding/json"

	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/rudransh/distributed-commerce/payment/internal/dto"
	"github.com/rudransh/distributed-commerce/payment/internal/model"

)

type SimulatorGateway struct {
	Config SimulatorConfig

	Client *http.Client

	WebhookURL string

}

func NewSimulatorGateway(config SimulatorConfig,webhookURL string) *SimulatorGateway {

	return &SimulatorGateway{
		Config: config,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
		WebhookURL: webhookURL,
	}
}

func (s *SimulatorGateway) CreateIntent(payment *model.Payment) (*CreateIntentResult, error) {

	reference := uuid.NewString()

	// Your roadmap currently does this.
	// Later we'll move this to Capture().
	//go s.sendWebhook(reference)

	return &CreateIntentResult{

		ProviderReference: reference,

		PaymentURL: "/payments/simulator/" + reference,
	}, nil
}

func (s *SimulatorGateway) ProcessPayment(intentID string) (*PaymentResult, error) {

	return &PaymentResult{

		Success: true,

		Status: model.StatusSuccess,

		Reference: uuid.NewString(),

		Message: "payment successful",
	}, nil
}

func (s *SimulatorGateway) Capture(payment *model.Payment,request dto.ProcessPaymentRequest) (*PaymentResult, error) {
	
	if request.ForceFailure {

		go s.sendWebhook(
			payment.ProviderReference,
			"FAILED",
		)

		return &PaymentResult{

			Success: false,

			Status: model.StatusFailed,

			Message: "simulated payment failure",
		}, nil
	}

	go s.sendWebhook(
		payment.ProviderReference,
		"SUCCESS",
	)

	return &PaymentResult{

		Success: true,

		Status: model.StatusSuccess,

		Reference: uuid.NewString(),

		Message: "payment captured successfully",
	}, nil
}

func (s *SimulatorGateway) sendWebhook(providerReference string ,status string,) {

	time.Sleep(10 * time.Second)

	//----------------------------------------
	// Payload used for signature
	//----------------------------------------

	payload := webhookPayload{

		EventID: uuid.NewString(),

		Provider: "SIMULATOR",

		ProviderReference: providerReference,

		Status: status,
	}

	body, err := json.Marshal(payload)

	if err != nil {

		log.Println(err)

		return
	}

	//----------------------------------------
	// Generate Signature
	//----------------------------------------

	signature := GenerateSignature(body,s.Config.WebhookSecret)

	//----------------------------------------
	// Actual webhook request
	//----------------------------------------

	request := dto.WebhookRequest{
		EventID: payload.EventID,
		Provider: payload.Provider,
		ProviderReference: payload.ProviderReference,
		Status: payload.Status,
		Signature: signature,
	}

	requestBody, err := json.Marshal(request)

	if err != nil {
		log.Println(err)
		return
	}

	//----------------------------------------
	// Send webhook
	//----------------------------------------

	resp, err := s.Client.Post(s.WebhookURL,
		"application/json",
		bytes.NewReader(requestBody),
	)

	if err != nil {
		log.Println("Webhook POST failed:", err)
		return
	}

	log.Println("Webhook response:", resp.Status)
	log.Println("Webhook JSON:")
	log.Println(string(requestBody))

	if err != nil {
		log.Println(err)
	}
}
