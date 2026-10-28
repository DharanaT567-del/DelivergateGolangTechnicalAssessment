package controllers

import (
	"icecream-shop/services"
	"icecream-shop/utils"

	"github.com/gofiber/fiber/v2"
)

type CustomerController struct{ customers *services.CustomerService }

func NewCustomerController(customers *services.CustomerService) *CustomerController {
	return &CustomerController{customers: customers}
}

// List returns paginated customer records and supports searching contact information.
// @Summary List customers
// @Tags customers
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Page size, maximum 100"
// @Param search query string false "Name, phone, or email search"
// @Success 200 {object} utils.Envelope
// @Router /api/v1/customers [get]
func (h *CustomerController) List(c *fiber.Ctx) error {
	page, err := queryInt(c, "page", 1)
	if err != nil {
		return utils.Failure(c, fiber.StatusBadRequest, "page must be a positive integer")
	}
	limit, err := queryInt(c, "limit", 10)
	if err != nil {
		return utils.Failure(c, fiber.StatusBadRequest, "limit must be an integer")
	}
	customers, total, err := h.customers.List(c.Query("search"), page, limit)
	if err != nil {
		return handleServiceError(c, err)
	}
	return utils.Success(c, fiber.StatusOK, "Customers retrieved", pageData{
		Items: customers, Total: total, Page: page, Limit: limit,
	})
}

// Get returns a customer by ID.
// @Summary Get customer
// @Tags customers
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {object} utils.Envelope
// @Router /api/v1/customers/{id} [get]
func (h *CustomerController) Get(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return utils.Failure(c, fiber.StatusBadRequest, "Invalid customer ID")
	}
	customer, err := h.customers.Get(id)
	if err != nil {
		return handleServiceError(c, err)
	}
	return utils.Success(c, fiber.StatusOK, "Customer retrieved", customer)
}

// Create adds a customer.
// @Summary Create customer
// @Tags customers
// @Security BearerAuth
// @Param customer body services.CustomerInput true "Customer details"
// @Success 201 {object} utils.Envelope
// @Failure 400 {object} utils.Envelope
// @Router /api/v1/customers [post]
func (h *CustomerController) Create(c *fiber.Ctx) error {
	var input services.CustomerInput
	if err := c.BodyParser(&input); err != nil {
		return utils.Failure(c, fiber.StatusBadRequest, "Invalid request body")
	}
	customer, err := h.customers.Create(input)
	if err != nil {
		return handleServiceError(c, err)
	}
	return utils.Success(c, fiber.StatusCreated, "Customer created", customer)
}

// Update replaces a customer's editable fields.
// @Summary Update customer
// @Tags customers
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Param customer body services.CustomerInput true "Customer details"
// @Success 200 {object} utils.Envelope
// @Router /api/v1/customers/{id} [put]
func (h *CustomerController) Update(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return utils.Failure(c, fiber.StatusBadRequest, "Invalid customer ID")
	}
	var input services.CustomerInput
	if err := c.BodyParser(&input); err != nil {
		return utils.Failure(c, fiber.StatusBadRequest, "Invalid request body")
	}
	customer, err := h.customers.Update(id, input)
	if err != nil {
		return handleServiceError(c, err)
	}
	return utils.Success(c, fiber.StatusOK, "Customer updated", customer)
}

// Delete soft-deletes a customer.
// @Summary Delete customer
// @Tags customers
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {object} utils.Envelope
// @Router /api/v1/customers/{id} [delete]
func (h *CustomerController) Delete(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return utils.Failure(c, fiber.StatusBadRequest, "Invalid customer ID")
	}
	if err := h.customers.Delete(id); err != nil {
		return handleServiceError(c, err)
	}
	return utils.Success(c, fiber.StatusOK, "Customer deleted", nil)
}
