package provider

import (
	"github.com/rudransh/distributed-commerce/payment/internal/dto"
	"github.com/rudransh/distributed-commerce/payment/internal/model"
)

type PaymentGateway interface {
	CreateIntent(payment *model.Payment) (*CreateIntentResult, error)

	Capture(payment *model.Payment,request dto.ProcessPaymentRequest) (*PaymentResult, error)
}
