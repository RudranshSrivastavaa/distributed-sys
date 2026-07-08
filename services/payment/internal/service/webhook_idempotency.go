package service

import "github.com/rudransh/distributed-commerce/payment/internal/repository"

func (s *paymentService) isDuplicateWebhook(repo repository.WebhookEventRepository,eventID string) (bool, error) {

	return repo.Exists(eventID)

}