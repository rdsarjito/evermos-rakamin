package main

import (
	"log"
	"os"

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

	// Check for force flag
	force := len(os.Args) > 1 && os.Args[1] == "--force"

	// Seed database
	if force {
		log.Println("force seeding database...")
		// For force, we can clear existing data first if needed
		// For now, just proceed with seed which checks if data exists
	}

	if err := repositories.SeedDatabase(); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	log.Println("seeding completed successfully")
}

