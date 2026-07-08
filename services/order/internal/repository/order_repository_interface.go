package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/rudransh/distributed-commerce/order/internal/model"
)

type OrderRepository interface {

	Create(order *model.Order) error

	FindAll() ([]model.Order, error)

	FindByID(id uuid.UUID) (*model.Order, error)

	Update(order *model.Order) error

	Delete(id uuid.UUID) error

	FindByIdempotencyKey(key string) (*model.IdempotencyKey, error)

    SaveIdempotencyKey(record *model.IdempotencyKey) error
	
	WithTransaction(
        fn func(tx *gorm.DB) error,
    ) error

}