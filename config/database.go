package config

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"icecream-shop/models"
)

func OpenDatabase(name string) (*gorm.DB, error) {
	if name == "" {
		return nil, fmt.Errorf("database name must not be empty")
	}

	db, err := gorm.Open(sqlite.Open(name), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Item{}, &models.Customer{}); err != nil {
		return nil, fmt.Errorf("migrate database: %w", err)
	}
	return db, nil
}
