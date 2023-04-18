// Package main
package main

import (
	"3d-print-account/database"
	"3d-print-account/router"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Start web app
func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())
	database.Connect()
	router.SetupRoutes(app)
	app.Static("/", "./files")
	log.Fatal(app.Listen(":3000"))
}
