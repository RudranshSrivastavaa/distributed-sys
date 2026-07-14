package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rudransh/distributed-commerce/inventory/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type outboxRepository struct {
	db *gorm.DB
}

func NewOutboxRepository(
	db *gorm.DB,
) OutboxRepository {

	return &outboxRepository{
		db: db,
	}
}

func NewOutboxRepositoryWithDB(
	db *gorm.DB,
) OutboxRepository {

	return &outboxRepository{
		db: db,
	}
}

func (r *outboxRepository) Create(ctx context.Context,event *model.OutboxEvent) error {

	return r.db.WithContext(ctx).Create(event).Error

}

func (r *outboxRepository) LockPending(ctx context.Context,limit int) ([]model.OutboxEvent, error) {

	var events []model.OutboxEvent

	err := r.db.
		WithContext(ctx).
		Clauses(clause.Locking{
			Strength: "UPDATE",
			Options:  "SKIP LOCKED",
		}).
		Where("status = ?", model.OutboxPending).
		Order("created_at ASC").
		Limit(limit).
		Find(&events).
		Error

	return events, err
}

func (r *outboxRepository) MarkPublished(ctx context.Context,id uuid.UUID) error {

	now := time.Now().UTC()

	return r.db.
		WithContext(ctx).
		Model(&model.OutboxEvent{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":       model.OutboxPublished,
			"published_at": now,
		}).
		Error
}

func (r *outboxRepository) IncrementRetry(ctx context.Context,id uuid.UUID,nextRetry time.Time,
	lastError string) error {

	return r.db.
		WithContext(ctx).
		Model(&model.OutboxEvent{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"retry_count": gorm.Expr("retry_count + 1"),
			"next_retry_at": nextRetry,
			"last_error":  lastError,
		}).
		Error
}

func (r *outboxRepository) MarkFailed(ctx context.Context,id uuid.UUID,lastError string) error {

	return r.db.
		WithContext(ctx).
		Model(&model.OutboxEvent{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":     model.OutboxFailed,
			"last_error": lastError,
		}).
		Error
}