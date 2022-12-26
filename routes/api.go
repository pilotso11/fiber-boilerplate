package routes

import (
	Controller "fiber-boilerplate/app/controllers/api"
	"fiber-boilerplate/database"
	hashing "github.com/thomasvvugt/fiber-hashing"

	"github.com/gofiber/fiber/v2"
)

func RegisterAPI(api fiber.Router, db *database.Database, hasher hashing.Driver) {
	registerRoles(api, db)
	registerUsers(api, db, hasher)
}

func registerRoles(api fiber.Router, db *database.Database) {
	roles := api.Group("/roles")

	roles.Get("/", Controller.GetAllRoles(db))
	roles.Get("/:id", Controller.GetRole(db))
	roles.Post("/", Controller.AddRole(db))
	roles.Put("/:id", Controller.EditRole(db))
	roles.Delete("/:id", Controller.DeleteRole(db))
}

func registerUsers(api fiber.Router, db *database.Database, hasher hashing.Driver) {
	users := api.Group("/users")

	users.Get("/", Controller.GetAllUsers(db))
	users.Get("/:id", Controller.GetUser(db))
	users.Post("/", Controller.AddUser(db, hasher))
	users.Put("/:id", Controller.EditUser(db, hasher))
	users.Delete("/:id", Controller.DeleteUser(db))
}
