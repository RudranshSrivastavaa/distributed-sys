package state

import (
	"errors"

	"github.com/rudransh/distributed-commerce/inventory/internal/model"
)

var transitions = map[model.ReservationStatus][]model.ReservationStatus{

	model.StatusReserved: {
		model.StatusConfirmed,
		model.StatusReleased,
		model.StatusExpired,
	},

	model.StatusConfirmed: {},

	model.StatusReleased: {},

	model.StatusExpired: {},
}

func CanTransition(
	current model.ReservationStatus,
	next model.ReservationStatus,
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
	reservation *model.Reservation,
	next model.ReservationStatus,
) error {

	if !CanTransition(
		reservation.Status,
		next,
	) {

		return errors.New(
			"invalid reservation state transition",
		)

	}

	reservation.Status = next

	return nil

}