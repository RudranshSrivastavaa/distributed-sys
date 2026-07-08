package service

import (
	"github.com/google/uuid"
	"github.com/rudransh/distributed-commerce/notification/internal/dto"
	"github.com/rudransh/distributed-commerce/notification/internal/mapper"
	"github.com/rudransh/distributed-commerce/notification/internal/model"
	"github.com/rudransh/distributed-commerce/notification/internal/provider"
	"github.com/rudransh/distributed-commerce/notification/internal/repository"
	"github.com/rudransh/distributed-commerce/notification/internal/state"
	"github.com/rudransh/distributed-commerce/pkg/circuitbreaker"
	"github.com/rudransh/distributed-commerce/pkg/retry"
	"gorm.io/gorm"
)

type notificationService struct {
	repository         repository.NotificationRepository
	transactionManager *repository.TransactionManager
	emailProvider      provider.EmailProvider
	retryExecutor      *retry.Executor
	breaker            *circuitbreaker.Breaker
}

func NewNotificationService(
	repository repository.NotificationRepository,
	transactionManager *repository.TransactionManager,
	emailProvider provider.EmailProvider,
	retryExecutor *retry.Executor,
	breaker *circuitbreaker.Breaker) NotificationService {

	return &notificationService{
		repository:         repository,
		transactionManager: transactionManager,
		emailProvider:      emailProvider,
		retryExecutor:      retryExecutor,
		breaker:            breaker,
	}
}

func (s *notificationService) CreateNotification(
	request dto.CreateNotificationRequest,
) (dto.NotificationResponse, error) {

	if err := validateCreateNotification(request); err != nil {
		return dto.NotificationResponse{}, err
	}

	//----------------------------------------------------
	// Create Notification
	//----------------------------------------------------

	notification := &model.Notification{

		EventID: request.EventID,

		Recipient: request.Recipient,

		Subject: request.Subject,

		Body: request.Body,

		Channel: model.NotificationChannel(request.Channel),

		Status: model.StatusPending,
	}

	//----------------------------------------------------
	// Transaction
	//----------------------------------------------------

	err := s.transactionManager.Execute(
		func(tx *gorm.DB) error {

			repo := repository.NewNotificationRepository(tx)

			//----------------------------------------
			// Inbox Pattern
			//----------------------------------------

			if err := s.processInbox(
				tx,
				request.EventID,
				"NotificationCreated",
				"notification-service",
			); err != nil {
				return err
			}

			//----------------------------------------
			// Save Notification
			//----------------------------------------

			if err := repo.Create(notification); err != nil {
				return err
			}

			//----------------------------------------
			// Send Notification
			//----------------------------------------

			return s.sendNotification(
				tx,
				repo,
				notification,
			)
		},
	)

	if err != nil {
		return dto.NotificationResponse{}, err
	}

	//----------------------------------------------------
	// Response
	//----------------------------------------------------

	return mapper.ToNotificationResponse(
		notification,
	), nil
}

func (s *notificationService) ListNotifications() ([]dto.NotificationResponse, error) {

	notifications, err := s.repository.List()

	if err != nil {
		return nil, err
	}

	response := make([]dto.NotificationResponse, 0, len(notifications))

	for i := range notifications {
		response = append(
			response,
			mapper.ToNotificationResponse(&notifications[i]),
		)
	}

	return response, nil
}

func (s *notificationService) DeleteNotification(id uuid.UUID) error {

	return s.repository.Delete(id)

}

func (s *notificationService) GetNotification(
	id uuid.UUID,
) (dto.NotificationResponse, error) {

	notification, err := s.repository.FindByID(id)
	if err != nil {
		return dto.NotificationResponse{}, err
	}

	return mapper.ToNotificationResponse(notification), nil
}

func (s *notificationService) sendNotification(tx *gorm.DB, repo repository.NotificationRepository,
	notification *model.Notification,
) error {

	if err := state.Transition(
		notification,
		model.StatusProcessing,
	); err != nil {
		return err
	}

	if err := repo.Update(notification); err != nil {
		return err
	}

	err := s.breaker.Execute(func() error {

		return s.retryExecutor.Do(func() error {

			return s.emailProvider.Send(notification)

		})

	})

	if err != nil {

		notification.FailureReason = err.Error()

		if err := state.Transition(
			notification,
			model.StatusFailed,
		); err != nil {
			return err
		}

		if err := repo.Update(notification); err != nil {
			return err
		}

		return s.moveToDLQ(
			tx,
			notification,
		)

	}

	if err := state.Transition(
		notification,
		model.StatusSent,
	); err != nil {
		return err
	}

	return repo.Update(notification)

}
