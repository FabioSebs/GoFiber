package main

import (
	"./database"
	"./routes"
	fiber "github.com/gofiber/fiber/v2"
)

func main() {
	database.Connect()
	app := fiber.New()
	routes.Setup(app)
	app.Listen(":3000")
}
