package main

import (
	"github.com/FabioSebs/GoFiber/backend/database"
	"github.com/FabioSebs/GoFiber/backend/routes"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	database.Connect()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
	routes.Setup(app)
	app.Listen(":3000")
}
