package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rudransh/distributed-commerce/payment/internal/dto"
)

func (s *paymentService) ProcessPayment(
	paymentID uuid.UUID,
	request dto.ProcessPaymentRequest,
) (dto.PaymentResponse, error) {

	payment, err := s.paymentRepository.FindByID(paymentID)
	if err != nil {
		return dto.PaymentResponse{}, err
	}

	if payment.Status.IsFinal() {
		return dto.PaymentResponse{},
			errors.New("payment already completed")
	}



err = s.breaker.Execute(func() error {

	return s.retryExecutor.Do(func() error {

		var err error

		_, err = s.gateway.Capture(
			payment,
			request,
		)

		return err
	})

})

if err != nil {
	return dto.PaymentResponse{}, err
}

	if err != nil {
		return dto.PaymentResponse{}, err
	}

	// Payment will be updated when the webhook arrives.

	return s.toPaymentResponse(
		payment,
		nil,
		"",
	), nil
}
