package service

import (
	"github.com/rudransh/distributed-commerce/payment/internal/model"
	"github.com/rudransh/distributed-commerce/payment/internal/provider"
	"github.com/rudransh/distributed-commerce/payment/internal/repository"
)

func (s *paymentService) createAttempt(repo repository.PaymentAttemptRepository,payment *model.Payment,event provider.WebhookEvent) (*model.PaymentAttempt, error) {

	attemptNumber := 1

	lastAttempt, err := repo.GetLatestAttempt(payment.ID)

	if err != nil {
		return nil, err
	}

	if lastAttempt != nil {
		attemptNumber = lastAttempt.AttemptNumber + 1
	}

	attempt := &model.PaymentAttempt{
		PaymentID:       payment.ID,
		AttemptNumber:   attemptNumber,
		Status:          event.Status,
		GatewayResponse: event.EventID,
	}

	return attempt, nil
}