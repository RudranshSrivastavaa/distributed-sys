package service

import (
	"errors"
	"strings"

	"github.com/rudransh/distributed-commerce/inventory/internal/dto"
)

func validateCreateProduct(request dto.CreateProductRequest) error {

	if strings.TrimSpace(request.SKU) == "" {
		return errors.New("sku is required")
	}

	if strings.TrimSpace(request.Name) == "" {
		return errors.New("name is required")
	}

	if request.Price <= 0 {
		return errors.New("price must be greater than zero")
	}

	return nil

}

func validateStockQuantity(quantity int64) error {

	if quantity <= 0 {
		return errors.New(
			"quantity must be greater than zero",
		)
	}

	return nil
}