package proxy

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func Forward(target string) fiber.Handler {

	return func(c *fiber.Ctx) error {

		return proxy.Do(c, target+c.OriginalURL())
	}
}