package main

import (
	"log"
	"main/src/db"

	"github.com/gofiber/fiber/v2"
)

func main() {
	db.Connect()

	app := fiber.New()
	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
