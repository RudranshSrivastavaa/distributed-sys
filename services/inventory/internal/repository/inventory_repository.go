package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/rudransh/distributed-commerce/inventory/internal/errors"
	"github.com/rudransh/distributed-commerce/inventory/internal/model"
)

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(
	db *gorm.DB,
) InventoryRepository {

	return &inventoryRepository{
		db: db,
	}
}

func (r *inventoryRepository) Create(
	inventory *model.Inventory,
) error {

	return r.db.Create(inventory).Error

}

func (r *inventoryRepository) FindByProductID(
	productID uuid.UUID,
) (*model.Inventory, error) {

	var inventory model.Inventory

	err := r.db.
		Where("product_id = ?", productID).
		First(&inventory).
		Error

	if err != nil {
		return nil, err
	}

	return &inventory, nil

}

func (r *inventoryRepository) Update(
	inventory *model.Inventory,
) error {

	return r.db.Save(inventory).Error

}
func (r *inventoryRepository) UpdateWithVersion(
	inventory *model.Inventory,
) error {

	result := r.db.Model(&model.Inventory{}).
		Where(
			"product_id = ? AND version = ?",
			inventory.ProductID,
			inventory.Version,
		).
		Updates(map[string]interface{}{
			"available_quantity": inventory.AvailableQuantity,
			"reserved_quantity":  inventory.ReservedQuantity,
			"version":            inventory.Version + 1,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.ErrOptimisticLockConflict
	}

	inventory.Version++

	return nil
}