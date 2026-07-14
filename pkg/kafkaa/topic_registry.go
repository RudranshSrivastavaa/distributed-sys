package kafkaa

var (
	OrderEvents = Topic{

		Name: "order-events",

		ConsumerGroup: "inventory-group",

		DeadLetterTopic: "order-events-dlq",

		Aggregate: OrderAggregate,

		Description: "Order domain events",
	}

	InventoryEvents = Topic{

		Name: "inventory-events",

		ConsumerGroup: "payment-group",

		DeadLetterTopic: "inventory-events-dlq",

		Aggregate: InventoryAggregate,

		Description: "Inventory events",
	}

	PaymentEvents = Topic{

		Name: "payment-events",

		ConsumerGroup: "notification-group",

		DeadLetterTopic: "payment-events-dlq",

		Aggregate: PaymentAggregate,

		Description: "Payment events",
	}

	NotificationEvents = Topic{

		Name: "notification-events",

		ConsumerGroup: "analytics-group",

		DeadLetterTopic: "notification-events-dlq",

		Aggregate: NotificationAggregate,

		Description: "Notification events",
	}

	SagaCommands = Topic{

		Name: "saga-commands",

		ConsumerGroup: "saga-group",

		DeadLetterTopic: "saga-commands-dlq",

		Aggregate: SagaAggregate,

		Description: "Saga orchestration commands",
	}
)
