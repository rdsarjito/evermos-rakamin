package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"
	"github.com/rdsarjito/evermos-rakamin/domain/dto"
	"github.com/rdsarjito/evermos-rakamin/middleware"
	"github.com/rdsarjito/evermos-rakamin/repositories"
	"github.com/rdsarjito/evermos-rakamin/services"
	"github.com/rdsarjito/evermos-rakamin/utils"
)

var validate = validator.New()

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(),
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register Request"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "invalid request body"), fiber.StatusBadRequest)
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		return utils.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()), fiber.StatusBadRequest)
	}

	// Register user
	response, err := h.authService.Register(req)
	if err != nil {
		if err == repositories.ErrUserExists {
			return utils.ErrorResponse(c, err, fiber.StatusConflict)
		}
		return utils.ErrorResponse(c, err, fiber.StatusInternalServerError)
	}

	return utils.SuccessResponse(c, response, "user registered successfully", fiber.StatusCreated)
}

// Login handles user login
// @Summary Login user
// @Description Authenticate user and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "invalid request body"), fiber.StatusBadRequest)
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		return utils.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()), fiber.StatusBadRequest)
	}

	// Login user
	response, err := h.authService.Login(req)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			return utils.ErrorResponse(c, err, fiber.StatusUnauthorized)
		}
		return utils.ErrorResponse(c, err, fiber.StatusInternalServerError)
	}

	return utils.SuccessResponse(c, response, "login successful", fiber.StatusOK)
}

// Profile gets authenticated user profile
// @Summary Get user profile
// @Description Get authenticated user profile information
// @Tags auth
// @Security Bearer
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /auth/profile [get]
func (h *AuthHandler) Profile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.NewError(fiber.StatusUnauthorized, "user not authenticated"), fiber.StatusUnauthorized)
	}

	profile, err := h.authService.GetUserProfile(userID)
	if err != nil {
		if err == repositories.ErrUserNotFound {
			return utils.ErrorResponse(c, err, fiber.StatusNotFound)
		}
		return utils.ErrorResponse(c, err, fiber.StatusInternalServerError)
	}

	return utils.SuccessResponse(c, profile, "profile retrieved successfully", fiber.StatusOK)
}

