package services

import (
	"errors"
	"fmt"
	"strings"

	"icecream-shop/models"
	"icecream-shop/repository"
	"icecream-shop/utils"

	"gorm.io/gorm"
)

type CustomerInput struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone" validate:"required"`
	Email     string `json:"email" validate:"omitempty,email"`
	Address   string `json:"address"`
}

type CustomerService struct {
	customers *repository.CustomerRepository
}

func NewCustomerService(customers *repository.CustomerRepository) *CustomerService {
	return &CustomerService{customers: customers}
}

func (s *CustomerService) List(search string, page, limit int) ([]models.Customer, int64, error) {
	if page < 1 || limit < 1 || limit > 100 {
		return nil, 0, fmt.Errorf("%w: page must be positive and limit must be between 1 and 100", ErrInvalidInput)
	}
	return s.customers.List(strings.TrimSpace(search), page, limit)
}

func (s *CustomerService) Get(id uint) (models.Customer, error) { return s.customers.GetByID(id) }

func (s *CustomerService) Create(input CustomerInput) (models.Customer, error) {
	trimCustomerInput(&input)
	if err := utils.Validate(input); err != nil {
		return models.Customer{}, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}
	customer := customerFromInput(input)
	if err := s.customers.Create(&customer); err != nil {
		return models.Customer{}, fmt.Errorf("create customer: %w", err)
	}
	return customer, nil
}

func (s *CustomerService) Update(id uint, input CustomerInput) (models.Customer, error) {
	trimCustomerInput(&input)
	if err := utils.Validate(input); err != nil {
		return models.Customer{}, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}
	customer, err := s.customers.GetByID(id)
	if err != nil {
		return models.Customer{}, err
	}
	updated := customerFromInput(input)
	customer.FirstName, customer.LastName = updated.FirstName, updated.LastName
	customer.Phone, customer.Email, customer.Address = updated.Phone, updated.Email, updated.Address
	if err := s.customers.Update(&customer); err != nil {
		return models.Customer{}, fmt.Errorf("update customer: %w", err)
	}
	return customer, nil
}

func (s *CustomerService) Delete(id uint) error {
	err := s.customers.Delete(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if err != nil {
		return fmt.Errorf("delete customer: %w", err)
	}
	return nil
}

func customerFromInput(input CustomerInput) models.Customer {
	return models.Customer{
		FirstName: input.FirstName, LastName: input.LastName,
		Phone: input.Phone, Email: input.Email, Address: input.Address,
	}
}

func trimCustomerInput(input *CustomerInput) {
	input.FirstName = strings.TrimSpace(input.FirstName)
	input.LastName = strings.TrimSpace(input.LastName)
	input.Phone = strings.TrimSpace(input.Phone)
	input.Email = strings.TrimSpace(input.Email)
	input.Address = strings.TrimSpace(input.Address)
}
