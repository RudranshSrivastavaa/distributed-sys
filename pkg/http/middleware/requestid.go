package middleware

import (
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2"
)
func RequestID() fiber.Handler {
	return requestid.New()
}