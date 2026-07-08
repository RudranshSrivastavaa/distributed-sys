package middleware

import (
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2"
)

func Recover() fiber.Handler {
	return recover.New()
}