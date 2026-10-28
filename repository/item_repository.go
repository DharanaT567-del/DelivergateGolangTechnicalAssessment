package repository

import (
	"icecream-shop/models"

	"gorm.io/gorm"
)

type ItemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) *ItemRepository { return &ItemRepository{db: db} }

func (r *ItemRepository) List(search string, page, limit int, sort string) ([]models.Item, int64, error) {
	query := r.db.Model(&models.Item{})
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	sortColumns := map[string]string{
		"name": "name", "category": "category", "price": "price",
		"quantity": "quantity", "created_at": "created_at",
		"-name": "name DESC", "-category": "category DESC", "-price": "price DESC",
		"-quantity": "quantity DESC", "-created_at": "created_at DESC",
	}
	order, ok := sortColumns[sort]
	if !ok {
		order = "id"
	}
	var items []models.Item
	err := query.Order(order).Offset((page - 1) * limit).Limit(limit).Find(&items).Error
	return items, total, err
}

func (r *ItemRepository) GetByID(id uint) (models.Item, error) {
	var item models.Item
	err := r.db.First(&item, id).Error
	return item, err
}

func (r *ItemRepository) Create(item *models.Item) error { return r.db.Create(item).Error }

func (r *ItemRepository) Update(item *models.Item) error {
	return r.db.Model(item).Select("Name", "Category", "Quantity", "Price", "Description").Updates(item).Error
}

func (r *ItemRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Item{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
