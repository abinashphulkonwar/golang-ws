package routes

import "github.com/gofiber/fiber/v2"

func getNodes(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "hiiiii 👋!"})
}

func postNodes(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "hiiiii 👋!"})
}

func Node(router fiber.Router) {

	router.Get("/get", getNodes)
	router.Post("/post", postNodes)
	router.All("*", func(c *fiber.Ctx) error {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{"message": "route not found 👋!"})
	})
}
