package app

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/rudransh/distributed-commerce/notification/internal/database"
	"github.com/rudransh/distributed-commerce/notification/internal/handlers"
	"github.com/rudransh/distributed-commerce/notification/internal/provider"
	"github.com/rudransh/distributed-commerce/notification/internal/repository"
	"github.com/rudransh/distributed-commerce/notification/internal/routes"
	"github.com/rudransh/distributed-commerce/notification/internal/service"

	"github.com/rudransh/distributed-commerce/pkg/circuitbreaker"
	"github.com/rudransh/distributed-commerce/pkg/config"
	"github.com/rudransh/distributed-commerce/pkg/http/middleware"
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

	routes.Register(app,notificationHandler)

	//----------------------------------------------------
	// Start Server
	//----------------------------------------------------

	log.Println("Notification Service Started")

	log.Fatal(
		app.Listen(config.Services.Notification.Port),
	)
}
