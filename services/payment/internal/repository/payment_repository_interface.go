package repository

import (
	"github.com/google/uuid"

	"github.com/rudransh/distributed-commerce/payment/internal/model"
)

type PaymentRepository interface {

	Create(payment *model.Payment) error

	FindByID(id uuid.UUID) (*model.Payment, error)

	FindByOrderID(orderID uuid.UUID) (*model.Payment, error)

	Update(payment *model.Payment) error

	FindByProviderReference(reference string) (*model.Payment, error)

}