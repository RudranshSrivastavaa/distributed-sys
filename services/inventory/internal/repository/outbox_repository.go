package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rudransh/distributed-commerce/inventory/internal/model"
)

type OutboxRepository interface {

	Create(
		ctx context.Context,
		event *model.OutboxEvent,
	) error

	LockPending(
		ctx context.Context,
		limit int,
	) ([]model.OutboxEvent, error)

	MarkPublished(
		ctx context.Context,
		id uuid.UUID,
	) error

	IncrementRetry(
		ctx context.Context,
		id uuid.UUID,
		nextRetry time.Time,
		lastError string,
	) error

	MarkFailed(
		ctx context.Context,
		id uuid.UUID,
		lastError string,
	) error
}