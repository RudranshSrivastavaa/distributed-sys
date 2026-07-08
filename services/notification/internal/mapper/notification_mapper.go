package mapper

import (
	"github.com/rudransh/distributed-commerce/notification/internal/dto"
	"github.com/rudransh/distributed-commerce/notification/internal/model"
)

func ToNotificationResponse(notification *model.Notification) dto.NotificationResponse {

	return dto.NotificationResponse{

		ID: notification.ID.String(),

		EventID: notification.EventID,

		Recipient: notification.Recipient,

		Subject: notification.Subject,

		Body: notification.Body,

		Channel: string(notification.Channel),

		Status: string(notification.Status),

		FailureReason: notification.FailureReason,

		CreatedAt: notification.CreatedAt,

		UpdatedAt: notification.UpdatedAt,
	}

}