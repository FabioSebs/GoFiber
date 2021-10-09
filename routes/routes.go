package routes

import fiber "github.com/gofiber/fiber/v2"

func Setup(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!")
	})
}
