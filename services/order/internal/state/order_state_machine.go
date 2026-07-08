package state

import (
    orderErrors "github.com/rudransh/distributed-commerce/order/internal/errors"
	"github.com/rudransh/distributed-commerce/order/internal/model"
)

var transitions = map[model.OrderStatus][]model.OrderStatus{

	model.StatusCreated: {
		model.StatusPendingPayment,
		model.StatusCancelled,
	},

	model.StatusPendingPayment: {
		model.StatusPaid,
		model.StatusCancelled,
	},

	model.StatusPaid: {
		model.StatusReserved,
	},

	model.StatusReserved: {
		model.StatusConfirmed,
	},

	model.StatusConfirmed: {
		model.StatusShipped,
	},

	model.StatusShipped: {
		model.StatusDelivered,
	},
}


func CanTransition(
	current model.OrderStatus,
	next model.OrderStatus,
) bool {

	allowed := transitions[current]

	for _, status := range allowed {

		if status == next {
			return true
		}

	}

	return false
}


func Transition(
	order *model.Order,
	next model.OrderStatus,
) error {

	if !CanTransition(
		order.Status,
		next,
	) {

		return orderErrors.ErrInvalidStatusTransition

	}

	order.Status = next

	return nil
}