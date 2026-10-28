package controllers

import (
	"strconv"

	"icecream-shop/services"
	"icecream-shop/utils"

	"github.com/gofiber/fiber/v2"
)

type ItemController struct{ items *services.ItemService }

func NewItemController(items *services.ItemService) *ItemController {
	return &ItemController{items: items}
}

type pageData struct {
	Items any   `json:"items"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}

// List returns paginated inventory and supports name search and allowlisted sort columns.
// @Summary List inventory items
// @Tags inventory
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Page size, maximum 100"
// @Param search query string false "Name search"
// @Param sort query string false "Sort column; prefix with - for descending"
// @Success 200 {object} utils.Envelope
// @Router /api/v1/items [get]
func (h *ItemController) List(c *fiber.Ctx) error {
	page, err := queryInt(c, "page", 1)
	if err != nil {
		return utils.Failure(c, fiber.StatusBadRequest, "page must be a positive integer")
	}
	limit, err := queryInt(c, "limit", 10)
	if err != nil {
		return utils.Failure(c, fiber.StatusBadRequest, "limit must be an integer")
	}
	items, total, err := h.items.List(c.Query("search"), page, limit, c.Query("sort"))
	if err != nil {
		return handleServiceError(c, err)
	}
	return utils.Success(c, fiber.StatusOK, "Items retrieved", pageData{Items: items, Total: total, Page: page, Limit: limit})
}

// Get returns an inventory item by ID.
// @Summary Get inventory item
// @Tags inventory
// @Security BearerAuth
// @Param id path int true "Item ID"
// @Success 200 {object} utils.Envelope
// @Failure 404 {object} utils.Envelope
// @Router /api/v1/items/{id} [get]
func (h *ItemController) Get(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return utils.Failure(c, fiber.StatusBadRequest, "Invalid item ID")
	}
	item, err := h.items.Get(id)
	if err != nil {
		return handleServiceError(c, err)
	}
	return utils.Success(c, fiber.StatusOK, "Item retrieved", item)
}

// Create adds an inventory item.
// @Summary Create inventory item
// @Tags inventory
// @Security BearerAuth
// @Param item body services.ItemInput true "Inventory item"
// @Success 201 {object} utils.Envelope
// @Failure 400 {object} utils.Envelope
// @Router /api/v1/items [post]
func (h *ItemController) Create(c *fiber.Ctx) error {
	var input services.ItemInput
	if err := c.BodyParser(&input); err != nil {
		return utils.Failure(c, fiber.StatusBadRequest, "Invalid request body")
	}
	item, err := h.items.Create(input)
	if err != nil {
		return handleServiceError(c, err)
	}
	return utils.Success(c, fiber.StatusCreated, "Item created", item)
}

// Update replaces the editable fields for an inventory item.
// @Summary Update inventory item
// @Tags inventory
// @Security BearerAuth
// @Param id path int true "Item ID"
// @Param item body services.ItemInput true "Inventory item"
// @Success 200 {object} utils.Envelope
// @Router /api/v1/items/{id} [put]
func (h *ItemController) Update(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return utils.Failure(c, fiber.StatusBadRequest, "Invalid item ID")
	}
	var input services.ItemInput
	if err := c.BodyParser(&input); err != nil {
		return utils.Failure(c, fiber.StatusBadRequest, "Invalid request body")
	}
	item, err := h.items.Update(id, input)
	if err != nil {
		return handleServiceError(c, err)
	}
	return utils.Success(c, fiber.StatusOK, "Item updated", item)
}

// Delete soft-deletes an inventory item.
// @Summary Delete inventory item
// @Tags inventory
// @Security BearerAuth
// @Param id path int true "Item ID"
// @Success 200 {object} utils.Envelope
// @Router /api/v1/items/{id} [delete]
func (h *ItemController) Delete(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return utils.Failure(c, fiber.StatusBadRequest, "Invalid item ID")
	}
	if err := h.items.Delete(id); err != nil {
		return handleServiceError(c, err)
	}
	return utils.Success(c, fiber.StatusOK, "Item deleted", nil)
}

func queryInt(c *fiber.Ctx, key string, fallback int) (int, error) {
	value := c.Query(key)
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(value)
	return parsed, err
}
