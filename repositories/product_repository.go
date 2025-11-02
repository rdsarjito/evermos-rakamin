package repositories

import (
	"errors"
	"strings"

	"github.com/rdsarjito/evermos-rakamin/domain"
	"gorm.io/gorm"
)

// ProductRepository handles product database operations
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		db: GetDB(),
	}
}

// FindAll finds all products with pagination and search
func (r *ProductRepository) FindAll(offset, limit int, search string, categoryID uint) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	query := r.db.Model(&domain.Product{})

	// Apply search filter
	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchPattern, searchPattern)
	}

	// Apply category filter
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results with category preload
	if err := query.Preload("Category").Offset(offset).Limit(limit).Order("created_at DESC").Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// FindByID finds product by ID
func (r *ProductRepository) FindByID(id uint) (*domain.Product, error) {
	var product domain.Product
	if err := r.db.Preload("Category").First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

// Create creates a new product
func (r *ProductRepository) Create(product *domain.Product) error {
	// Check if category exists
	var category domain.Category
	if err := r.db.First(&category, product.CategoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCategoryNotFound
		}
		return err
	}

	if err := r.db.Create(product).Error; err != nil {
		return err
	}
	return nil
}

// Update updates a product
func (r *ProductRepository) Update(product *domain.Product) error {
	// Check if product exists
	if _, err := r.FindByID(product.ID); err != nil {
		return err
	}

	// Check if category exists (if category_id is being updated)
	if product.CategoryID > 0 {
		var category domain.Category
		if err := r.db.First(&category, product.CategoryID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrCategoryNotFound
			}
			return err
		}
	}

	// Update product
	updateData := make(map[string]interface{})
	if product.Name != "" {
		updateData["name"] = product.Name
	}
	if product.Description != "" {
		updateData["description"] = product.Description
	}
	if product.Price > 0 {
		updateData["price"] = product.Price
	}
	if product.Stock >= 0 {
		updateData["stock"] = product.Stock
	}
	if product.CategoryID > 0 {
		updateData["category_id"] = product.CategoryID
	}

	if err := r.db.Model(&domain.Product{}).Where("id = ?", product.ID).Updates(updateData).Error; err != nil {
		return err
	}
	return nil
}

// Delete soft deletes a product
func (r *ProductRepository) Delete(id uint) error {
	// Check if product exists
	if _, err := r.FindByID(id); err != nil {
		return err
	}

	if err := r.db.Delete(&domain.Product{}, id).Error; err != nil {
		return err
	}
	return nil
}

