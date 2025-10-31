package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/rdsarjito/evermos-rakamin/config"
)

func main() {
	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	cfg := config.Get()

	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "service healthy",
		})
	})

	addr := fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port)
	log.Printf("server starting on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
