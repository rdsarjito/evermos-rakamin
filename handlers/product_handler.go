package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"
	"github.com/rdsarjito/evermos-rakamin/domain/dto"
	"github.com/rdsarjito/evermos-rakamin/repositories"
	"github.com/rdsarjito/evermos-rakamin/services"
	"github.com/rdsarjito/evermos-rakamin/utils"
)

var validateProduct = validator.New()

// ProductHandler handles product HTTP requests
type ProductHandler struct {
	productService *services.ProductService
}

// NewProductHandler creates a new product handler
func NewProductHandler() *ProductHandler {
	return &ProductHandler{
		productService: services.NewProductService(),
	}
}

// ListProducts handles GET /products with pagination and search
// @Summary List all products
// @Description Get list of products with pagination, search, and category filter
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search query"
// @Param category_id query int false "Filter by category ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /products [get]
func (h *ProductHandler) ListProducts(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search", "")
	categoryID := c.Query("category_id", "")

	pagination, products, err := h.productService.GetAllProducts(page, limit, search, categoryID)
	if err != nil {
		return utils.ErrorResponse(c, err, fiber.StatusInternalServerError)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"products":   products,
		"pagination": pagination,
	}, "products retrieved successfully", fiber.StatusOK)
}

// GetProduct handles GET /products/:id
// @Summary Get product by ID
// @Description Get product details by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "invalid product ID"), fiber.StatusBadRequest)
	}

	product, err := h.productService.GetProductByID(uint(id))
	if err != nil {
		if err == repositories.ErrProductNotFound {
			return utils.ErrorResponse(c, err, fiber.StatusNotFound)
		}
		return utils.ErrorResponse(c, err, fiber.StatusInternalServerError)
	}

	return utils.SuccessResponse(c, product, "product retrieved successfully", fiber.StatusOK)
}

// CreateProduct handles POST /products
// @Summary Create a new product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.CreateProductRequest true "Create Product Request"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req dto.CreateProductRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "invalid request body"), fiber.StatusBadRequest)
	}

	// Validate request
	if err := validateProduct.Struct(req); err != nil {
		return utils.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()), fiber.StatusBadRequest)
	}

	product, err := h.productService.CreateProduct(req)
	if err != nil {
		if err == repositories.ErrCategoryNotFound {
			return utils.ErrorResponse(c, err, fiber.StatusNotFound)
		}
		return utils.ErrorResponse(c, err, fiber.StatusInternalServerError)
	}

	return utils.SuccessResponse(c, product, "product created successfully", fiber.StatusCreated)
}

// UpdateProduct handles PUT /products/:id
// @Summary Update product
// @Description Update product by ID
// @Tags products
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Product ID"
// @Param request body dto.UpdateProductRequest true "Update Product Request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "invalid product ID"), fiber.StatusBadRequest)
	}

	var req dto.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "invalid request body"), fiber.StatusBadRequest)
	}

	// Validate request
	if err := validateProduct.Struct(req); err != nil {
		return utils.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()), fiber.StatusBadRequest)
	}

	product, err := h.productService.UpdateProduct(uint(id), req)
	if err != nil {
		if err == repositories.ErrProductNotFound || err == repositories.ErrCategoryNotFound {
			return utils.ErrorResponse(c, err, fiber.StatusNotFound)
		}
		return utils.ErrorResponse(c, err, fiber.StatusInternalServerError)
	}

	return utils.SuccessResponse(c, product, "product updated successfully", fiber.StatusOK)
}

// DeleteProduct handles DELETE /products/:id
// @Summary Delete product
// @Description Delete product by ID (soft delete)
// @Tags products
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Product ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "invalid product ID"), fiber.StatusBadRequest)
	}

	if err := h.productService.DeleteProduct(uint(id)); err != nil {
		if err == repositories.ErrProductNotFound {
			return utils.ErrorResponse(c, err, fiber.StatusNotFound)
		}
		return utils.ErrorResponse(c, err, fiber.StatusInternalServerError)
	}

	return utils.SuccessResponse(c, nil, "product deleted successfully", fiber.StatusOK)
}

