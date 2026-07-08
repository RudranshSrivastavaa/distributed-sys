package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Inventory struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	ProductID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"productId"`

	Product Product `gorm:"foreignKey:ProductID"`

	AvailableQuantity int64 `gorm:"not null" json:"availableQuantity"`

	ReservedQuantity int64 `gorm:"not null" json:"reservedQuantity"`

	Version int64 `gorm:"not null;default:1" json:"version"`

	CreatedAt time.Time `json:"createdAt"`

	UpdatedAt time.Time `json:"updatedAt"`
}

func (i *Inventory) BeforeCreate(tx *gorm.DB) error {
	i.ID = uuid.New()
	return nil
}

func (i *Inventory) AddStock(quantity int64) {
	i.AvailableQuantity += quantity
}

func (i *Inventory) Reserve(quantity int64) error {

	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	if i.AvailableQuantity < quantity {
		return errors.New("insufficient stock")
	}

	i.AvailableQuantity -= quantity
	i.ReservedQuantity += quantity

	return nil
}

func (i *Inventory) Release(quantity int64) error {

	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	if i.ReservedQuantity < quantity {
		return errors.New("invalid reserved quantity")
	}

	i.ReservedQuantity -= quantity
	i.AvailableQuantity += quantity

	return nil
}

func (i *Inventory) Confirm(quantity int64) error {

	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	if i.ReservedQuantity < quantity {
		return errors.New("invalid reserved quantity")
	}

	i.ReservedQuantity -= quantity

	return nil
}