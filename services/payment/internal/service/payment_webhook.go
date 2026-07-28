package service

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm"

	paymentevent "github.com/rudransh/distributed-commerce/payment/event"
	"github.com/rudransh/distributed-commerce/payment/internal/model"
	"github.com/rudransh/distributed-commerce/payment/internal/outbox"
	"github.com/rudransh/distributed-commerce/payment/internal/provider"
	"github.com/rudransh/distributed-commerce/payment/internal/repository"
	"github.com/rudransh/distributed-commerce/payment/internal/state"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
)

func (s *paymentService) HandleWebhook(event provider.WebhookEvent) error {

	log.Println("========== WEBHOOK RECEIVED ==========")
	log.Printf("EventID: %s", event.EventID)
	log.Printf("ProviderReference: %s", event.ProviderReference)

	return s.retryExecutor.Do(func() error {

		return s.transactionManager.Execute(func(tx *gorm.DB) error {

			paymentRepo := repository.NewPaymentRepository(tx)

			attemptRepo := repository.NewPaymentAttemptRepository(tx)

			webhookRepo := repository.NewWebhookEventRepository(tx)

			outboxRepo := repository.NewOutboxRepository(tx)

			publisher := outbox.NewOutboxPublisher(outboxRepo)

			//----------------------------------------------------
			// Duplicate webhook?
			//----------------------------------------------------

			duplicate, err := s.isDuplicateWebhook(webhookRepo,event.EventID)

			if err != nil {
				return err
			}

			if duplicate {
				log.Printf("duplicate webhook ignored: %s",event.EventID)
				return nil
			}

			log.Println("1. Duplicate checked")

			//----------------------------------------------------
			// Save webhook
			//----------------------------------------------------

			webhookEvent := &model.WebhookEvent{
				EventID:           event.EventID,
				Provider:          event.Provider,
				ProviderReference: event.ProviderReference,
				ProcessedAt:       time.Now(),
			}

			if err := webhookRepo.Create(webhookEvent); err != nil {
				return err
			}

			log.Println("2. Saving webhook event")

			//----------------------------------------------------
			// Find payment
			//----------------------------------------------------

			payment, err := paymentRepo.FindByProviderReference(event.ProviderReference)

			if err != nil {
				return err
			}

			log.Println("3. Finding payment")

			//----------------------------------------------------
			// Already processed?
			//----------------------------------------------------

			if payment.Status == event.Status {
				log.Printf("duplicate webhook ignored: %s",event.EventID)
				return nil
			}

			//----------------------------------------------------
			// Transition
			//----------------------------------------------------

			if err := state.Transition(payment,event.Status); err != nil {
				return err
			}

			log.Println("4. Transition payment")

			//----------------------------------------------------
			// Create Attempt
			//----------------------------------------------------

			attempt, err := s.createAttempt(attemptRepo,payment,event)

			if err != nil {
				return err
			}

			if err := attemptRepo.Create(attempt); err != nil {
				return err
			}

			log.Println("5. Create attempt done")

			//----------------------------------------------------
			// Update Payment
			//----------------------------------------------------

			if err := paymentRepo.Update(payment); err != nil {
				return err
			}

			//----------------------------------------------------
			// Publish Outbox Event
			//----------------------------------------------------

			var evt kafkaa.Event

			switch payment.Status {

			case model.StatusSuccess:

				evt, err = paymentevent.BuildPaymentCompletedEvent(payment)

				if err != nil {
					return err
				}

			case model.StatusFailed:

				evt, err = paymentevent.BuildPaymentFailedEvent(payment)

				if err != nil {
					return err
				}

			default:
				return nil
			}

			if err := publisher.Publish(
				context.Background(),
				kafkaa.PaymentEvents,
				payment.ID.String(),
				evt,
			); err != nil {
				return err
			}

			log.Println("6. Payment event written to Outbox")

			return nil
		})
	})
}