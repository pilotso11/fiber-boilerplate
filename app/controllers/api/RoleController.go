package api

import (
	"fiber-boilerplate/app/models"
	"fiber-boilerplate/database"

	"github.com/gofiber/fiber/v2"
)

// GetAllRoles
//
//	@Summary	Return all roles as JSON
//	@Router		/api/v1/roles [get]
//	@Produce	json
//	@Accept		json
//	@Success	200	{object}	[]models.RoleDto
func GetAllRoles(db *database.Database) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var Role []models.Role
		if response := db.Find(&Role); response.Error != nil {
			panic("Error occurred while retrieving roles from the database: " + response.Error.Error())
		}
		err := ctx.JSON(Role)
		if err != nil {
			panic("Error occurred when returning JSON of roles: " + err.Error())
		}
		return err
	}
}

// GetRole
//
//	@Summary	Return a single role as JSON
//	@Router		/api/v1/roles/{id} [get]
//	@Produce	json
//	@Accept		json
//	@Param		id	path		string	true	"Role ID"
//	@Success	200	{object}	models.RoleDto
func GetRole(db *database.Database) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		Role := new(models.Role)
		id := ctx.Params("id")
		if response := db.Find(&Role, id); response.Error != nil {
			panic("An error occurred when retrieving the role: " + response.Error.Error())
		}
		if Role.ID == 0 {
			// Send status not found
			err := ctx.SendStatus(fiber.StatusNotFound)
			if err != nil {
				panic("Cannot return status not found: " + err.Error())
			}
			// Set ID
			err = ctx.JSON(fiber.Map{
				"ID": id,
			})
			if err != nil {
				panic("Error occurred when returning JSON of a role: " + err.Error())
			}
			return err
		}
		err := ctx.JSON(Role)
		if err != nil {
			panic("Error occurred when returning JSON of a role: " + err.Error())
		}
		return err
	}
}

// AddRole
//
//	@Summary	Add a single role to the database
//	@Router		/api/v1/roles [post]
//	@Produce	json
//	@Accept		json
//	@Param		request	body		models.RoleDto	true	"Role data"
//	@Success	200		{object}	models.RoleDto
func AddRole(db *database.Database) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		Role := new(models.Role)
		if err := ctx.BodyParser(Role); err != nil {
			panic("An error occurred when parsing the new role: " + err.Error())
		}
		if response := db.Create(&Role); response.Error != nil {
			panic("An error occurred when storing the new role: " + response.Error.Error())
		}
		err := ctx.JSON(Role)
		if err != nil {
			panic("Error occurred when returning JSON of a role: " + err.Error())
		}
		return err
	}
}

// EditRole
//
//	@Summary	Edit a single role
//	@Router		/api/v1/roles/{id} [put]
//	@Produce	json
//	@Accept		json
//	@Param		id		path		string			true	"User ID"
//	@Param		request	body		models.RoleDto	true	"Role data"
//	@Success	200		{object}	models.RoleDto
func EditRole(db *database.Database) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		FormRole := new(models.Role)
		Role := new(models.Role)
		if err := ctx.BodyParser(FormRole); err != nil {
			panic("An error occurred when parsing the edited role: " + err.Error())
		}
		if response := db.Find(&Role, id); response.Error != nil {
			panic("An error occurred when retrieving the existing role: " + response.Error.Error())
		}
		// Role does not exist
		if Role.ID == 0 {
			err := ctx.SendStatus(fiber.StatusNotFound)
			if err != nil {
				panic("Cannot return status not found: " + err.Error())
			}
			err = ctx.JSON(fiber.Map{
				"ID": id,
			})
			if err != nil {
				panic("Error occurred when returning JSON of a role: " + err.Error())
			}
			return err
		}
		Role.Name = FormRole.Name
		Role.Description = FormRole.Description
		db.Save(&Role)

		err := ctx.JSON(Role)
		if err != nil {
			panic("Error occurred when returning JSON of a role: " + err.Error())
		}
		return err
	}
}

// DeleteRole
//
//	@Summary	Delete a single role
//	@Router		/api/v1/roles/{id} [delete]
//	@Produce	json
//	@Accept		json
//	@Param		id	path	string	true	"Role ID"
//	@Success	200
func DeleteRole(db *database.Database) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		var Role models.Role
		db.Find(&Role, id)
		if response := db.Find(&Role); response.Error != nil {
			panic("An error occurred when finding the role to be deleted: " + response.Error.Error())
		}
		db.Delete(&Role)

		err := ctx.JSON(fiber.Map{
			"ID":      id,
			"Deleted": true,
		})
		if err != nil {
			panic("Error occurred when returning JSON of a role: " + err.Error())
		}
		return err
	}
}
