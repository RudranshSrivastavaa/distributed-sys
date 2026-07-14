package kafkaa

type AggregateType string

const (

	OrderAggregate AggregateType = "ORDER"

	InventoryAggregate AggregateType = "INVENTORY"

	PaymentAggregate AggregateType = "PAYMENT"

	NotificationAggregate AggregateType = "NOTIFICATION"

	SagaAggregate        AggregateType = "SAGA"
)