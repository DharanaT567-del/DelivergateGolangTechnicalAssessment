package seed

import (
	"fmt"

	"icecream-shop/models"
	"icecream-shop/repository"
	"icecream-shop/utils"

	"gorm.io/gorm"
)

type account struct {
	username string
	password string
	role     string
}

var defaultAccounts = []account{
	{username: "owner", password: "owner123", role: "Owner"},
	{username: "manager", password: "manager123", role: "Manager"},
	{username: "cashier", password: "cashier123", role: "Cashier"},
}

func SeedUsers(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		users := repository.NewUserRepository(tx)
		count, err := users.Count()
		if err != nil {
			return fmt.Errorf("count users before seeding: %w", err)
		}
		if count != 0 {
			return nil
		}
		for _, account := range defaultAccounts {
			hash, err := utils.HashPassword(account.password)
			if err != nil {
				return fmt.Errorf("hash seed password for %s: %w", account.username, err)
			}
			user := models.User{Username: account.username, Password: hash, Role: account.role}
			if err := users.Create(&user); err != nil {
				return fmt.Errorf("seed user %s: %w", account.username, err)
			}
		}
		return nil
	})
}
