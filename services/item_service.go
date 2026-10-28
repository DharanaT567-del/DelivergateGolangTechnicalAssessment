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

type ItemInput struct {
	Name        string  `json:"name" validate:"required"`
	Category    string  `json:"category" validate:"required"`
	Quantity    int     `json:"quantity" validate:"gte=0"`
	Price       float64 `json:"price" validate:"gt=0"`
	Description string  `json:"description"`
}

var ErrInvalidInput = errors.New("invalid input")

type ItemService struct{ items *repository.ItemRepository }

func NewItemService(items *repository.ItemRepository) *ItemService { return &ItemService{items: items} }

func (s *ItemService) List(search string, page, limit int, sort string) ([]models.Item, int64, error) {
	if page < 1 || limit < 1 || limit > 100 {
		return nil, 0, fmt.Errorf("%w: page must be positive and limit must be between 1 and 100", ErrInvalidInput)
	}
	return s.items.List(strings.TrimSpace(search), page, limit, sort)
}

func (s *ItemService) Get(id uint) (models.Item, error) { return s.items.GetByID(id) }

func (s *ItemService) Create(input ItemInput) (models.Item, error) {
	input.Name, input.Category = strings.TrimSpace(input.Name), strings.TrimSpace(input.Category)
	input.Description = strings.TrimSpace(input.Description)
	if err := utils.Validate(input); err != nil {
		return models.Item{}, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}
	item := models.Item{
		Name: input.Name, Category: input.Category,
		Quantity: input.Quantity, Price: input.Price, Description: input.Description,
	}
	if err := s.items.Create(&item); err != nil {
		return models.Item{}, fmt.Errorf("create item: %w", err)
	}
	return item, nil
}

func (s *ItemService) Update(id uint, input ItemInput) (models.Item, error) {
	input.Name, input.Category = strings.TrimSpace(input.Name), strings.TrimSpace(input.Category)
	input.Description = strings.TrimSpace(input.Description)
	if err := utils.Validate(input); err != nil {
		return models.Item{}, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}
	item, err := s.items.GetByID(id)
	if err != nil {
		return models.Item{}, err
	}
	item.Name, item.Category = input.Name, input.Category
	item.Quantity, item.Price = input.Quantity, input.Price
	item.Description = input.Description
	if err := s.items.Update(&item); err != nil {
		return models.Item{}, fmt.Errorf("update item: %w", err)
	}
	return item, nil
}

func (s *ItemService) Delete(id uint) error {
	err := s.items.Delete(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if err != nil {
		return fmt.Errorf("delete item: %w", err)
	}
	return nil
}
