package service

import (
	"time"

	"github.com/rudransh/distributed-commerce/notification/internal/errors"
	"github.com/rudransh/distributed-commerce/notification/internal/model"
	"github.com/rudransh/distributed-commerce/notification/internal/repository"
	"gorm.io/gorm"
)

func (s *notificationService) processInbox(
	tx *gorm.DB,
	eventID string,
	eventType string,
	source string,
) error {

	repository := repository.NewInboxRepository(tx)

	exists, err := repository.Exists(
		eventID,
	)

	if err != nil {
		return err
	}

	if exists {

		return errors.ErrDuplicateEvent;

	}

	message := &model.InboxMessage{

		EventID: eventID,

		EventType: eventType,

		Source: source,

		ProcessedAt: time.Now(),
	}

	return repository.Create(
		message,
	)

}