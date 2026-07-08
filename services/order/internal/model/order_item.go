package model

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type OrderItem struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	OrderID uuid.UUID `gorm:"type:uuid;not null;index" json:"orderId"`

	ProductID uuid.UUID `gorm:"type:uuid;not null" json:"productId"`

	ProductName string `gorm:"size:255;not null" json:"productName"`

	Quantity int `gorm:"not null" json:"quantity"`

	Price float64 `gorm:"type:numeric(12,2);not null" json:"price"`

	CreatedAt time.Time `json:"createdAt"`

	UpdatedAt time.Time `json:"updatedAt"`
}

func (i *OrderItem) BeforeCreate(tx *gorm.DB) error {
	i.ID = uuid.New()
	return nil
}