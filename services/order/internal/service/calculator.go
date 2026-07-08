package service

import "github.com/rudransh/distributed-commerce/order/internal/model"

func CalculateTotal(order *model.Order) {

	var total float64

	for _, item := range order.Items {

		total += item.Price * float64(item.Quantity)

	}

	order.TotalAmount = total

}