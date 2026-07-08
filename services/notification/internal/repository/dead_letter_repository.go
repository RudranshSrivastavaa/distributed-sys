package repository

import (
	"github.com/rudransh/distributed-commerce/notification/internal/model"
	"gorm.io/gorm"
)

type deadLetterRepository struct {
	db *gorm.DB
}

func NewDeadLetterRepository(db *gorm.DB) DeadLetterRepository {

	return &deadLetterRepository{
		db: db,
	}

}

func (r *deadLetterRepository) Create(
	notification *model.DeadLetterNotification,
) error {

	return r.db.Create(notification).Error

}

func (r *deadLetterRepository) List() (

	[]model.DeadLetterNotification,

	error,
) {

	var notifications []model.DeadLetterNotification

	err := r.db.
		Order("created_at DESC").
		Find(&notifications).
		Error

	if err != nil {
		return nil, err
	}

	return notifications, nil

}