package app

import (
	"github.com/gofiber/fiber/v2"

	"github.com/rudransh/distributed-commerce/gateway/internal/routes"

	"github.com/rudransh/distributed-commerce/pkg/config"
	"github.com/rudransh/distributed-commerce/pkg/logger"
	"github.com/rudransh/distributed-commerce/pkg/http/middleware"
)

func Start() {

	log := logger.New(config.Services.Gateway.Name)

	app := fiber.New()

	app.Use(middleware.Recover())
	app.Use(middleware.RequestID())
	app.Use(middleware.CORS())
	app.Use(middleware.Logger())

	routes.Register(app)

	log.Println("Gateway Started")

	log.Fatal(app.Listen(config.Services.Gateway.Port))

}
