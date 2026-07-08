package app

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/rudransh/distributed-commerce/inventory/internal/database"
	"github.com/rudransh/distributed-commerce/inventory/internal/handlers"
	"github.com/rudransh/distributed-commerce/inventory/internal/repository"
	"github.com/rudransh/distributed-commerce/inventory/internal/routes"
	"github.com/rudransh/distributed-commerce/inventory/internal/service"
	"github.com/rudransh/distributed-commerce/inventory/internal/worker"

	"github.com/rudransh/distributed-commerce/pkg/config"
	"github.com/rudransh/distributed-commerce/pkg/http/middleware"
	"github.com/rudransh/distributed-commerce/pkg/logger"
)

func Start() {

	log := logger.New(config.Services.Inventory.Name)

	database.Connect()

	app := fiber.New()

	app.Use(middleware.Recover())
	app.Use(middleware.RequestID())
	app.Use(middleware.CORS())
	app.Use(middleware.Logger())

	productRepo := repository.NewProductRepository(database.DB)

	inventoryRepo := repository.NewInventoryRepository(database.DB)

	reservationRepo := repository.NewReservationRepository(database.DB)

	txManager := repository.NewTransactionManager(database.DB)

	inventoryService := service.NewInventoryService(
		productRepo,
		inventoryRepo,
		reservationRepo,
		txManager,
	)
	inventoryWorker := worker.NewReservationExpiryWorker(
		inventoryService,
		time.Minute,
	)
	ctx := context.Background()
	go inventoryWorker.Start(ctx)

	handler := handlers.NewInventoryHandler(inventoryService)

	routes.Register(
		app,
		handler,
	)

	log.Println("Inventory Service Started")

	log.Fatal(app.Listen(config.Services.Inventory.Port))
}
