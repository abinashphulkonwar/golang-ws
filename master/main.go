package main

import (
	"github.com/abinashphulkonwar/master/db"
	"github.com/abinashphulkonwar/master/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	defer db.DB.Close()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hiiii 👋!")
	})

	app.Route("/nodes", routes.Node)
	app.Route("/connections", routes.Connection)

	app.Listen(":3000")
}
