package service

import (
	"github.com/google/uuid"
	orderErrors "github.com/rudransh/distributed-commerce/order/internal/errors"
	"github.com/rudransh/distributed-commerce/order/internal/model"
)

func ValidateOrder(
	order *model.Order,
) error {

	if order.CustomerID == uuid.Nil {
		return orderErrors.ErrCustomerIDRequired
	}

	if len(order.Items) == 0 {
		return orderErrors.ErrEmptyOrder
	}

	for _, item := range order.Items {

		if item.Quantity <= 0 {
			return orderErrors.ErrInvalidQuantity
		}

		if item.Price <= 0 {
			return orderErrors.ErrInvalidPrice
		}

	}

	return nil
}