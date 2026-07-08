package state

import (
	"errors"

	"github.com/rudransh/distributed-commerce/payment/internal/model"
)

var paymentTransitions = map[model.PaymentStatus][]model.PaymentStatus{

	model.StatusCreated: {
		model.StatusPending,
	},

	model.StatusPending: {
		model.StatusSuccess,
		model.StatusFailed,
	},

	model.StatusSuccess: {
		model.StatusRefunded,
	},

	model.StatusFailed: {},

	model.StatusRefunded: {},
}

func CanTransition(
	current model.PaymentStatus,
	next model.PaymentStatus,
) bool {

	allowed := paymentTransitions[current]

	for _, status := range allowed {

		if status == next {
			return true
		}

	}

	return false

}

func Transition(
	payment *model.Payment,
	next model.PaymentStatus,
) error {

	if !CanTransition(
		payment.Status,
		next,
	) {

		return errors.New(
			"invalid payment state transition",
		)

	}

	payment.Status = next

	return nil

}