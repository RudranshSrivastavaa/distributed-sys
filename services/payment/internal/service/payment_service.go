package service

import (
	"github.com/google/uuid"

	"github.com/rudransh/distributed-commerce/pkg/circuitbreaker"
	"github.com/rudransh/distributed-commerce/payment/internal/dto"
	"github.com/rudransh/distributed-commerce/payment/internal/provider"
	"github.com/rudransh/distributed-commerce/payment/internal/repository"
	"github.com/rudransh/distributed-commerce/pkg/retry"
)

type PaymentService interface {

	CreatePayment(
		request dto.CreatePaymentRequest,
	) (dto.PaymentResponse, error)

	ProcessPayment(
		paymentID uuid.UUID,
		request dto.ProcessPaymentRequest,
	) (dto.PaymentResponse, error)

	GetPayment(
		id uuid.UUID,
	) (dto.PaymentResponse, error)

	GetPaymentByOrderID(
		orderID uuid.UUID,
	) (dto.PaymentResponse, error)

	HandleWebhook(
	event provider.WebhookEvent,
) error
}

type paymentService struct {
	paymentRepository repository.PaymentRepository

	attemptRepository repository.PaymentAttemptRepository

	transactionManager *repository.TransactionManager

	gateway provider.PaymentGateway

	retryExecutor *retry.Executor

	breaker *circuitbreaker.Breaker
}

func NewPaymentService(
	paymentRepository repository.PaymentRepository,
	attemptRepository repository.PaymentAttemptRepository,
	transactionManager *repository.TransactionManager,
	gateway provider.PaymentGateway,
	retryExecutor *retry.Executor,
	breaker *circuitbreaker.Breaker,
) PaymentService {

	return &paymentService{
		paymentRepository: paymentRepository,
		attemptRepository: attemptRepository,
		transactionManager: transactionManager,
		gateway:            gateway,
		retryExecutor: retryExecutor,
		breaker: breaker,
	}
}