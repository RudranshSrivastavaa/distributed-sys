package repository

import (

	"github.com/rudransh/distributed-commerce/payment/internal/model"
)

type WebhookEventRepository interface {

	Create(event *model.WebhookEvent) error

	Exists(eventID string) (bool, error)
}