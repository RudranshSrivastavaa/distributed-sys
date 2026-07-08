package repository

import "github.com/rudransh/distributed-commerce/notification/internal/model"

type InboxRepository interface {

	Exists(eventID string) (bool, error)

	Create(message *model.InboxMessage) error
}