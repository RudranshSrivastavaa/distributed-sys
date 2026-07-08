package app

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/rudransh/distributed-commerce/pkg/circuitbreaker"
	"github.com/rudransh/distributed-commerce/payment/internal/database"
	"github.com/rudransh/distributed-commerce/payment/internal/handlers"
	"github.com/rudransh/distributed-commerce/payment/internal/provider"
	"github.com/rudransh/distributed-commerce/payment/internal/repository"
	"github.com/rudransh/distributed-commerce/pkg/retry"
	"github.com/rudransh/distributed-commerce/payment/internal/routes"
	"github.com/rudransh/distributed-commerce/payment/internal/service"
	"github.com/rudransh/distributed-commerce/payment/internal/webhook"

	"github.com/rudransh/distributed-commerce/pkg/config"
	"github.com/rudransh/distributed-commerce/pkg/logger"
)

func Start() {

	log := logger.New(config.Services.Payment.Name)

	database.Connect()

	db := database.DB

	app := fiber.New()

	//----------------------------------------------------
	// Dependencies
	//----------------------------------------------------

	configg := provider.SimulatorConfig{
		Delay:         10 * time.Second,
		SuccessRate:   100,
		WebhookSecret: "super-secret-key",
	}

	gateway := provider.NewSimulatorGateway(
		configg,
		"http://localhost:8083/payments/webhook",
	)

	paymentRepo := repository.NewPaymentRepository(db)

	attemptRepo := repository.NewPaymentAttemptRepository(db)

	txManager := repository.NewTransactionManager(db)

	//----------------------------------------------------
	// Service
	//----------------------------------------------------

	retryExecutor := retry.NewExecutor(
		retry.DefaultConfig(),
	)
	breaker := circuitbreaker.New(
		circuitbreaker.DefaultConfig(),
	)

	paymentService := service.NewPaymentService(
		paymentRepo,
		attemptRepo,
		txManager,
		gateway,
		retryExecutor,
		breaker,
	)

	//----------------------------------------------------
	// Handler
	//----------------------------------------------------

	verifier := webhook.NewVerifier(
		configg.WebhookSecret,
	)

	paymentHandler := handlers.NewPaymentHandler(

		paymentService,

		verifier,
	)

	//----------------------------------------------------
	// Routes
	//----------------------------------------------------

	routes.Register(
		app,
		paymentHandler,
	)

	log.Println("Payment Service Started")

	log.Fatal(
		app.Listen(config.Services.Payment.Port),
	)
}
