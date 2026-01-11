package controllers

import (
	"api-reniec/database"
	"api-reniec/models"
	"api-reniec/services"
	"github.com/gofiber/fiber/v2"
)

// FUNCIÓN MAESTRA PARA EMPRESAS
func obtenerEmpresaGarantizada(ruc string) (*models.Empresa, error) {
	var empresa models.Empresa

	// 1. Buscar en Local
	err := database.DB.Preload("Contactos").First(&empresa, "ruc = ?", ruc).Error
	if err == nil {
		return &empresa, nil
	}

	// 2. Si no está, buscar en SUNAT (Decolecta)
	datosApi, err := services.FetchFromSunat(ruc)
	if err != nil {
		return nil, err
	}

	// 3. Guardar en MariaDB con todos los datos nuevos
	nuevaEmpresa := models.Empresa{
		RUC:               datosApi.NumeroDocumento,
		RazonSocial:      datosApi.RazonSocial,
		Estado:           datosApi.Estado,
		Condicion:        datosApi.Condicion,
		Direccion:        datosApi.Direccion,
		Ubigeo:           datosApi.Ubigeo,
		Distrito:         datosApi.Distrito,
		Provincia:        datosApi.Provincia,
		Departamento:     datosApi.Departamento,
		EsAgenteRetencion: datosApi.EsAgenteRetencion,
	}

	if err := database.DB.Create(&nuevaEmpresa).Error; err != nil {
		return nil, err
	}

	return &nuevaEmpresa, nil
}

// Handler: JSON Completo de Empresa
func GetEmpresa(c *fiber.Ctx) error {
	ruc := c.Params("ruc")
	empresa, err := obtenerEmpresaGarantizada(ruc)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(empresa)
}

func GetEmpresaReciente(c *fiber.Ctx) error {
	ruc := c.Params("ruc")
	empresa, err := obtenerEmpresaGarantizada(ruc)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	var correo, telefono string
	for _, con := range empresa.Contactos {
		if con.Tipo == "correo" && correo == "" { correo = con.Valor }
		if con.Tipo == "celular" && telefono == "" { telefono = con.Valor }
	}

	return c.JSON(fiber.Map{
		"ruc":          empresa.RUC,
		"razon_social": empresa.RazonSocial,
		"estado":       empresa.Estado,
		"direccion":    empresa.Direccion,
		"correo":       correo,
		"telefono":     telefono,
	})
}