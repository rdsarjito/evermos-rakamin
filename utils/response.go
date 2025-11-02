package utils

import (
	"github.com/gofiber/fiber/v2"
)

// Response represents standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// SuccessResponse sends a success response
func SuccessResponse(c *fiber.Ctx, data interface{}, message string, statusCode int) error {
	return c.Status(statusCode).JSON(Response{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *fiber.Ctx, err error, statusCode int) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Error:   err.Error(),
	})
}

