package mapper

import (
	"github.com/rudransh/distributed-commerce/inventory/internal/dto"
	"github.com/rudransh/distributed-commerce/inventory/internal/model"
)

func ToProduct(request dto.CreateProductRequest) *model.Product {

	return &model.Product{
		SKU:         request.SKU,
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
	}

}

func ToProductResponse(product *model.Product) dto.ProductResponse {

	return dto.ProductResponse{
		ID:          product.ID,
		SKU:         product.SKU,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}

}