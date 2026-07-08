package repository

import (
	stdErrors "errors"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/rudransh/distributed-commerce/order/internal/model"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func NewOrderRepositoryWithDB(db *gorm.DB) OrderRepository {

	return &orderRepository{
		db: db,
	}

}

func (r *orderRepository) Create(order *model.Order) error {
	if r.db == nil {
		panic("repository db is nil")
	}
	return r.db.Create(order).Error

}

func (r *orderRepository) FindAll() ([]model.Order, error) {

	var orders []model.Order

	err := r.db.
		Preload("Items").
		Find(&orders).Error

	if err != nil {
		return nil, err
	}

	return orders, nil

}

func (r *orderRepository) FindByID(id uuid.UUID) (*model.Order, error) {

	var order model.Order

	err := r.db.
		Preload("Items").
		First(&order, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &order, nil

}

func (r *orderRepository) Update(order *model.Order) error {

	return r.db.Save(order).Error

}

func (r *orderRepository) Delete(id uuid.UUID) error {

	return r.db.Delete(
		&model.Order{},
		"id = ?",
		id,
	).Error

}

func (r *orderRepository) FindByIdempotencyKey(
    key string,
) (*model.IdempotencyKey, error) {

    var record model.IdempotencyKey

    err := r.db.Where("key = ?", key).First(&record).Error

    if stdErrors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }

    if err != nil {
        return nil, err
    }

    return &record, nil
}
func (r *orderRepository) SaveIdempotencyKey(record *model.IdempotencyKey) error {
	return r.db.Create(record).Error
}

func (r *orderRepository) WithTransaction(
	fn func(tx *gorm.DB) error,
) error {

	return r.db.Transaction(func(tx *gorm.DB) error {

		return fn(tx)

	})

}
