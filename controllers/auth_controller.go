package controllers

import (
	"errors"

	"icecream-shop/services"
	"icecream-shop/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct{ auth *services.AuthService }

func NewAuthController(auth *services.AuthService) *AuthController {
	return &AuthController{auth: auth}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login authenticates a seeded user and returns a 24-hour bearer token.
// @Summary Login
// @Router /api/v1/auth/login [post]
func (h *AuthController) Login(c *fiber.Ctx) error {
	var request loginRequest
	if err := c.BodyParser(&request); err != nil || request.Username == "" || request.Password == "" {
		return utils.Failure(c, fiber.StatusBadRequest, "Username and password are required")
	}
	token, err := h.auth.Login(request.Username, request.Password)
	if errors.Is(err, services.ErrUnauthorized) {
		return utils.Failure(c, fiber.StatusUnauthorized, "Invalid username or password")
	}
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true, "message": "Login successful", "token": token,
	})
}

// Logout acknowledges logout. JWTs are stateless and expire after 24 hours.
// @Summary Logout
// @Router /api/v1/auth/logout [post]
func (h *AuthController) Logout(c *fiber.Ctx) error {
	return utils.Success(c, fiber.StatusOK, "Logged out", nil)
}
