package routes

import (
	_ "fiber-boilerplate/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger" // swagger handler
)

// RegisterSwagger func for describe group of API Docs routes.
func RegisterSwagger(a *fiber.App) {
	// Create routes group.
	//route := a.Group("/swagger")

	// Routes for GET method:
	a.Get("/swagger/*", swagger.HandlerDefault)    // get one user by ID
	a.Get("/swagger-ui/*", swagger.HandlerDefault) // get one user by ID

}
