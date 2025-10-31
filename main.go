package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/rdsarjito/evermos-rakamin/config"
	"github.com/rdsarjito/evermos-rakamin/repositories"
)

func main() {
	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize database
	if err := repositories.InitDatabase(); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer func() {
		if err := repositories.CloseDatabase(); err != nil {
			log.Printf("error closing database: %v", err)
		}
	}()

	cfg := config.Get()

	app := fiber.New()

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "service healthy",
		})
	})

	// Database health check endpoint
	app.Get("/health/db", func(c *fiber.Ctx) error {
		if err := repositories.HealthCheck(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":  "error",
				"message": "database connection failed",
				"error":   err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "database connection healthy",
		})
	})

	addr := fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port)
	log.Printf("server starting on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
