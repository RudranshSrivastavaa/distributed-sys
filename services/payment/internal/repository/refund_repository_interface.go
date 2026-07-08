package repository

import (
	"github.com/google/uuid"

	"github.com/rudransh/distributed-commerce/payment/internal/model"
)

type RefundRepository interface {

	Create(
		refund *model.Refund,
	) error

	Update(
		refund *model.Refund,
	) error

	FindByID(
		id uuid.UUID,
	) (*model.Refund, error)

	FindByPaymentID(
		paymentID uuid.UUID,
	) ([]model.Refund, error)
}