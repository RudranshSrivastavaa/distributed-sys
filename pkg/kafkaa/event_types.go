package kafkaa

type EventType string

const (

	OrderCreated EventType = "ORDER_CREATED"

	OrderCancelled EventType = "ORDER_CANCELLED"

	InventoryReserved EventType = "INVENTORY_RESERVED"

	InventoryReleased EventType = "INVENTORY_RELEASED"

	PaymentCreated EventType = "PAYMENT_CREATED"

	PaymentSucceeded EventType = "PAYMENT_SUCCEEDED"

	PaymentFailed EventType = "PAYMENT_FAILED"

	NotificationSent EventType = "NOTIFICATION_SENT"

	NotificationFailed EventType = "NOTIFICATION_FAILED"

	//commands

	ReserveInventory EventType = "RESERVE_INVENTORY"

	ProcessPayment EventType = "PROCESS_PAYMENT"

	ReleaseInventory EventType = "RELEASE_INVENTORY"

	CompleteOrder EventType = "COMPLETE_ORDER"

	CancelOrder EventType = "CANCEL_ORDER"

	InventoryReservationFailed EventType = "INVENTORY_RESERVATION_FAILED"
)