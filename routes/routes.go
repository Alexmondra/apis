package routes

import (
	"api-reniec/controllers"
	"api-reniec/middleware"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	// --- 1. RUTAS DE ADMIN (Usa MASTER_TOKEN vía Header Authorization: Bearer) ---
	admin := api.Group("/admin", middleware.AdminAuth)
	admin.Post("/clients", controllers.CreateClient)
	admin.Get("/clients", controllers.GetClients)
	admin.Patch("/clients/:id/limit", controllers.UpdateLimit)
	admin.Patch("/clients/:id/status", controllers.UpdateClientStatus)

	// --- 2. RUTAS DE CONSULTA (X-API-KEY: Maestro o Cliente) ---
	consultas := api.Group("/consulta", middleware.ClientAuth)
	
	// Personas (DNI)
	consultas.Get("/persona/:dni", controllers.GetPersona)
	consultas.Get("/persona-resumen/:dni", controllers.GetPersonaReciente)

	// Empresas (RUC)
	consultas.Get("/empresa/:ruc", controllers.GetEmpresa)
	consultas.Get("/empresa-resumen/:ruc", controllers.GetEmpresaReciente)
	
	// Estado del Token (Para que el cliente vea cuánto le queda)
	consultas.Get("/status", controllers.GetMyStatus)
}