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

var validateCategory = validator.New()

// CategoryHandler handles category HTTP requests
type CategoryHandler struct {
	categoryService *services.CategoryService
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{
		categoryService: services.NewCategoryService(),
	}
}

// ListCategories handles GET /categories with pagination
// @Summary List all categories
// @Description Get list of categories with pagination
// @Tags categories
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /categories [get]
func (h *CategoryHandler) ListCategories(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	pagination, categories, err := h.categoryService.GetAllCategories(page, limit)
	if err != nil {
		return utils.ErrorResponse(c, err, fiber.StatusInternalServerError)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"categories": categories,
		"pagination": pagination,
	}, "categories retrieved successfully", fiber.StatusOK)
}

// GetCategory handles GET /categories/:id
// @Summary Get category by ID
// @Description Get category details by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategory(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "invalid category ID"), fiber.StatusBadRequest)
	}

	category, err := h.categoryService.GetCategoryByID(uint(id))
	if err != nil {
		if err == repositories.ErrCategoryNotFound {
			return utils.ErrorResponse(c, err, fiber.StatusNotFound)
		}
		return utils.ErrorResponse(c, err, fiber.StatusInternalServerError)
	}

	return utils.SuccessResponse(c, category, "category retrieved successfully", fiber.StatusOK)
}

// CreateCategory handles POST /categories
// @Summary Create a new category
// @Description Create a new category
// @Tags categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.CreateCategoryRequest true "Create Category Request"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req dto.CreateCategoryRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "invalid request body"), fiber.StatusBadRequest)
	}

	// Validate request
	if err := validateCategory.Struct(req); err != nil {
		return utils.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()), fiber.StatusBadRequest)
	}

	category, err := h.categoryService.CreateCategory(req)
	if err != nil {
		if err == repositories.ErrCategoryExists {
			return utils.ErrorResponse(c, err, fiber.StatusConflict)
		}
		return utils.ErrorResponse(c, err, fiber.StatusInternalServerError)
	}

	return utils.SuccessResponse(c, category, "category created successfully", fiber.StatusCreated)
}

// UpdateCategory handles PUT /categories/:id
// @Summary Update category
// @Description Update category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Category ID"
// @Param request body dto.UpdateCategoryRequest true "Update Category Request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "invalid category ID"), fiber.StatusBadRequest)
	}

	var req dto.UpdateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "invalid request body"), fiber.StatusBadRequest)
	}

	// Validate request
	if err := validateCategory.Struct(req); err != nil {
		return utils.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()), fiber.StatusBadRequest)
	}

	category, err := h.categoryService.UpdateCategory(uint(id), req)
	if err != nil {
		if err == repositories.ErrCategoryNotFound {
			return utils.ErrorResponse(c, err, fiber.StatusNotFound)
		}
		if err == repositories.ErrCategoryExists {
			return utils.ErrorResponse(c, err, fiber.StatusConflict)
		}
		return utils.ErrorResponse(c, err, fiber.StatusInternalServerError)
	}

	return utils.SuccessResponse(c, category, "category updated successfully", fiber.StatusOK)
}

// DeleteCategory handles DELETE /categories/:id
// @Summary Delete category
// @Description Delete category by ID (soft delete)
// @Tags categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Category ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "invalid category ID"), fiber.StatusBadRequest)
	}

	if err := h.categoryService.DeleteCategory(uint(id)); err != nil {
		if err == repositories.ErrCategoryNotFound {
			return utils.ErrorResponse(c, err, fiber.StatusNotFound)
		}
		return utils.ErrorResponse(c, err, fiber.StatusInternalServerError)
	}

	return utils.SuccessResponse(c, nil, "category deleted successfully", fiber.StatusOK)
}

