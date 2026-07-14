package app

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	paymentevent "github.com/rudransh/distributed-commerce/payment/internal/command"

	"github.com/rudransh/distributed-commerce/payment/internal/database"
	"github.com/rudransh/distributed-commerce/payment/internal/handlers"
	"github.com/rudransh/distributed-commerce/payment/internal/outbox"
	"github.com/rudransh/distributed-commerce/payment/internal/provider"
	"github.com/rudransh/distributed-commerce/payment/internal/repository"
	"github.com/rudransh/distributed-commerce/payment/internal/routes"
	"github.com/rudransh/distributed-commerce/payment/internal/service"
	"github.com/rudransh/distributed-commerce/payment/internal/webhook"

	"github.com/rudransh/distributed-commerce/pkg/circuitbreaker"
	"github.com/rudransh/distributed-commerce/pkg/config"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
	"github.com/rudransh/distributed-commerce/pkg/logger"
	"github.com/rudransh/distributed-commerce/pkg/retry"
)

func Start() {

	log := logger.New(config.Services.Payment.Name)

	//----------------------------------------------------
	// Database
	//----------------------------------------------------

	database.Connect()

	db := database.DB

	//----------------------------------------------------
	// Kafka
	//----------------------------------------------------

	kafkaConfig := kafkaa.DefaultConfig()
	kafkaConfig.Consumer.GroupID = "payment-group"


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

		log.Println("Starting Payment Outbox Relay...")

		if err := relay.Start(ctx); err != nil {
			log.Printf("Payment relay stopped: %v", err)
		}
	}()

	//----------------------------------------------------
	// Fiber
	//----------------------------------------------------

	app := fiber.New()

	//----------------------------------------------------
	// Payment Provider
	//----------------------------------------------------

	configg := provider.SimulatorConfig{
		Delay:         10 * time.Second,
		SuccessRate:  100,
		WebhookSecret: "super-secret-key",
	}

	gateway := provider.NewSimulatorGateway(
		configg,
		"http://localhost:8083/payments/webhook",
	)

	//----------------------------------------------------
	// Repositories
	//----------------------------------------------------

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
	// Kafka Consumer
	//----------------------------------------------------

	dispatcher := kafkaa.NewDispatcher()

	dispatcher.Register(
		kafkaa.ProcessPayment,
		kafkaa.WrapHandler(
			paymentevent.NewProcessPaymentHandler(
				paymentService,
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

	go func() {

		log.Println("Starting Payment Consumer...")

		if err := host.Start(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	//----------------------------------------------------
	// HTTP Handlers
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