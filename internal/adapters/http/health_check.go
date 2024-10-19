package http

import "github.com/gofiber/fiber/v2"

func NewHealthCheckController(app *fiber.App) {
	group := app.Group("/health")
	group.Get("/alive", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"ok": true,
		})
	},
	)
}
