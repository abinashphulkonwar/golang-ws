package routes

import (
	"encoding/json"

	"github.com/abinashphulkonwar/master/db"
	"github.com/gofiber/fiber/v2"
)

func getConnection(c *fiber.Ctx) error {
	query := c.Query("id")
	println(query)
	if query == "" {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"message": "id not found π!"})

	}
	node, err := db.GetNodes(query)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"message": "node not found π!"})
	}
	return c.JSON(fiber.Map{"message": "hiiiii π!", "node": node.IP})
}

func postConnection(c *fiber.Ctx) error {
	body := db.Connection{}

	bytes := c.Body()
	err := json.Unmarshal(bytes, &body)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"message": "error in parsing π!", "err": err.Error()})
	}

	if body.Id == "" || body.Node == "" {
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{"message": "error in validationπ!"})
	}
	err = db.SetConnection(&body)

	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"message": "error in parsing π!", "err": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "add π!", "body": body})
}

func Connection(router fiber.Router) {

	router.Get("/get", getConnection)
	router.Post("/post", postConnection)
	router.All("*", func(c *fiber.Ctx) error {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{"message": "route not found π!"})
	})
}
