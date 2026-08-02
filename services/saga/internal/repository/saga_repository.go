package repository

import (
	"github.com/rudransh/distributed-commerce/saga/internal/model"
)

type SagaRepository interface {

	Create(saga *model.Saga) error

	FindByOrderID(orderID string) (*model.Saga, error)

	Update(saga *model.Saga) error
}