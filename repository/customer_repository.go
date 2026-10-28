package repository

import (
	"icecream-shop/models"

	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) List(search string, page, limit int) ([]models.Customer, int64, error) {
	query := r.db.Model(&models.Customer{})
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("first_name LIKE ? OR last_name LIKE ? OR phone LIKE ? OR email LIKE ?", like, like, like, like)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var customers []models.Customer
	err := query.Order("id").Offset((page - 1) * limit).Limit(limit).Find(&customers).Error
	return customers, total, err
}

func (r *CustomerRepository) GetByID(id uint) (models.Customer, error) {
	var customer models.Customer
	err := r.db.First(&customer, id).Error
	return customer, err
}

func (r *CustomerRepository) Create(customer *models.Customer) error {
	return r.db.Create(customer).Error
}

func (r *CustomerRepository) Update(customer *models.Customer) error {
	return r.db.Model(customer).Select("FirstName", "LastName", "Phone", "Email", "Address").Updates(customer).Error
}

func (r *CustomerRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Customer{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
