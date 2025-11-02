package repositories

import (
	"log"

	"github.com/rdsarjito/evermos-rakamin/constants"
	"github.com/rdsarjito/evermos-rakamin/domain"
	"golang.org/x/crypto/bcrypt"
)

// SeedDatabase seeds the database with sample data
func SeedDatabase() error {
	db := GetDB()

	// Check if data already exists
	var count int64
	db.Model(&domain.User{}).Count(&count)
	if count > 0 {
		log.Println("database already seeded, skipping...")
		return nil
	}

	// Seed Users
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	users := []domain.User{
		{
			Name:         "Admin User",
			Email:        "admin@example.com",
			PasswordHash: string(hashedPassword),
			Role:         constants.RoleAdmin,
		},
		{
			Name:         "Regular User",
			Email:        "user@example.com",
			PasswordHash: string(hashedPassword),
			Role:         constants.RoleUser,
		},
		{
			Name:         "Seller User",
			Email:        "seller@example.com",
			PasswordHash: string(hashedPassword),
			Role:         constants.RoleSeller,
		},
	}

	if err := db.Create(&users).Error; err != nil {
		return err
	}
	log.Println("seeded users")

	// Seed Categories
	categories := []domain.Category{
		{
			Name:        "Electronics",
			Description: "Electronic devices and accessories",
		},
		{
			Name:        "Clothing",
			Description: "Fashion and apparel",
		},
		{
			Name:        "Food & Beverages",
			Description: "Food and drink products",
		},
		{
			Name:        "Books",
			Description: "Books and reading materials",
		},
	}

	if err := db.Create(&categories).Error; err != nil {
		return err
	}
	log.Println("seeded categories")

	// Seed Products
	products := []domain.Product{
		{
			Name:        "Laptop",
			Description: "High-performance laptop for work and gaming",
			Price:       8999.99,
			Stock:       10,
			CategoryID:  1, // Electronics
		},
		{
			Name:        "Smartphone",
			Description: "Latest smartphone with advanced features",
			Price:       5999.99,
			Stock:       25,
			CategoryID:  1, // Electronics
		},
		{
			Name:        "T-Shirt",
			Description: "Comfortable cotton t-shirt",
			Price:       199.99,
			Stock:       50,
			CategoryID:  2, // Clothing
		},
		{
			Name:        "Jeans",
			Description: "Classic blue jeans",
			Price:       499.99,
			Stock:       30,
			CategoryID:  2, // Clothing
		},
		{
			Name:        "Coffee Beans",
			Description: "Premium arabica coffee beans",
			Price:       149.99,
			Stock:       100,
			CategoryID:  3, // Food & Beverages
		},
		{
			Name:        "Programming Book",
			Description: "Learn Go programming language",
			Price:       299.99,
			Stock:       20,
			CategoryID:  4, // Books
		},
	}

	if err := db.Create(&products).Error; err != nil {
		return err
	}
	log.Println("seeded products")

	log.Println("database seeded successfully")
	return nil
}

