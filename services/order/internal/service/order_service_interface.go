package service

import (
	"context"

	"github.com/google/uuid"

	//ordercommand "github.com/rudransh/distributed-commerce/order/command"
	"github.com/rudransh/distributed-commerce/order/internal/dto"
	sagaevent "github.com/rudransh/distributed-commerce/saga/event"
)

type OrderService interface {
	Create(request dto.CreateOrderRequest) (dto.OrderResponse, bool, error)

	GetAll() ([]dto.OrderResponse, error)

	GetByID(id uuid.UUID) (dto.OrderResponse, error)

	Update(id uuid.UUID, request dto.UpdateOrderRequest) (dto.OrderResponse, error)

	Delete(id uuid.UUID) error

	UpdateStatus(id uuid.UUID, request dto.UpdateOrderStatusRequest) (dto.OrderResponse, error)

	HandleCompleteOrder(
		ctx context.Context,
		request sagaevent.CompleteOrderPayload,
	) error

	HandleCancelOrderCommand(
		ctx context.Context,
		request sagaevent.CancelOrderPayload,
	) error
}
