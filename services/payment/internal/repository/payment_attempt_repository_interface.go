package repository

import (
	"github.com/google/uuid"

	"github.com/rudransh/distributed-commerce/payment/internal/model"
)

type PaymentAttemptRepository interface {

	Create(
		attempt *model.PaymentAttempt,
	) error

	FindByPaymentID(
		paymentID uuid.UUID,
	) ([]model.PaymentAttempt, error)

	GetLatestAttempt(
	paymentID uuid.UUID,
) (*model.PaymentAttempt, error)

}