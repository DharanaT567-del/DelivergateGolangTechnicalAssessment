package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"icecream-shop/config"
	"icecream-shop/routes"
	"icecream-shop/seed"
	"icecream-shop/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("load environment: %w", err)
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return errors.New("JWT_SECRET must be set")
	}
	databaseName := os.Getenv("DB_NAME")
	if databaseName == "" {
		databaseName = "icecream.db"
	}
	db, err := config.OpenDatabase(databaseName)
	if err != nil {
		return err
	}
	if err := seed.SeedUsers(db); err != nil {
		return fmt.Errorf("seed users: %w", err)
	}

	app := fiber.New(fiber.Config{
		BodyLimit: 1024 * 1024,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := fiber.ErrInternalServerError.Message
			var fiberError *fiber.Error
			if errors.As(err, &fiberError) {
				code = fiberError.Code
				message = fiberError.Message
			} else {
				log.Printf("http error: %v", err)
			}
			return utils.Failure(c, code, message)
		},
	})
	routes.Register(app, db, secret)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	if err := app.Listen(":" + port); err != nil {
		return fmt.Errorf("start HTTP server: %w", err)
	}
	return nil
}
