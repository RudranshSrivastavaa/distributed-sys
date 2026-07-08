package service

import (
	"github.com/rudransh/distributed-commerce/payment/internal/dto"
	"github.com/rudransh/distributed-commerce/payment/internal/model"
	"github.com/rudransh/distributed-commerce/payment/internal/repository"
	"github.com/rudransh/distributed-commerce/payment/internal/state"
	"gorm.io/gorm"
)

func (s *paymentService) CreatePayment(
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
		OrderID: request.OrderID,

		Money: money,

		Status: model.StatusCreated,

		Provider: request.Provider,
	}

	//----------------------------------------------------
	// First Transaction
	//----------------------------------------------------

	err = s.transactionManager.Execute(func(tx *gorm.DB) error {

		repo := repository.NewPaymentRepository(tx)

		return repo.Create(payment)
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
	// Update Payment
	//----------------------------------------------------

	payment.ProviderReference = intent.ProviderReference

	if err := state.Transition(
		payment,
		model.StatusPending,
	); err != nil {
		return dto.PaymentResponse{}, err
	}

	err = s.transactionManager.Execute(func(tx *gorm.DB) error {

		repo := repository.NewPaymentRepository(tx)

		return repo.Update(payment)
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
