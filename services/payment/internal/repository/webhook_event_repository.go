package repository

import (
	"github.com/rudransh/distributed-commerce/payment/internal/model"
	"gorm.io/gorm"
)

type webhookEventRepository struct {
	db *gorm.DB
}

func NewWebhookEventRepository(
	db *gorm.DB,
) WebhookEventRepository {

	return &webhookEventRepository{
		db: db,
	}

}

func (r *webhookEventRepository) Exists(eventID string) (bool, error) {

	var count int64

	err := r.db.Model(&model.WebhookEvent{}).Where("event_id = ?",eventID).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *webhookEventRepository) Create(event *model.WebhookEvent) error {

	return r.db.Create(event).Error
}