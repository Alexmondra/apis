package main

import (
	"api-reniec/database"
	"api-reniec/jobs"
	"api-reniec/routes"
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	database.Connect()

	jobs.InitCron()

	app := fiber.New()
	routes.Setup(app)

	log.Println("🚀 Servidor Fiber corriendo en el puerto 9090")
	log.Fatal(app.Listen(":9090"))
}