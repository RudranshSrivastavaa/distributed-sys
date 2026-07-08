package repository

import (
	"github.com/google/uuid"

	"github.com/rudransh/distributed-commerce/inventory/internal/model"
)

type ProductRepository interface {

	Create(product *model.Product) error

	FindAll() ([]model.Product, error)

	FindByID(id uuid.UUID) (*model.Product, error)

	FindBySKU(sku string) (*model.Product, error)

	Update(product *model.Product) error

	Delete(id uuid.UUID) error
}