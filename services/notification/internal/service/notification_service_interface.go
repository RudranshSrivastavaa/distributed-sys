package service

import (
	"github.com/google/uuid"

	"github.com/rudransh/distributed-commerce/notification/internal/dto"
)

type NotificationService interface {

	CreateNotification(request dto.CreateNotificationRequest) (dto.NotificationResponse, error)

	GetNotification(id uuid.UUID) (dto.NotificationResponse, error)

	ListNotifications() ([]dto.NotificationResponse,error)

	DeleteNotification(id uuid.UUID) error

}