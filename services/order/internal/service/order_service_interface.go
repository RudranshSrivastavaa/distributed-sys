package service

import (
	"github.com/google/uuid"

	"github.com/rudransh/distributed-commerce/order/internal/dto"
)

type OrderService interface {

	Create(request dto.CreateOrderRequest) (dto.OrderResponse, bool ,error)

	GetAll() ([]dto.OrderResponse, error)

	GetByID(id uuid.UUID) (dto.OrderResponse, error)

	Update(id uuid.UUID,request dto.UpdateOrderRequest) (dto.OrderResponse, error)

	Delete(id uuid.UUID) error

	UpdateStatus(id uuid.UUID,request dto.UpdateOrderStatusRequest) (dto.OrderResponse, error)
}