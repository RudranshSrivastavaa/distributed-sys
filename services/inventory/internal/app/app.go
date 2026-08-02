package app

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"


	command "github.com/rudransh/distributed-commerce/inventory/internal/command"
	"github.com/rudransh/distributed-commerce/inventory/internal/database"
	"github.com/rudransh/distributed-commerce/inventory/internal/handlers"
	"github.com/rudransh/distributed-commerce/inventory/internal/outbox"
	"github.com/rudransh/distributed-commerce/inventory/internal/repository"
	"github.com/rudransh/distributed-commerce/inventory/internal/routes"
	"github.com/rudransh/distributed-commerce/inventory/internal/service"
	"github.com/rudransh/distributed-commerce/inventory/internal/worker"

	"github.com/rudransh/distributed-commerce/pkg/config"
	"github.com/rudransh/distributed-commerce/pkg/http/middleware"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
	"github.com/rudransh/distributed-commerce/pkg/logger"
)

func Start() {

	log := logger.New(config.Services.Inventory.Name)

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

	kafkaConfig.Consumer.GroupID = "inventory-group"

	client := kafkaa.NewClient(kafkaConfig)

	producer, err := kafkaa.NewProducer(client)
	if err != nil {
		log.Fatal(err)
	}


	defer producer.Close()


	//----------------------------------------------------
	// Outbox Relay
	//----------------------------------------------------

	outboxRepo := repository.NewOutboxRepository(db)

	relay := outbox.NewRelay(
		outbox.DefaultConfig(),
		outboxRepo,
		producer,
	)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()
	defer relay.Stop()

	go func() {

		log.Println("Starting Inventory Outbox Relay...")

		if err := relay.Start(ctx); err != nil {
			log.Printf("Inventory relay stopped: %v", err)
		}
	}()

	//----------------------------------------------------
	// Fiber
	//----------------------------------------------------

	app := fiber.New()

	app.Use(middleware.Recover())
	app.Use(middleware.RequestID())
	app.Use(middleware.CORS())
	app.Use(middleware.Logger())

	//----------------------------------------------------
	// Repositories
	//----------------------------------------------------

	productRepo := repository.NewProductRepository(db)

	inventoryRepo := repository.NewInventoryRepository(db)

	reservationRepo := repository.NewReservationRepository(db)

	txManager := repository.NewTransactionManager(db)

	//----------------------------------------------------
	// Service
	//----------------------------------------------------

	inventoryService := service.NewInventoryService(
		productRepo,
		inventoryRepo,
		reservationRepo,
		txManager,
	)

	//----------------------------------------------------
	// Kafka Consumer
	//----------------------------------------------------

	dispatcher := kafkaa.NewDispatcher()

	dispatcher.Register(
		kafkaa.ReserveInventory,
		kafkaa.WrapHandler(
			command.NewReserveInventoryHandler(
				inventoryService,
			),
		),
	)

	dispatcher.Register(
	kafkaa.ReleaseInventory,
	kafkaa.WrapHandler(
		command.NewReleaseInventoryHandler(
			inventoryService,
		),
	),
)

	host,_ := kafkaa.NewConsumerHost(client)

	host.Register(
		kafkaa.SagaCommands,
		dispatcher,
	)

	go func() {

		log.Println("Starting Inventory Consumer...")

		if err := host.Start(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	//----------------------------------------------------
	// Reservation Expiry Worker
	//----------------------------------------------------

	inventoryWorker := worker.NewReservationExpiryWorker(
		inventoryService,
		time.Minute,
	)

	go inventoryWorker.Start(ctx)

	//----------------------------------------------------
	// HTTP
	//----------------------------------------------------

	handler := handlers.NewInventoryHandler(
		inventoryService,
	)

	routes.Register(
		app,
		handler,
	)

	log.Println("Inventory Service Started")

	log.Fatal(
		app.Listen(config.Services.Inventory.Port),
	)
}