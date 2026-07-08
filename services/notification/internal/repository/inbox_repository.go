package repository

import (
	"github.com/rudransh/distributed-commerce/notification/internal/model"
	"gorm.io/gorm"
)

type inboxRepository struct {
	db *gorm.DB
}

func NewInboxRepository(
	db *gorm.DB,
) InboxRepository {

	return &inboxRepository{
		db: db,
	}
}

func (r *inboxRepository) Exists(
	eventID string,
) (bool, error) {

	var count int64

	err := r.db.
		Model(&model.InboxMessage{}).
		Where(
			"event_id = ?",
			eventID,
		).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *inboxRepository) Create(
	message *model.InboxMessage,
) error {

	return r.db.Create(message).Error
}