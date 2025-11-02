package dto

import "time"

// CreateProductRequest represents create product request
type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=200"`
	Description string  `json:"description" validate:"max=1000"`
	Price       float64 `json:"price" validate:"required,min=0"`
	Stock       int     `json:"stock" validate:"required,min=0"`
	CategoryID  uint    `json:"category_id" validate:"required"`
}

// UpdateProductRequest represents update product request
type UpdateProductRequest struct {
	Name        string  `json:"name" validate:"omitempty,min=3,max=200"`
	Description string  `json:"description" validate:"max=1000"`
	Price       float64 `json:"price" validate:"omitempty,min=0"`
	Stock       int     `json:"stock" validate:"omitempty,min=0"`
	CategoryID  uint    `json:"category_id" validate:"omitempty"`
}

// ProductResponse represents product response
type ProductResponse struct {
	ID          uint                `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Price       float64             `json:"price"`
	Stock       int                 `json:"stock"`
	CategoryID  uint                `json:"category_id"`
	Category    *CategoryResponse    `json:"category,omitempty"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

