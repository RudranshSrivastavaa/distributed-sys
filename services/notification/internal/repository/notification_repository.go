package repository

import (
	"github.com/google/uuid"
	"github.com/rudransh/distributed-commerce/notification/internal/model"
	"gorm.io/gorm"
)

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(
	db *gorm.DB,
) NotificationRepository {

	return &notificationRepository{
		db: db,
	}

}

func (r *notificationRepository) Create(notification *model.Notification) error {

	return r.db.Create(notification).Error

}

func (r *notificationRepository) FindByID(id uuid.UUID) (*model.Notification, error) {

	var notification model.Notification

	err := r.db.
		First(
			&notification,
			"id = ?",
			id,
		).
		Error

	if err != nil {
		return nil, err
	}

	return &notification, nil

}

func (r *notificationRepository) Update(notification *model.Notification) error {

	return r.db.Save(notification).Error

}

func (r *notificationRepository) List() ([]model.Notification,error) {

	var notifications []model.Notification

	err := r.db.
		Order("created_at DESC").
		Find(&notifications).
		Error

	if err != nil {
		return nil, err
	}

	return notifications, nil

}

func (r *notificationRepository) Delete(id uuid.UUID) error {

	return r.db.Delete(&model.Notification{},"id = ?",id).Error
}