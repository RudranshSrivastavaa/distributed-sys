package service

import (
	"errors"

	"github.com/rudransh/distributed-commerce/notification/internal/dto"
	"github.com/rudransh/distributed-commerce/notification/internal/model"
)

func validateCreateNotification(
	request dto.CreateNotificationRequest,
) error {

	if request.Recipient == "" {
		return errors.New("recipient is required")
	}

	if request.Subject == "" {
		return errors.New("subject is required")
	}

	if request.Body == "" {
		return errors.New("body is required")
	}

	channel := model.NotificationChannel(request.Channel)

	if !channel.IsSupported() {
		return errors.New("unsupported notification channel")
	}

	return nil
}