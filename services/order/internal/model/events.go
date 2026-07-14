package model

import (
	orderevents "github.com/rudransh/distributed-commerce/pkg/events/order"
)

func (o *Order) CreatedEvent() orderevents.OrderCreatedPayload {

	items := make([]orderevents.OrderItem, 0, len(o.Items))

	for _, item := range o.Items {

		items = append(items, orderevents.OrderItem{

			ProductID: item.ProductID,

			Quantity: item.Quantity,
		})
	}

	return orderevents.OrderCreatedPayload{

    OrderID: o.ID.String(),

    CustomerID: o.CustomerID,

    TotalPrice: o.TotalAmount,

    Currency: "INR",

    Items: items,
}
}