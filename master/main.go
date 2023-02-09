package main

import (
	"github.com/abinashphulkonwar/master/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Route("/nodes", routes.Node)

	app.Listen(":3000")
}
