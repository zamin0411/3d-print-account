package router

import (
	"3d-print-account/handler"
	"3d-print-account/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	// Middleware
	api := app.Group("/api", logger.New())
	api.Get("/", handler.Ping)

	// Auth
	auth := api.Group("/auth")
	auth.Post("/login", middleware.LoginWithToken(), handler.LoginWithToken)

	auth.Post("/register", handler.Register)

	// User
	user := api.Group("/user")
	user.Put("/:id", handler.UpdateUser)
}
