package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/rdsarjito/evermos-rakamin/config"
	"github.com/rdsarjito/evermos-rakamin/handlers"
	"github.com/rdsarjito/evermos-rakamin/middleware"
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

	// Initialize handlers
	authHandler := handlers.NewAuthHandler()
	categoryHandler := handlers.NewCategoryHandler()
	productHandler := handlers.NewProductHandler()
	wilayahHandler := handlers.NewWilayahHandler()

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

	// Auth routes
	auth := app.Group("/auth")
	{
		auth.Post("/register", authHandler.Register)
		auth.Post("/login", authHandler.Login)
		auth.Get("/profile", middleware.AuthMiddleware(), authHandler.Profile)
	}

	// Category routes
	categories := app.Group("/categories")
	{
		categories.Get("/", categoryHandler.ListCategories)                    // Public
		categories.Get("/:id", categoryHandler.GetCategory)                    // Public
		categories.Post("/", middleware.AuthMiddleware(), categoryHandler.CreateCategory)
		categories.Put("/:id", middleware.AuthMiddleware(), categoryHandler.UpdateCategory)
		categories.Delete("/:id", middleware.AuthMiddleware(), categoryHandler.DeleteCategory)
	}

	// Product routes
	products := app.Group("/products")
	{
		products.Get("/", productHandler.ListProducts)                        // Public (with search & filter)
		products.Get("/:id", productHandler.GetProduct)                        // Public
		products.Post("/", middleware.AuthMiddleware(), productHandler.CreateProduct)
		products.Put("/:id", middleware.AuthMiddleware(), productHandler.UpdateProduct)
		products.Delete("/:id", middleware.AuthMiddleware(), productHandler.DeleteProduct)
	}

	// Location/Wilayah routes (API Wilayah Indonesia integration)
	locations := app.Group("/locations")
	{
		locations.Get("/provinces", wilayahHandler.GetProvinces)               // Public
		locations.Get("/regencies", wilayahHandler.GetRegencies)                // Public (requires province_id query param)
		locations.Get("/districts", wilayahHandler.GetDistricts)               // Public (requires regency_id query param)
	}

	addr := fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port)
	log.Printf("server starting on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
