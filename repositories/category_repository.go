package repositories

import (
	"errors"

	"github.com/rdsarjito/evermos-rakamin/domain"
	"gorm.io/gorm"
)

// CategoryRepository handles category database operations
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{
		db: GetDB(),
	}
}

// FindAll finds all categories with pagination
func (r *CategoryRepository) FindAll(offset, limit int) ([]domain.Category, int64, error) {
	var categories []domain.Category
	var total int64

	// Count total
	if err := r.db.Model(&domain.Category{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

// FindByID finds category by ID
func (r *CategoryRepository) FindByID(id uint) (*domain.Category, error) {
	var category domain.Category
	if err := r.db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

// Create creates a new category
func (r *CategoryRepository) Create(category *domain.Category) error {
	// Check if name already exists
	var existingCategory domain.Category
	if err := r.db.Where("name = ?", category.Name).First(&existingCategory).Error; err == nil {
		return ErrCategoryExists
	}

	if err := r.db.Create(category).Error; err != nil {
		return err
	}
	return nil
}

// Update updates a category
func (r *CategoryRepository) Update(category *domain.Category) error {
	// Check if category exists
	if _, err := r.FindByID(category.ID); err != nil {
		return err
	}

	// Check if new name conflicts with existing category
	if category.Name != "" {
		var existingCategory domain.Category
		if err := r.db.Where("name = ? AND id != ?", category.Name, category.ID).First(&existingCategory).Error; err == nil {
			return ErrCategoryExists
		}
	}

	if err := r.db.Model(category).Updates(category).Error; err != nil {
		return err
	}
	return nil
}

// Delete soft deletes a category
func (r *CategoryRepository) Delete(id uint) error {
	// Check if category exists
	if _, err := r.FindByID(id); err != nil {
		return err
	}

	if err := r.db.Delete(&domain.Category{}, id).Error; err != nil {
		return err
	}
	return nil
}

