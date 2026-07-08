package provider

import "github.com/rudransh/distributed-commerce/notification/internal/model"

type EmailProvider interface {
	Send(notification *model.Notification) error
}