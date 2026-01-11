package controllers

import (
	"api-reniec/database"
	"api-reniec/models"
	"crypto/rand"
	"os"
	"encoding/hex"
	"github.com/gofiber/fiber/v2"
)

func generateApiKey() string {
	b := make([]byte, 20)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func CreateClient(c *fiber.Ctx) error {
	type Request struct {
		Empresa string `json:"empresa"`
		Ruc     string `json:"ruc"`
		Limit   int    `json:"limit"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	newClient := models.Client{
		Empresa: req.Empresa,
		Ruc:     req.Ruc,
		ApiKey:  generateApiKey(),
		Limit:   req.Limit,
	}

	if err := database.DB.Create(&newClient).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al crear cliente o RUC duplicado"})
	}

	return c.JSON(newClient)
}

func GetClients(c *fiber.Ctx) error {
	var clients []models.Client
	database.DB.Find(&clients)
	return c.JSON(clients)
}

// Aumentar límite a un cliente
func UpdateLimit(c *fiber.Ctx) error {
	id := c.Params("id")
	type Request struct {
		NewLimit int `json:"new_limit"`
	}
	var req Request
	c.BodyParser(&req)

	database.DB.Model(&models.Client{}).Where("id = ?", id).Update("limit", req.NewLimit)
	return c.JSON(fiber.Map{"message": "Límite actualizado correctamente"})
}

func GetMyStatus(c *fiber.Ctx) error {
	apiKey := c.Get("X-API-KEY")
	masterToken := os.Getenv("MASTER_TOKEN")

	if apiKey == masterToken {
		return c.JSON(fiber.Map{
			"empresa": "ADMINISTRADOR",
			"limit":   "ILIMITADO",
			"usage":   "N/A",
		})
	}

	var client models.Client
	database.DB.Where("api_key = ?", apiKey).First(&client)
	
	return c.JSON(fiber.Map{
		"empresa": client.Empresa,
		"ruc":     client.Ruc,
		"limit":   client.Limit,
		"usage":   client.Usage,
		"restantes": client.Limit - client.Usage,
	})
}

// Desactivar o Activar un cliente (Baja/Alta)
func UpdateClientStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	type Request struct {
		Status bool `json:"status"` // true para activar, false para dar de baja
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	if err := database.DB.Model(&models.Client{}).Where("id = ?", id).Update("status", req.Status).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "No se pudo actualizar el estado"})
	}

	mensaje := "Cliente desactivado correctamente"
	if req.Status {
		mensaje = "Cliente activado correctamente"
	}

	return c.JSON(fiber.Map{"message": mensaje, "id": id, "nuevo_status": req.Status})
}