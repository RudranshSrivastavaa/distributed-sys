package service

import (
	"log"

	"github.com/rudransh/distributed-commerce/order/internal/model"
)

func CalculateTotal(order *model.Order) {

	var total float64

	for _, item := range order.Items {

		total += item.Price * float64(item.Quantity)

	}
	log.Printf("Calculated Total = %f", total)
	order.TotalAmount = total

}