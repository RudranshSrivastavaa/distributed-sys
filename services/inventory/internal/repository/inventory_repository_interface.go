package repository

import (
	"github.com/google/uuid"

	"github.com/rudransh/distributed-commerce/inventory/internal/model"
)

type InventoryRepository interface {

	Create(inventory *model.Inventory) error

	FindByProductID(productID uuid.UUID) (*model.Inventory, error)

	Update(inventory *model.Inventory) error

	UpdateWithVersion(inventory *model.Inventory) error
}