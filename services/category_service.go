package services

import (
	"github.com/rdsarjito/evermos-rakamin/domain"
	"github.com/rdsarjito/evermos-rakamin/domain/dto"
	"github.com/rdsarjito/evermos-rakamin/repositories"
	"github.com/rdsarjito/evermos-rakamin/utils"
)

// CategoryService handles category business logic
type CategoryService struct {
	categoryRepo *repositories.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService() *CategoryService {
	return &CategoryService{
		categoryRepo: repositories.NewCategoryRepository(),
	}
}

// GetAllCategories gets all categories with pagination
func (s *CategoryService) GetAllCategories(page, limit int) (*utils.PaginationResponse, []dto.CategoryResponse, error) {
	offset, limitVal := utils.GetPaginationParams(page, limit)

	categories, total, err := s.categoryRepo.FindAll(offset, limitVal)
	if err != nil {
		return nil, nil, err
	}

	// Convert to DTOs
	responses := make([]dto.CategoryResponse, len(categories))
	for i, cat := range categories {
		responses[i] = dto.CategoryResponse{
			ID:          cat.ID,
			Name:        cat.Name,
			Description: cat.Description,
			CreatedAt:   cat.CreatedAt,
			UpdatedAt:   cat.UpdatedAt,
		}
	}

	pagination := &utils.PaginationResponse{
		Page:       page,
		Limit:      limitVal,
		Total:      total,
		TotalPages: utils.GetTotalPages(total, limitVal),
	}

	return pagination, responses, nil
}

// GetCategoryByID gets category by ID
func (s *CategoryService) GetCategoryByID(id uint) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(req dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	category := &domain.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		if err == repositories.ErrCategoryExists {
			return nil, err
		}
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

// UpdateCategory updates a category
func (s *CategoryService) UpdateCategory(id uint, req dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	// Get existing category
	existing, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Description != "" {
		existing.Description = req.Description
	}

	if err := s.categoryRepo.Update(existing); err != nil {
		if err == repositories.ErrCategoryExists {
			return nil, err
		}
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:          existing.ID,
		Name:        existing.Name,
		Description: existing.Description,
		CreatedAt:   existing.CreatedAt,
		UpdatedAt:   existing.UpdatedAt,
	}, nil
}

// DeleteCategory deletes a category
func (s *CategoryService) DeleteCategory(id uint) error {
	return s.categoryRepo.Delete(id)
}

