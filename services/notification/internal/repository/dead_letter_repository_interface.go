package repository

import (
	"github.com/rudransh/distributed-commerce/notification/internal/model"
)


type DeadLetterRepository interface {

	Create(notification *model.DeadLetterNotification) error

	List() ([]model.DeadLetterNotification,error)
}