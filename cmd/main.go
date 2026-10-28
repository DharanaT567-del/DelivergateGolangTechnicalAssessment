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
	// Load .env
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to load .env: %w", err)
	}

	// JWT Secret
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return errors.New("JWT_SECRET must be set in .env")
	}

	// Database name
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "icecream.db"
	}

	// Open database
	db, err := config.OpenDatabase(dbName)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Seed users
	if err := seed.SeedUsers(db); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit: 1024 * 1024,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "Internal Server Error"

			var fiberErr *fiber.Error
			if errors.As(err, &fiberErr) {
				code = fiberErr.Code
				message = fiberErr.Message
			} else {
				log.Println("ERROR:", err)
			}

			return utils.Failure(c, code, message)
		},
	})

	// Register routes
	routes.Register(app, db, secret)

	// Print all registered routes
	fmt.Println("Registered Routes")

	for _, route := range app.GetRoutes() {
		fmt.Printf("%-6s %s\n", route.Method, route.Path)
	}

	// Server port
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("\n Server running at http://localhost:%s\n", port)
	fmt.Println("Swagger UI: http://localhost:" + port + "/swagger/index.html")

	// Start server
	if err := app.Listen(":" + port); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
