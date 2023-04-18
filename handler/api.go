package handler

import "github.com/gofiber/fiber/v2"

// Hello hanlde api status
func Ping(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "success", "message": "pong", "data": nil})
}
