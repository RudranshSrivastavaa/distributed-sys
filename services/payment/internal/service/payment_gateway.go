package service

import (
	"github.com/rudransh/distributed-commerce/payment/internal/model"
	"github.com/rudransh/distributed-commerce/payment/internal/provider"
)

func (s *paymentService) createIntent(payment *model.Payment) (*provider.CreateIntentResult, error) {

	return s.gateway.CreateIntent(payment)
}