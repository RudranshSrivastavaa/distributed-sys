package app

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/rudransh/distributed-commerce/notification/event"
	"github.com/rudransh/distributed-commerce/notification/internal/database"
	"github.com/rudransh/distributed-commerce/notification/internal/handlers"
	"github.com/rudransh/distributed-commerce/notification/internal/provider"
	"github.com/rudransh/distributed-commerce/notification/internal/repository"
	"github.com/rudransh/distributed-commerce/notification/internal/routes"
	"github.com/rudransh/distributed-commerce/notification/internal/service"

	"github.com/rudransh/distributed-commerce/pkg/circuitbreaker"
	"github.com/rudransh/distributed-commerce/pkg/config"
	"github.com/rudransh/distributed-commerce/pkg/http/middleware"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
	"github.com/rudransh/distributed-commerce/pkg/logger"
	"github.com/rudransh/distributed-commerce/pkg/retry"
)

func Start() {

	log := logger.New(config.Services.Notification.Name)

	//----------------------------------------------------
	// Database
	//----------------------------------------------------

	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}

	if err := database.Migrate(); err != nil {
		log.Fatal(err)
	}

	db := database.DB

	//----------------------------------------------------
	// Repositories
	//----------------------------------------------------

	notificationRepository := repository.NewNotificationRepository(db)

	transactionManager := repository.NewTransactionManager(db)

	emailProvider := provider.NewConsoleEmailProvider(
		90,
		500*time.Millisecond,
		2*time.Second,
	)

	retryExecutor := retry.NewExecutor(
		retry.DefaultConfig(),
	)
	breaker := circuitbreaker.New(
		circuitbreaker.DefaultConfig(),
	)

	notificationService := service.NewNotificationService(
		notificationRepository,
		transactionManager,
		emailProvider,
		retryExecutor,
		breaker)

	notificationHandler := handlers.NewNotificationHandler(
		notificationService,
	)
	// Temporary until service layer is added
	_ = notificationRepository
	_ = transactionManager
	_ = emailProvider
	_ = notificationService

	cfg := kafkaa.DefaultConfig()

	cfg.TLS.Enabled = true

	cfg.TLS.CAFile = "../../certs/ca/ca.crt"

	cfg.Consumer.GroupID = "notification-group"

	client := kafkaa.NewClient(cfg)

	host, _ := kafkaa.NewConsumerHost(client)

	dispatcher := kafkaa.NewDispatcher()

	//orders

	dispatcher.Register(
		kafkaa.OrderCreated,
		kafkaa.WrapHandler(
			event.NewOrderCreatedHandler(notificationService),
		),
	)

	//payments

	dispatcher.Register(
		kafkaa.PaymentCreated,
		kafkaa.WrapHandler(
			event.NewPaymentCreatedHandler(notificationService),
		),
	)

	dispatcher.Register(
		kafkaa.PaymentSucceeded,
		kafkaa.WrapHandler(
			event.NewPaymentSucceededHandler(notificationService),
		),
	)

	dispatcher.Register(
		kafkaa.PaymentFailed,
		kafkaa.WrapHandler(
			event.NewPaymentFailedHandler(notificationService),
		),
	)

	// inventory reservation
	
	dispatcher.Register(
		kafkaa.InventoryReserved,
		kafkaa.WrapHandler(
			event.NewInventoryReservedHandler(notificationService),
		),
	)
	dispatcher.Register(
		kafkaa.InventoryReleased,
		kafkaa.WrapHandler(
			event.NewInventoryReleasedHandler(notificationService),
		),
	)

	host.Register(
		kafkaa.PaymentEvents,
		dispatcher,
	)

	host.Register(
		kafkaa.OrderEvents,
		dispatcher,
	)

	host.Register(
		kafkaa.InventoryEvents,
		dispatcher,
	)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {

		if err := host.Start(ctx); err != nil {

			log.Printf(
				"kafka consumer stopped",
				"error",
				err,
			)
		}
	}()

	defer func() {

		cancel()

		host.Close()
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
	// Routes
	//----------------------------------------------------

	routes.Register(app, notificationHandler)

	//----------------------------------------------------
	// Start Server
	//----------------------------------------------------

	log.Println("Notification Service Started")

	log.Fatal(
		app.Listen(config.Services.Notification.Port),
	)
}
