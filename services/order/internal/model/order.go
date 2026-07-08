package model

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type Order struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	CustomerID uuid.UUID `gorm:"type:uuid;not null;index" json:"customerId"`

	TotalAmount float64 `gorm:"type:numeric(12,2);not null" json:"totalAmount"`

	Status OrderStatus `gorm:"type:varchar(30);not null" json:"status"`

	Items []OrderItem `gorm:"foreignKey:OrderID" json:"items"`

	CreatedAt time.Time `json:"createdAt"`

	UpdatedAt time.Time `json:"updatedAt"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	o.ID = uuid.New()
	return nil
}