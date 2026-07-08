package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	SKU string `gorm:"size:100;uniqueIndex;not null" json:"sku"`

	Name string `gorm:"size:255;not null" json:"name"`

	Description string `gorm:"type:text" json:"description"`

	Price int64 `gorm:"not null" json:"price"`

	CreatedAt time.Time `json:"createdAt"`

	UpdatedAt time.Time `json:"updatedAt"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.New()
	return nil
}