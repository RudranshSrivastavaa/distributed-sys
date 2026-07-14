package service

import (
	"context"

	inventoryevents "github.com/rudransh/distributed-commerce/pkg/events/inventory"
	orderevents "github.com/rudransh/distributed-commerce/pkg/events/order"
	paymentevents "github.com/rudransh/distributed-commerce/pkg/events/payment"
)

type SagaService interface {
	StartSaga(
		ctx context.Context,
		payload orderevents.OrderCreatedPayload,
	) error

	HandleInventoryReserved(
		ctx context.Context,
		payload inventoryevents.InventoryReservedPayload,
	) error

	HandleInventoryReservationFailed(
		ctx context.Context,
		payload inventoryevents.InventoryReservationFailedPayload,
	) error

	HandlePaymentCompleted(
		ctx context.Context,
		payload paymentevents.PaymentSucceededPayload,
	) error

	HandlePaymentFailed(
		ctx context.Context,
		payload paymentevents.PaymentFailedPayload,
	) error

	HandleInventoryReleased(
		context.Context,
		inventoryevents.InventoryReleasedPayload,
	) error
}
