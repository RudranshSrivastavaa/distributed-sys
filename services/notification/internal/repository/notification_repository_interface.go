package repository

import (
	"github.com/google/uuid"

	"github.com/rudransh/distributed-commerce/notification/internal/model"
)

type NotificationRepository interface {

	Create(notification *model.Notification) error

	FindByID(id uuid.UUID) (*model.Notification, error)

	Update(notification *model.Notification) error

	Delete(id uuid.UUID) error

	List() ([]model.Notification, error)
}