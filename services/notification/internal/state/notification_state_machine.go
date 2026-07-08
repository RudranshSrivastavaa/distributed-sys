package state

import (
	"errors"

	"github.com/rudransh/distributed-commerce/notification/internal/model"
)

var transitions = map[model.NotificationStatus][]model.NotificationStatus{

	model.StatusPending: {

		model.StatusProcessing,
	},

	model.StatusProcessing: {

		model.StatusSent,

		model.StatusFailed,
	},

	model.StatusFailed: {

		model.StatusProcessing,
	},

	model.StatusSent: {},
}

func CanTransition(

	current model.NotificationStatus,

	next model.NotificationStatus,

) bool {

	allowed := transitions[current]

	for _, state := range allowed {

		if state == next {
			return true
		}

	}

	return false

}

func Transition(

	notification *model.Notification,

	next model.NotificationStatus,

) error {

	if !CanTransition(

		notification.Status,

		next,

	) {

		return errors.New(
			"invalid notification state transition",
		)

	}

	notification.Status = next

	return nil

}