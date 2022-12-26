package api

import (
	"fiber-boilerplate/app/models"
	"fiber-boilerplate/database"

	"github.com/gofiber/fiber/v2"
	hashing "github.com/thomasvvugt/fiber-hashing"
	_ "gorm.io/gorm"
)

// GetAllUsers
//
//	@Summary	Return all users as JSON
//	@Produce	json
//	@Router		/api/v1/users [get]
//	@Success	200	{object}	[]models.UserDto
func GetAllUsers(db *database.Database) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var Users []models.User
		if response := db.Find(&Users); response.Error != nil {
			panic("Error occurred while retrieving users from the database: " + response.Error.Error())
		}
		// Match roles to users
		for index, User := range Users {
			if User.RoleID != 0 {
				Role := new(models.Role)
				if response := db.Find(&Role, User.RoleID); response.Error != nil {
					panic("An error occurred when retrieving the role: " + response.Error.Error())
				}
				if Role.ID != 0 {
					Users[index].Role = *Role
				}
			}
		}
		err := ctx.JSON(Users)
		if err != nil {
			panic("Error occurred when returning JSON of users: " + err.Error())
		}
		return err
	}
}

// GetUser
//
//	@Summary	Return a single user as JSON
//	@Produce	json
//	@Router		/api/v1/users/{id} [get]
//	@Param		id	path		string	true	"User ID"
//	@Success	200	{object}	models.UserDto
func GetUser(db *database.Database) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		User := new(models.User)
		id := ctx.Params("id")
		if response := db.Find(&User, id); response.Error != nil {
			panic("An error occurred when retrieving the user: " + response.Error.Error())
		}
		if User.ID == 0 {
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
		// Match role to user
		if User.RoleID != 0 {
			Role := new(models.Role)
			if response := db.Find(&Role, User.RoleID); response.Error != nil {
				panic("An error occurred when retrieving the role: " + response.Error.Error())
			}
			if Role.ID != 0 {
				User.Role = *Role
			}
		}
		err := ctx.JSON(User)
		if err != nil {
			panic("Error occurred when returning JSON of a user: " + err.Error())
		}
		return err
	}
}

// AddUser
//
//	@Summary	Add a single user to the database
//	@Produce	json
//	@Accept		json
//	@Router		/api/v1/users [post]
//	@Param		request	body		models.UserDto	true	"User data"
//	@Success	200		{object}	models.UserDto
func AddUser(db *database.Database, hasher hashing.Driver) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		FormUser := new(models.UserDto)
		if err := ctx.BodyParser(FormUser); err != nil {
			panic("An error occurred when parsing the new user: " + err.Error())
		}
		User := new(models.User)
		User.Name = FormUser.Name
		User.Email = FormUser.Email
		User.RoleID = FormUser.RoleID
		User.Password, _ = hasher.CreateHash(FormUser.Password)
		if response := db.Create(&User); response.Error != nil {
			panic("An error occurred when storing the new user: " + response.Error.Error())
		}
		// Match role to user
		if User.RoleID != 0 {
			Role := new(models.Role)
			if response := db.Find(&Role, User.RoleID); response.Error != nil {
				panic("An error occurred when retrieving the role" + response.Error.Error())
			}
			if Role.ID != 0 {
				User.Role = *Role
			}
		}
		err := ctx.JSON(User)
		if err != nil {
			panic("Error occurred when returning JSON of a user: " + err.Error())
		}
		return err
	}
}

// EditUser
//
//	@Summary	Edit a single user
//	@Produce	json
//	@Accept		json
//	@Router		/api/v1/users/{id} [put]
//	@Param		id		path		string			true	"User ID"
//	@Param		request	body		models.UserDto	true	"User data"
//	@Success	200		{object}	models.UserDto
func EditUser(db *database.Database, hasher hashing.Driver) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		FormUser := new(models.UserDto)
		User := new(models.User)
		if err := ctx.BodyParser(FormUser); err != nil {
			panic("An error occurred when parsing the edited user: " + err.Error())
		}
		if response := db.Find(&User, id); response.Error != nil {
			panic("An error occurred when retrieving the existing user: " + response.Error.Error())
		}
		// User does not exist
		if User.ID == 0 {
			err := ctx.SendStatus(fiber.StatusNotFound)
			if err != nil {
				panic("Cannot return status not found: " + err.Error())
			}
			err = ctx.JSON(fiber.Map{
				"ID": id,
			})
			if err != nil {
				panic("Error occurred when returning JSON of a user: " + err.Error())
			}
			return err
		}
		User.Name = FormUser.Name
		User.Email = FormUser.Email
		User.RoleID = FormUser.RoleID
		if len(FormUser.Password) > 0 {
			User.Password, _ = hasher.CreateHash(FormUser.Password)
		}
		// Match role to user
		if User.RoleID != 0 {
			Role := new(models.Role)
			if response := db.Find(&Role, User.RoleID); response.Error != nil {
				panic("An error occurred when retrieving the role" + response.Error.Error())
			}
			if Role.ID != 0 {
				User.Role = *Role
			}
		}
		// Save user
		db.Save(&User)

		err := ctx.JSON(User)
		if err != nil {
			panic("Error occurred when returning JSON of a user: " + err.Error())
		}
		return err
	}
}

// DeleteUser
//
//	@Summary	Delete a single user
//	@Produce	json
//	@Accept		json
//	@Router		/api/v1/users/{id} [delete]
//	@Param		id	path	string	true	"User ID"
//	@Success	200
func DeleteUser(db *database.Database) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		var User models.User
		db.Find(&User, id)
		if response := db.Find(&User); response.Error != nil {
			panic("An error occurred when finding the user to be deleted" + response.Error.Error())
		}
		db.Delete(&User)

		err := ctx.JSON(fiber.Map{
			"ID":      id,
			"Deleted": true,
		})
		if err != nil {
			panic("Error occurred when returning JSON of a user: " + err.Error())
		}
		return err
	}
}
