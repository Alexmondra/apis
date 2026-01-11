package controllers

import (
	"api-reniec/database"
	"api-reniec/models"
	"api-reniec/services"
	"github.com/gofiber/fiber/v2"
)

// FUNCIÓN MAESTRA: Maneja la lógica de búsqueda y guardado automático
func obtenerPersonaGarantizada(dni string) (*models.Persona, error) {
	var persona models.Persona
	err := database.DB.Preload("Contactos").First(&persona, "dni = ?", dni).Error
	if err == nil {
		return &persona, nil // Lo encontró en local, lo devuelve
	}

	datosApi, err := services.FetchFromReniec(dni)
	if err != nil {
		return nil, err // No existe en ninguna parte
	}
	nuevaPersona := models.Persona{
		DNI:     datosApi.DocumentNumber,
		Nombres: datosApi.FirstName,
		Paterno: datosApi.FirstLastName,
		Materno: datosApi.SecondLastName,
	}
	
	if err := database.DB.Create(&nuevaPersona).Error; err != nil {
		return nil, err
	}

	return &nuevaPersona, nil
}

func GetPersona(c *fiber.Ctx) error {
	dni := c.Params("dni")
	persona, err := obtenerPersonaGarantizada(dni)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(persona)
}

func GetPersonaReciente(c *fiber.Ctx) error {
	dni := c.Params("dni")
	persona, err := obtenerPersonaGarantizada(dni)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	var correo, telefono string
	for _, con := range persona.Contactos {
		if con.Tipo == "correo" && correo == "" {
			correo = con.Valor
		}
		if con.Tipo == "celular" && telefono == "" {
			telefono = con.Valor
		}
	}

	return c.JSON(fiber.Map{
		"dni":              persona.DNI,
		"nombres":          persona.Nombres,
		"apellido_paterno": persona.Paterno,
		"apellido_materno": persona.Materno,
		"correo":           correo,
		"telefono":         telefono,
	})
}