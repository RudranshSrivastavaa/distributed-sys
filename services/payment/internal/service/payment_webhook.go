package service

import (
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/rudransh/distributed-commerce/payment/internal/model"
	"github.com/rudransh/distributed-commerce/payment/internal/provider"
	"github.com/rudransh/distributed-commerce/payment/internal/repository"
	"github.com/rudransh/distributed-commerce/payment/internal/state"
)

func (s *paymentService) HandleWebhook(
	event provider.WebhookEvent,
) error {
	log.Println("========== WEBHOOK RECEIVED ==========")
	log.Printf("EventID: %s", event.EventID)
	log.Printf("ProviderReference: %s", event.ProviderReference)

	return s.retryExecutor.Do(func() error {
		return s.transactionManager.Execute(
			func(tx *gorm.DB) error {

				paymentRepo := repository.NewPaymentRepository(tx)

				attemptRepo := repository.NewPaymentAttemptRepository(tx)

				webhookRepo := repository.NewWebhookEventRepository(tx)

				//----------------------------------------
				// Duplicate webhook?
				//----------------------------------------

				duplicate, err := s.isDuplicateWebhook(
					webhookRepo,
					event.EventID,
				)

				if err != nil {
					return err
				}

				if duplicate {
					log.Printf(
						"duplicate webhook ignored: %s",
						event.EventID,
					)
					return nil
				}
				log.Println("1. Duplicate check")

				//----------------------------------------
				// Save webhook event
				//----------------------------------------

				webhookEvent := &model.WebhookEvent{
					EventID:           event.EventID,
					Provider:          event.Provider,
					ProviderReference: event.ProviderReference,
					ProcessedAt:       time.Now(),
				}

				if err := webhookRepo.Create(
					webhookEvent,
				); err != nil {
					return err
				}
				log.Println("2. Saving webhook event")
				//----------------------------------------
				// Find payment
				//----------------------------------------

				payment, err := paymentRepo.FindByProviderReference(
					event.ProviderReference,
				)

				if err != nil {
					return err
				}
				log.Println("3. Finding payment")
				//----------------------------------------
				// Already processed?
				//----------------------------------------

				if payment.Status == event.Status {
					log.Printf(
						"duplicate webhook ignored: %s",
						event.EventID,
					)
					return nil
				}

				//----------------------------------------
				// Transition state
				//----------------------------------------

				if err := state.Transition(
					payment,
					event.Status,
				); err != nil {
					return err
				}
				log.Println("4. Transition payment")
				//----------------------------------------
				// Create payment attempt
				//----------------------------------------

				attempt, err := s.createAttempt(
					attemptRepo,
					payment,
					event,
				)

				if err != nil {
					return err
				}

				if err := attemptRepo.Create(
					attempt,
				); err != nil {
					return err
				}
				log.Println("5. Create attempt")
				//----------------------------------------
				// Save payment
				//----------------------------------------

				if err := paymentRepo.Update(
					payment,
				); err != nil {
					return err
				}

				return nil
			},
		)
	})
}
