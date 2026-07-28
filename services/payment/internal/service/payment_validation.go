package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rudransh/distributed-commerce/payment/internal/dto"
)

func (s *paymentService) validateCreatePayment(request dto.CreatePaymentRequest) error {

	if request.OrderID == uuid.Nil {
		return errors.New("order id is required")
	}

	if request.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	if request.Currency == "" {
		return errors.New("currency is required")
	}

	if !request.Provider.IsSupported() {
		return errors.New("unsupported payment provider")
	}

	return nil
}