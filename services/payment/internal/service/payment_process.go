package service

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/rudransh/distributed-commerce/payment/internal/dto"
	"github.com/rudransh/distributed-commerce/payment/internal/model"
	sagaevent "github.com/rudransh/distributed-commerce/saga/event"
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

	// Payment will be updated when the webhook arrives.

	return s.toPaymentResponse(
		payment,
		nil,
		"",
	), nil
}

func (s *paymentService) HandleProcessPaymentCommand(
	ctx context.Context,
	request sagaevent.ProcessPaymentPayload,
) error {

	orderID, err := uuid.Parse(request.OrderID)
	if err != nil {
		return err
	}

	log.Println("Creating payment...")

	payment, err := s.CreatePayment(
		ctx,
		dto.CreatePaymentRequest{
			OrderID:  orderID,
			Amount:   request.Amount,
			Currency: request.Currency,
			Provider: model.ProviderSimulator,
		},
	)
	if err != nil {
		return err
	}
	log.Println("Payment created")
	log.Println("Calling payment gateway...")
	_, err = s.ProcessPayment(
		payment.ID,
		dto.ProcessPaymentRequest{
			ForceFailure: false, //true for failing the payment and false for payment done
		},
	)

	return err
}
