package routes

import (
	"embed"

	"icecream-shop/controllers"
	"icecream-shop/middleware"
	"icecream-shop/repository"
	"icecream-shop/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

//go:embed swagger.json
var swaggerDocument embed.FS

func Register(app *fiber.App, db *gorm.DB, jwtSecret string) {
	app.Use(middleware.Recovery(), middleware.Logging(), cors.New())
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		data, err := swaggerDocument.ReadFile("swagger.json")
		if err != nil {
			return err
		}
		return c.Type("json").Send(data)
	})
	app.Get("/swagger/index.html", func(c *fiber.Ctx) error {
		return c.Type("html").SendString(`<!doctype html><html><head><title>Ice Cream Shop API</title><link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css"></head><body><div id="swagger-ui"></div><script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script><script>SwaggerUIBundle({url:"/swagger/doc.json",dom_id:"#swagger-ui"});</script></body></html>`)
	})

	userRepository := repository.NewUserRepository(db)
	itemRepository := repository.NewItemRepository(db)
	customerRepository := repository.NewCustomerRepository(db)
	authController := controllers.NewAuthController(services.NewAuthService(userRepository, jwtSecret))
	itemController := controllers.NewItemController(services.NewItemService(itemRepository))
	customerController := controllers.NewCustomerController(services.NewCustomerService(customerRepository))

	api := app.Group("/api/v1")
	auth := api.Group("/auth")
	auth.Post("/login", authController.Login)

	protected := api.Group("", middleware.Authenticate(jwtSecret))
	protected.Post("/auth/logout", authController.Logout)

	items := protected.Group("/items")
	items.Get("/", itemController.List)
	items.Get("/:id", itemController.Get)
	items.Post("/", middleware.AllowRoles("Owner", "Manager"), itemController.Create)
	items.Put("/:id", middleware.AllowRoles("Owner", "Manager", "Cashier"), itemController.Update)
	items.Delete("/:id", middleware.AllowRoles("Owner"), itemController.Delete)

	customers := protected.Group("/customers")
	customers.Get("/", customerController.List)
	customers.Get("/:id", customerController.Get)
	customers.Post("/", middleware.AllowRoles("Owner", "Manager", "Cashier"), customerController.Create)
	customers.Put("/:id", middleware.AllowRoles("Owner", "Manager", "Cashier"), customerController.Update)
	customers.Delete("/:id", middleware.AllowRoles("Owner"), customerController.Delete)

}
