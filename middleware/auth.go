package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rdsarjito/evermos-rakamin/helpers"
	"github.com/rdsarjito/evermos-rakamin/utils"
)

// AuthMiddleware validates JWT token and attaches user info to context
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.ErrorResponse(c, fiber.NewError(fiber.StatusUnauthorized, "authorization header required"), fiber.StatusUnauthorized)
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.ErrorResponse(c, fiber.NewError(fiber.StatusUnauthorized, "invalid authorization header format"), fiber.StatusUnauthorized)
		}

		token := parts[1]

		// Validate token
		claims, err := helpers.ValidateToken(token)
		if err != nil {
			statusCode := fiber.StatusUnauthorized
			if err == helpers.ErrExpiredToken {
				statusCode = fiber.StatusUnauthorized
			}
			return utils.ErrorResponse(c, err, statusCode)
		}

		// Attach user info to context
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("user_role", claims.Role)

		return c.Next()
	}
}

// GetUserID extracts user ID from context
func GetUserID(c *fiber.Ctx) uint {
	if userID, ok := c.Locals("user_id").(uint); ok {
		return userID
	}
	return 0
}

// GetUserRole extracts user role from context
func GetUserRole(c *fiber.Ctx) string {
	if role, ok := c.Locals("user_role").(string); ok {
		return role
	}
	return ""
}

