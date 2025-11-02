package services

import (
	"strconv"

	"github.com/rdsarjito/evermos-rakamin/domain"
	"github.com/rdsarjito/evermos-rakamin/domain/dto"
	"github.com/rdsarjito/evermos-rakamin/repositories"
	"github.com/rdsarjito/evermos-rakamin/utils"
)

// ProductService handles product business logic
type ProductService struct {
	productRepo  *repositories.ProductRepository
	categoryRepo *repositories.CategoryRepository
}

// NewProductService creates a new product service
func NewProductService() *ProductService {
	return &ProductService{
		productRepo:  repositories.NewProductRepository(),
		categoryRepo: repositories.NewCategoryRepository(),
	}
}

// GetAllProducts gets all products with pagination and search
func (s *ProductService) GetAllProducts(page, limit int, search string, categoryIDStr string) (*utils.PaginationResponse, []dto.ProductResponse, error) {
	offset, limitVal := utils.GetPaginationParams(page, limit)

	var categoryID uint = 0
	if categoryIDStr != "" {
		if parsedID, err := strconv.ParseUint(categoryIDStr, 10, 32); err == nil {
			categoryID = uint(parsedID)
		}
	}

	products, total, err := s.productRepo.FindAll(offset, limitVal, search, categoryID)
	if err != nil {
		return nil, nil, err
	}

	// Convert to DTOs
	responses := make([]dto.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = dto.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			CategoryID:  product.CategoryID,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		}

		// Include category if loaded
		if product.Category.ID > 0 {
			responses[i].Category = &dto.CategoryResponse{
				ID:          product.Category.ID,
				Name:        product.Category.Name,
				Description: product.Category.Description,
				CreatedAt:   product.Category.CreatedAt,
				UpdatedAt:   product.Category.UpdatedAt,
			}
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

// GetProductByID gets product by ID
func (s *ProductService) GetProductByID(id uint) (*dto.ProductResponse, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := &dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		CategoryID:  product.CategoryID,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}

	// Include category if loaded
	if product.Category.ID > 0 {
		response.Category = &dto.CategoryResponse{
			ID:          product.Category.ID,
			Name:        product.Category.Name,
			Description: product.Category.Description,
			CreatedAt:   product.Category.CreatedAt,
			UpdatedAt:   product.Category.UpdatedAt,
		}
	}

	return response, nil
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	// Verify category exists
	_, err := s.categoryRepo.FindByID(req.CategoryID)
	if err != nil {
		if err == repositories.ErrCategoryNotFound {
			return nil, repositories.ErrCategoryNotFound
		}
		return nil, err
	}

	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, err
	}

	// Reload with category
	created, err := s.productRepo.FindByID(product.ID)
	if err != nil {
		return nil, err
	}

	response := &dto.ProductResponse{
		ID:          created.ID,
		Name:        created.Name,
		Description: created.Description,
		Price:       created.Price,
		Stock:       created.Stock,
		CategoryID:  created.CategoryID,
		CreatedAt:   created.CreatedAt,
		UpdatedAt:   created.UpdatedAt,
	}

	if created.Category.ID > 0 {
		response.Category = &dto.CategoryResponse{
			ID:          created.Category.ID,
			Name:        created.Category.Name,
			Description: created.Category.Description,
			CreatedAt:   created.Category.CreatedAt,
			UpdatedAt:   created.Category.UpdatedAt,
		}
	}

	return response, nil
}

// UpdateProduct updates a product
func (s *ProductService) UpdateProduct(id uint, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	// Get existing product
	existing, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Verify category exists if being updated
	if req.CategoryID > 0 {
		_, err := s.categoryRepo.FindByID(req.CategoryID)
		if err != nil {
			if err == repositories.ErrCategoryNotFound {
				return nil, repositories.ErrCategoryNotFound
			}
			return nil, err
		}
	}

	// Update fields if provided
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.Price > 0 {
		existing.Price = req.Price
	}
	if req.Stock >= 0 {
		existing.Stock = req.Stock
	}
	if req.CategoryID > 0 {
		existing.CategoryID = req.CategoryID
	}

	if err := s.productRepo.Update(existing); err != nil {
		return nil, err
	}

	// Reload with category
	updated, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := &dto.ProductResponse{
		ID:          updated.ID,
		Name:        updated.Name,
		Description: updated.Description,
		Price:       updated.Price,
		Stock:       updated.Stock,
		CategoryID:  updated.CategoryID,
		CreatedAt:   updated.CreatedAt,
		UpdatedAt:   updated.UpdatedAt,
	}

	if updated.Category.ID > 0 {
		response.Category = &dto.CategoryResponse{
			ID:          updated.Category.ID,
			Name:        updated.Category.Name,
			Description: updated.Category.Description,
			CreatedAt:   updated.Category.CreatedAt,
			UpdatedAt:   updated.Category.UpdatedAt,
		}
	}

	return response, nil
}

// DeleteProduct deletes a product
func (s *ProductService) DeleteProduct(id uint) error {
	return s.productRepo.Delete(id)
}

