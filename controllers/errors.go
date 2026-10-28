package controllers

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"icecream-shop/services"
	"icecream-shop/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func handleServiceError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return utils.Failure(c, fiber.StatusNotFound, "Resource not found")
	case errors.Is(err, services.ErrInvalidInput):
		return utils.Failure(c, fiber.StatusBadRequest, strings.TrimPrefix(err.Error(), services.ErrInvalidInput.Error()+": "))
	default:
		log.Printf("request failed: %v", err)
		return utils.Failure(c, fiber.StatusInternalServerError, "Internal server error")
	}
}

func parseID(c *fiber.Ctx) (uint, error) {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || id == 0 {
		return 0, errors.New("invalid id")
	}
	return uint(id), nil
}
