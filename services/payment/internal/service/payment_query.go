package service

import (
	"github.com/google/uuid"
	"github.com/rudransh/distributed-commerce/payment/internal/dto"
)

func (s *paymentService) GetPayment(id uuid.UUID) (dto.PaymentResponse, error) {

	payment, err := s.paymentRepository.FindByID(id)

	if err != nil {
		return dto.PaymentResponse{}, err
	}

	attempts, err := s.attemptRepository.FindByPaymentID(payment.ID)

	if err != nil {
		return dto.PaymentResponse{}, err
	}

	return s.toPaymentResponse(payment,attempts,""), nil

}

func (s *paymentService) GetPaymentByOrderID(orderID uuid.UUID) (dto.PaymentResponse, error) {

	payment, err := s.paymentRepository.FindByOrderID(orderID)

	if err != nil {
		return dto.PaymentResponse{}, err
	}

	attempts, err := s.attemptRepository.FindByPaymentID(payment.ID)

	if err != nil {
		return dto.PaymentResponse{}, err
	}

	return s.toPaymentResponse(payment,attempts,""), nil

}