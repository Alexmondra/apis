package middleware

import (
	"os"
	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
	masterToken := os.Getenv("MASTER_TOKEN")
	authHeader := c.Get("Authorization")

	if authHeader != "Bearer "+masterToken {
		return c.Status(401).JSON(fiber.Map{"error": "Acceso administrativo denegado"})
	}
	return c.Next()
}