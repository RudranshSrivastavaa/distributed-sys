package app

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"github.com/rudransh/distributed-commerce/order/internal/database"
	"github.com/rudransh/distributed-commerce/order/internal/handlers"
	"github.com/rudransh/distributed-commerce/order/internal/outbox"
	"github.com/rudransh/distributed-commerce/order/internal/repository"
	"github.com/rudransh/distributed-commerce/order/internal/routes"
	"github.com/rudransh/distributed-commerce/order/internal/service"

	commandevent "github.com/rudransh/distributed-commerce/order/internal/command"
	"github.com/rudransh/distributed-commerce/pkg/config"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
	"github.com/rudransh/distributed-commerce/pkg/logger"
)

func Start() {

	log := logger.New(config.Services.Order.Name)

	//----------------------------------------
	// Database
	//----------------------------------------

	database.Connect()

	db := database.DB

	//----------------------------------------
	// Kafka Producer (Used ONLY by Outbox Relay)
	//----------------------------------------

	kafkaConfig := kafkaa.DefaultConfig()

	kafkaConfig.Consumer.GroupID = "order-group"

	client := kafkaa.NewClient(kafkaConfig)

	producer, err := kafkaa.NewProducer(client)
	if err != nil {
		log.Fatal(err)
	}

	defer producer.Close()

	//----------------------------------------
	// Repositories
	//----------------------------------------

	orderRepo := repository.NewOrderRepository(db)

	outboxRepo := repository.NewOutboxRepository(db)

	//----------------------------------------
	// Outbox Relay
	//----------------------------------------

	relay := outbox.NewRelay(
		outbox.DefaultConfig(),
		outboxRepo,
		producer,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		log.Println("Starting Outbox Relay...")
		if err := relay.Start(ctx); err != nil {
			log.Printf("outbox relay stopped: %v", err)
		}

	}()

	defer relay.Stop()

	//----------------------------------------
	// Services
	//----------------------------------------

	orderService := service.NewOrderService(
		orderRepo,
	)

	dispatcher := kafkaa.NewDispatcher()

	dispatcher.Register(
		kafkaa.CompleteOrder,
		kafkaa.WrapHandler(
			commandevent.NewCompleteOrderHandler(
				orderService,
			),
		),
	)


	dispatcher.Register(
		kafkaa.CancelOrder,
		kafkaa.WrapHandler(
			commandevent.NewCancelOrderHandler(
				orderService,
			),
		),
	)

	host, err := kafkaa.NewConsumerHost(client)
	if err != nil {
		log.Fatal(err)
	}

	defer host.Close()

	host.Register(
		kafkaa.SagaCommands,
		dispatcher,
	)

	handler := handlers.NewOrderHandler(orderService)

	//----------------------------------------
	// Fiber
	//----------------------------------------

	app := fiber.New()

	routes.Register(app, handler)

	go func() {

		log.Println("Starting Order Consumer...")

		if err := host.Start(ctx); err != nil {
			log.Fatal(err)
		}

	}()

	log.Println("Order Service Started")

	log.Fatal(app.Listen(config.Services.Order.Port))
}
