package service

import (
	"log"

	"github.com/rudransh/distributed-commerce/notification/internal/model"
	"github.com/rudransh/distributed-commerce/notification/internal/repository"
	"gorm.io/gorm"
)


func (s *notificationService) moveToDLQ(tx *gorm.DB,notification *model.Notification) error {

	repository := repository.NewDeadLetterRepository(tx)

	entry := &model.DeadLetterNotification{

		NotificationID: notification.ID,

		EventID: notification.EventID,

		Recipient: notification.Recipient,

		Subject: notification.Subject,

		Body: notification.Body,

		Channel: notification.Channel,

		FailureReason: notification.FailureReason,

		RetryCount: s.retryExecutor.Config().MaxAttempts,
	}
	log.Println("Moving notification to DLQ")
	
	return repository.Create(entry)

}