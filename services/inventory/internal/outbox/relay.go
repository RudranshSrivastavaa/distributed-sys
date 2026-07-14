package outbox

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/rudransh/distributed-commerce/inventory/internal/model"
	"github.com/rudransh/distributed-commerce/inventory/internal/repository"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
)

type Relay interface {
	Start(ctx context.Context) error
	Stop() error
}

type relay struct {
	config     Config
	repository repository.OutboxRepository
	producer   kafkaa.Producer
	stop       chan struct{}
}

func NewRelay(cfg Config, repository repository.OutboxRepository, producer kafkaa.Producer) Relay {

	return &relay{
		config:     cfg,
		repository: repository,
		producer:   producer,
		stop:       make(chan struct{}),
	}
}

func (r *relay) Start(ctx context.Context) error {

	ticker := time.NewTicker(
		r.config.PollInterval,
	)

	defer ticker.Stop()

	for {
		select {

		case <-ctx.Done():
			return nil

		case <-r.stop:
			return nil

		case <-ticker.C:
			if err := r.processBatch(ctx); err != nil {

				log.Error(
					"relay batch failed",
					"error",
					err,
				)
			}
		}
	}
}

func (r *relay) processBatch(ctx context.Context) error {

	events, err := r.repository.LockPending(
		ctx,
		r.config.BatchSize,
	)

	if err != nil {
		return err
	}

	if len(events) == 0 {
		time.Sleep(500 * time.Millisecond)
	}
	log.Infof("Found %d pending outbox events", len(events))
	for _, event := range events {

		r.processEvent(
			ctx,
			event,
		)
	}

	return nil
}
func (r *relay) processEvent(ctx context.Context, event model.OutboxEvent) {

	var kafkaEvent kafkaa.Event

	err := json.Unmarshal(event.Payload, &kafkaEvent)

	if err != nil {

		r.repository.MarkFailed(
			ctx,
			event.ID,
			err.Error(),
		)

		return
	}

	log.Infof("Publishing outbox event %s",event.EventID)

	err = r.producer.Publish(
		ctx,
		event.Topic,
		kafkaEvent.Metadata.AggregateID,
		kafkaEvent,
	)

	if err == nil {
		r.repository.MarkPublished(
			ctx,
			event.ID,
		)
		log.Infof(
			"Published outbox event %s",
			event.EventID,
		)
		return
	}
	if event.RetryCount >= r.config.MaxRetries {

		r.repository.MarkFailed(
			ctx,
			event.ID,
			err.Error(),
		)

		return
	}
	next := CalculateNextRetry(

		event.RetryCount,
	)

	r.repository.IncrementRetry(
		ctx,
		event.ID,
		next,
		err.Error(),
	)

}

func (r *relay) Stop() error {

	close(r.stop)

	return nil
}
