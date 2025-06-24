package api

import (
	"prtask2/internal/config"
	"prtask2/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(taskService service.Service) *fiber.App {
	cfg := config.Load()

	app := fiber.New(fiber.Config{
		ServerHeader: cfg.ServerName,
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"app":    cfg.AppName,
		})
	})

	api := app.Group("/api/v1")

	api.Post("/tasks", taskService.CreateTask)                           // POST /api/v1/tasks
	api.Get("/tasks", taskService.GetAllTasks)                           // GET /api/v1/tasks
	api.Get("/tasks/:id", taskService.GetTaskByID)                       // GET /api/v1/tasks/:id
	api.Put("/tasks/:id", taskService.UpdateTask)                        // PUT /api/v1/tasks/:id
	api.Delete("/tasks/:id", taskService.DeleteTask)                     // DELETE /api/v1/tasks/:id
	api.Get("/tasks/user/:user_id", taskService.GetTasksByUserID)        // GET /api/v1/tasks/user/:user_id
	api.Get("/tasks/username/:username", taskService.GetTasksByUsername) // GET /api/v1/tasks/username/:username

	return app
}
