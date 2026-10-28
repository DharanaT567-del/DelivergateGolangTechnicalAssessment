package middleware

import (
	"icecream-shop/utils"

	"github.com/gofiber/fiber/v2"
)

func AllowRoles(roles ...string) fiber.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok {
			return utils.Failure(c, fiber.StatusUnauthorized, "Authentication required")
		}
		if _, ok := allowed[role]; !ok {
			return utils.Failure(c, fiber.StatusForbidden, "Forbidden")
		}
		return c.Next()
	}
}
