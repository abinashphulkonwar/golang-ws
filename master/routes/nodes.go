package routes

import (
	"encoding/json"

	"github.com/abinashphulkonwar/master/db"
	"github.com/gofiber/fiber/v2"
)

func getNodes(c *fiber.Ctx) error {
	query := c.Query("id")

	if query == "" {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"message": "id not found 👋!"})

	}
	node, err := db.GetNodes(query)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"message": "node not found 👋!"})
	}
	return c.JSON(fiber.Map{"message": "hiiiii 👋!", "node": node.IP})
}

func postNodes(c *fiber.Ctx) error {
	body := db.Node{}
	bytes := c.Body()
	err := json.Unmarshal(bytes, &body)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"message": "error in parsing 👋!", "err": err.Error()})
	}
	if body.IP == "" || body.NAME == "" || body.STATUS == "" {
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{"message": "error in validation👋!"})
	}

	_, err = db.SetNode(&body)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"message": "error in parsing 👋!", "err": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "add 👋!", "body": body})
}

func Node(router fiber.Router) {

	router.Get("/get", getNodes)
	router.Post("/post", postNodes)
	router.All("*", func(c *fiber.Ctx) error {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{"message": "route not found 👋!"})
	})
}
