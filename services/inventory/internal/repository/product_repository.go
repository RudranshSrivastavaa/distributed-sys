package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/rudransh/distributed-commerce/inventory/internal/model"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) Create(product *model.Product) error {

	return r.db.Create(product).Error

}

func (r *productRepository) FindByID(id uuid.UUID) (*model.Product, error) {

	var product model.Product

	err := r.db.First(&product, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &product, nil

}

func (r *productRepository) FindBySKU(sku string) (*model.Product, error) {

	var product model.Product

	err := r.db.Where("sku = ?", sku).First(&product).Error

	if err != nil {
		return nil, err
	}

	return &product, nil

}

func (r *productRepository) FindAll() ([]model.Product, error) {

	var products []model.Product

	err := r.db.Find(&products).Error

	if err != nil {
		return nil, err
	}

	return products, nil

}

func (r *productRepository) Update(product *model.Product) error {

	return r.db.Save(product).Error

}

func (r *productRepository) Delete(id uuid.UUID) error {

	return r.db.Delete(&model.Product{},"id = ?",id).Error

}