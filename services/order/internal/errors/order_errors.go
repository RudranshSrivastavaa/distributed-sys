package errors

import "errors"

var (

	ErrOrderNotFound = errors.New(
		"order not found",
	)

	ErrCustomerIDRequired = errors.New(
		"customer id is required",
	)

	ErrEmptyOrder = errors.New(
		"order must contain at least one item",
	)

	ErrInvalidQuantity = errors.New(
		"quantity must be greater than zero",
	)

	ErrInvalidPrice = errors.New(
		"price must be greater than zero",
	)

	ErrDeliveredOrderCannotBeDeleted = errors.New(
		"delivered orders cannot be deleted",
	)

	ErrInvalidStatusTransition = errors.New(
		"invalid order status transition",
	)

)