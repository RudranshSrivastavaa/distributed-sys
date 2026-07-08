package dto

import "github.com/rudransh/distributed-commerce/order/internal/model"

type UpdateOrderStatusRequest struct {
	Status model.OrderStatus `json:"status"`
}