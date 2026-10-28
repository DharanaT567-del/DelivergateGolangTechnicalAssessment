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
