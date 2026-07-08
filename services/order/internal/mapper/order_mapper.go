package mapper

import (
	"github.com/rudransh/distributed-commerce/order/internal/dto"
	"github.com/rudransh/distributed-commerce/order/internal/model"
)

func ToOrder(request dto.CreateOrderRequest) *model.Order {

	order := &model.Order{
		CustomerID: request.CustomerID,
	}

	items := make([]model.OrderItem, 0, len(request.Items))

	for _, item := range request.Items {

		items = append(items, model.OrderItem{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Price:       item.Price,
		})

	}

	order.Items = items

	return order
}

func ToOrderResponse(order *model.Order) dto.OrderResponse {

	response := dto.OrderResponse{
		ID:          order.ID,
		CustomerID:  order.CustomerID,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}

	items := make([]dto.OrderItemResponse, 0, len(order.Items))

	for _, item := range order.Items {

		items = append(items, dto.OrderItemResponse{
			ID:          item.ID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Price:       item.Price,
		})

	}

	response.Items = items

	return response

}

func ToOrderResponses(
	orders []model.Order,
) []dto.OrderResponse {

	response := make(
		[]dto.OrderResponse,
		0,
		len(orders),
	)

	for i := range orders {

		response = append(
			response,
			ToOrderResponse(&orders[i]),
		)

	}

	return response

}

func UpdateOrder(
	order *model.Order,
	request dto.UpdateOrderRequest,
) {

	items := make([]model.OrderItem, 0, len(request.Items))

	for _, item := range request.Items {

		items = append(items, model.OrderItem{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Price:       item.Price,
		})

	}

	order.Items = items
}