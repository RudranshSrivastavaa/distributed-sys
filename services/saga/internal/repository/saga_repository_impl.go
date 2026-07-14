package repository

import (
	"github.com/rudransh/distributed-commerce/saga/internal/model"
	"gorm.io/gorm"
)

type sagaRepository struct {
	db *gorm.DB
}

func NewSagaRepository(db *gorm.DB) SagaRepository {

	return &sagaRepository{
		db: db,
	}
}

func (r *sagaRepository) Create(saga *model.Saga) error {

	return r.db.Create(saga).Error
}

func (r *sagaRepository) FindByOrderID(orderID string) (*model.Saga, error) {

	var saga model.Saga

	err := r.db.Where("order_id = ?", orderID).First(&saga).Error

	if err != nil {
		return nil, err
	}

	return &saga, nil
}

func (r *sagaRepository) Update(saga *model.Saga) error {

	return r.db.Save(saga).Error
}