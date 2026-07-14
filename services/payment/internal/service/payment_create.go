package service

import (
	"context"

	"github.com/rudransh/distributed-commerce/payment/event"
	"github.com/rudransh/distributed-commerce/payment/internal/dto"
	"github.com/rudransh/distributed-commerce/payment/internal/model"
	"github.com/rudransh/distributed-commerce/payment/internal/outbox"
	"github.com/rudransh/distributed-commerce/payment/internal/repository"
	"github.com/rudransh/distributed-commerce/payment/internal/state"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
	"gorm.io/gorm"
)

func (s *paymentService) CreatePayment(
	ctx context.Context,
	request dto.CreatePaymentRequest,
) (dto.PaymentResponse, error) {

	//----------------------------------------------------
	// Validate
	//----------------------------------------------------

	if err := s.validateCreatePayment(request); err != nil {
		return dto.PaymentResponse{}, err
	}

	//----------------------------------------------------
	// Money
	//----------------------------------------------------

	money, err := model.NewMoney(
		request.Amount,
		request.Currency,
	)
	if err != nil {
		return dto.PaymentResponse{}, err
	}

	//----------------------------------------------------
	// Create Payment
	//----------------------------------------------------

	payment := &model.Payment{
		OrderID:  request.OrderID,
		Money:    money,
		Status:   model.StatusCreated,
		Provider: request.Provider,
	}

	//----------------------------------------------------
	// Transaction 1
	// Create payment + PAYMENT_CREATED Outbox
	//----------------------------------------------------

	err = s.transactionManager.Execute(func(tx *gorm.DB) error {

		repo := repository.NewPaymentRepository(tx)

		outboxRepo := repository.NewOutboxRepository(tx)

		publisher := outbox.NewOutboxPublisher(outboxRepo)

		if err := repo.Create(payment); err != nil {
			return err
		}

		evt, err := event.BuildPaymentCreatedEvent(payment)
		if err != nil {
			return err
		}

		if err := publisher.Publish(
			ctx,
			kafkaa.PaymentEvents,
			payment.ID.String(),
			evt,
		); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return dto.PaymentResponse{}, err
	}

	//----------------------------------------------------
	// Gateway Call
	//----------------------------------------------------

	intent, err := s.createIntent(payment)
	if err != nil {
		return dto.PaymentResponse{}, err
	}

	//----------------------------------------------------
	// Move Payment -> PENDING
	//----------------------------------------------------

	payment.ProviderReference = intent.ProviderReference

	if err := state.Transition(
		payment,
		model.StatusPending,
	); err != nil {
		return dto.PaymentResponse{}, err
	}

	//----------------------------------------------------
	// Transaction 2
	// Only update payment
	//----------------------------------------------------

	err = s.transactionManager.Execute(func(tx *gorm.DB) error {

		repo := repository.NewPaymentRepository(tx)

		if err := repo.Update(payment); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return dto.PaymentResponse{}, err
	}

	//----------------------------------------------------
	// Response
	//----------------------------------------------------

	return s.toPaymentResponse(
		payment,
		nil,
		intent.PaymentURL,
	), nil
}