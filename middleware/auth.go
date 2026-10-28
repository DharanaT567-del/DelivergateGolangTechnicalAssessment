package middleware

import (
	"strings"

	"icecream-shop/utils"

	"github.com/gofiber/fiber/v2"
)

func Authenticate(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := strings.TrimSpace(c.Get(fiber.HeaderAuthorization))
		parts := strings.Fields(header)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return utils.Failure(c, fiber.StatusUnauthorized, "Missing or invalid Authorization header")
		}
		claims, err := utils.ParseToken([]byte(secret), parts[1])
		if err != nil {
			return utils.Failure(c, fiber.StatusUnauthorized, "Invalid or expired token")
		}
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)
		return c.Next()
	}
}
