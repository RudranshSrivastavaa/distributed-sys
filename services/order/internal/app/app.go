package app

import (
	"github.com/gofiber/fiber/v2"

	"github.com/rudransh/distributed-commerce/order/internal/database"
	"github.com/rudransh/distributed-commerce/order/internal/handlers"
	"github.com/rudransh/distributed-commerce/order/internal/repository"
	"github.com/rudransh/distributed-commerce/order/internal/routes"
	"github.com/rudransh/distributed-commerce/order/internal/service"

	"github.com/rudransh/distributed-commerce/pkg/config"
	"github.com/rudransh/distributed-commerce/pkg/kafka"
	"github.com/rudransh/distributed-commerce/pkg/logger"
)

func Start() {
	log := logger.New(config.Services.Order.Name)

	app := fiber.New()

	database.Connect()
	
	db := database.DB

	producer := kafka.NewProducer(
	kafka.DefaultConfig(),
)

	repo := repository.NewOrderRepository(db)

	service := service.NewOrderService(repo,producer)

	handler := handlers.NewOrderHandler(service)

	routes.Register(app, handler)

	log.Println("Order Service Started")

	log.Fatal(app.Listen(config.Services.Order.Port))
}
