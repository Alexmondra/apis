package middleware

import (
	"api-reniec/database"
	"api-reniec/models"
	"os"
	"github.com/gofiber/fiber/v2"
)

func ClientAuth(c *fiber.Ctx) error {
	apiKey := c.Get("X-API-KEY")
	masterToken := os.Getenv("MASTER_TOKEN")

	if apiKey != "" && apiKey == masterToken {
		return c.Next()
	}

	// 2. SI NO HAY TOKEN: Error
	if apiKey == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Falta API Key"})
	}

	// 3. SI ES TOKEN DE CLIENTE: Validar en DB y contar uso
	var client models.Client
	if err := database.DB.Where("api_key = ? AND status = ?", apiKey, true).First(&client).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Token inválido o cliente inactivo"})
	}

	if client.Usage >= client.Limit {
		return c.Status(429).JSON(fiber.Map{"error": "Límite de consultas excedido"})
	}

	client.Usage++
	database.DB.Save(&client)

	return c.Next()
}