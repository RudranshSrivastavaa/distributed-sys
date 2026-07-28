package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/rudransh/distributed-commerce/payment/internal/model"
)

type paymentAttemptRepository struct {
	db *gorm.DB
}

func NewPaymentAttemptRepository(
	db *gorm.DB,
) PaymentAttemptRepository {

	return &paymentAttemptRepository{
		db: db,
	}
}

func (r *paymentAttemptRepository) Create(attempt *model.PaymentAttempt) error {

	return r.db.Create(attempt).Error

}

func (r *paymentAttemptRepository) FindByPaymentID(paymentID uuid.UUID) ([]model.PaymentAttempt, error) {

	var attempts []model.PaymentAttempt

	err := r.db.
		Where("payment_id = ?", paymentID).
		Order("created_at ASC").
		Find(&attempts).
		Error

	if err != nil {
		return nil, err
	}

	return attempts, nil

}

func (r *paymentAttemptRepository) GetLatestAttempt(paymentID uuid.UUID) (*model.PaymentAttempt, error) {

	var attempt model.PaymentAttempt

	err := r.db.
		Where("payment_id = ?",paymentID).
		Order("attempt_number DESC").
		First(&attempt).
		Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &attempt, nil
}