package routes

import (
	Controllers "fiber-boilerplate/app/controllers/web"
	"fiber-boilerplate/database"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	hashing "github.com/thomasvvugt/fiber-hashing"
)

func RegisterWeb(web fiber.Router, session *session.Store, db *database.Database, hasher hashing.Driver) {
	// Homepage
	web.Get("/", Controllers.Index(session))

	// Panic test route, this brings up an error
	web.Get("/panic", func(ctx *fiber.Ctx) error {
		panic("We have aa panic error!")
	})

	// Test to load static, compiled assets
	web.Get("/test", func(c *fiber.Ctx) error {
		return c.Render("test", fiber.Map{})
	})

	// Test to load static, compiled assets
	web.Get("/secured", Controllers.Secured(session))

	// Make a new hash
	web.Get("/hash/*", func(ctx *fiber.Ctx) error {
		hash, err := hasher.CreateHash(ctx.Params("*"))
		if err != nil {
			log.Fatalf("Error when creating hash: %v", err)
		}
		if err := ctx.SendString(hash); err != nil {
			panic(err.Error())
		}
		return err
	})

	// Auth routes
	web.Get("/login", Controllers.ShowLoginForm())
	web.Post("/login", Controllers.PostLoginForm(hasher, session, db))
	web.Get("/logout", Controllers.PostLogoutForm(session))
}
