package app

import (
	"context"

	"github.com/rudransh/distributed-commerce/pkg/config"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
	"github.com/rudransh/distributed-commerce/pkg/logger"

	//sagaevent "github.com/rudransh/distributed-commerce/saga/event"
	"github.com/rudransh/distributed-commerce/saga/internal/database"
	"github.com/rudransh/distributed-commerce/saga/internal/handlers"
	"github.com/rudransh/distributed-commerce/saga/internal/repository"
	"github.com/rudransh/distributed-commerce/saga/internal/service"
)

func Start() {

	log := logger.New(config.Services.Saga.Name)

	//----------------------------------------------------
	// Database
	//----------------------------------------------------

	database.Connect()

	db := database.DB

	//----------------------------------------------------
	// Kafka
	//----------------------------------------------------

	kafkaConfig := kafkaa.DefaultConfig()

	kafkaConfig.TLS.Enabled = true

	kafkaConfig.TLS.CAFile = "../../certs/ca/ca.crt"

	kafkaConfig.Consumer.GroupID= "saga-group";

	client := kafkaa.NewClient(kafkaConfig)

	producer, err := kafkaa.NewProducer(client)
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	//----------------------------------------------------
	// Repository
	//----------------------------------------------------

	sagaRepository := repository.NewSagaRepository(db)

	//----------------------------------------------------
	// Service
	//----------------------------------------------------

	sagaService := service.NewSagaService(
		sagaRepository,
		producer,
	)

	//----------------------------------------------------
	// Dispatcher
	//----------------------------------------------------

	dispatcher := kafkaa.NewDispatcher()

	// ORDER_CREATED
	dispatcher.Register(
		kafkaa.OrderCreated,
		kafkaa.WrapHandler(
			handlers.NewOrderCreatedHandler(
				sagaService,
			),
		),
	)

	// INVENTORY_RESERVED
	dispatcher.Register(
		kafkaa.InventoryReserved,
		kafkaa.WrapHandler(
			handlers.NewInventoryReservedHandler(
				sagaService,
			),
		),
	)

	dispatcher.Register(
	kafkaa.InventoryReservationFailed,
	kafkaa.WrapHandler(
		handlers.NewInventoryReservationFailedHandler(
			sagaService,
		),
	),
)

	// // INVENTORY_RELEASED (temporary placeholder)
	dispatcher.Register(
		kafkaa.InventoryReleased,
		kafkaa.WrapHandler(
			handlers.NewInventoryReleasedHandler(
				sagaService,
			),
		),
	)

	//PAYMENT_SUCCEEDED
	dispatcher.Register(
		kafkaa.PaymentSucceeded,
		kafkaa.WrapHandler(
			handlers.NewPaymentSucceededHandler(
				sagaService,
			),
		),
	)

	//PAYMENT_FAILED
	dispatcher.Register(
		kafkaa.PaymentFailed,
		kafkaa.WrapHandler(
			handlers.NewPaymentFailedHandler(
				sagaService,
			),
		),
	)

	//----------------------------------------------------
	// Consumer Host
	//----------------------------------------------------

	host, err := kafkaa.NewConsumerHost(client)
	if err != nil {
		log.Fatal(err)
	}

	host.Register(
		kafkaa.OrderEvents,
		dispatcher,
	)

	host.Register(
		kafkaa.InventoryEvents,
		dispatcher,
	)

	host.Register(
		kafkaa.PaymentEvents,
		dispatcher,
	)

	//----------------------------------------------------
	// Start Consumer
	//----------------------------------------------------

	ctx := context.Background()

	log.Println("Saga Service Started")

	if err := host.Start(ctx); err != nil {
		log.Fatal(err)
	}

	
}