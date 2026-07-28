package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/rudransh/distributed-commerce/payment/internal/model"
)

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(
	db *gorm.DB,
) PaymentRepository {

	return &paymentRepository{
		db: db,
	}
}

func (r *paymentRepository) Create(payment *model.Payment) error {

	return r.db.Create(payment).Error

}

func (r *paymentRepository) FindByID(id uuid.UUID) (*model.Payment, error) {

	var payment model.Payment

	err := r.db.First(&payment, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &payment, nil

}

func (r *paymentRepository) FindByOrderID(orderID uuid.UUID) (*model.Payment, error) {

	var payment model.Payment

	err := r.db.Where("order_id = ?", orderID).First(&payment).Error

	if err != nil {
		return nil, err
	}

	return &payment, nil

}

func (r *paymentRepository) Update(payment *model.Payment) error {

	return r.db.Save(payment).Error

}

func (r *paymentRepository) FindByProviderReference(reference string) (*model.Payment, error) {

	var payment model.Payment

	err := r.db.
		Where("provider_reference = ?",reference).First(&payment).Error

	if err != nil {
		return nil, err
	}

	return &payment, nil
}