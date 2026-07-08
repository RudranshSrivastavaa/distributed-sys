package kafka

type Topic string

const (

	OrderEvents Topic = "order.events"

	PaymentEvents Topic = "payment.events"

	InventoryEvents Topic = "inventory.events"

	NotificationEvents Topic = "notification.events"
)