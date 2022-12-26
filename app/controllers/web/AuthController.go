package web

import (
	"fiber-boilerplate/database"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	hashing "github.com/thomasvvugt/fiber-hashing"
	"log"
)

func IsAuthenticated(session *session.Store, ctx *fiber.Ctx) (authenticated bool) {
	store, err := session.Get(ctx)
	if err != nil {
		panic(err)
	}
	// Get User ID from session store
	userID, correct := store.Get("userid").(uint)
	if !correct {
		userID = 0
	}
	auth := false
	if userID > 0 {
		auth = true
	}
	return auth
}

func ShowLoginForm() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		err := ctx.Render("login", fiber.Map{})
		if err != nil {
			if err2 := ctx.Status(500).SendString(err.Error()); err2 != nil {
				panic(err2.Error())
			}
		}
		return err
	}
}

// PostLoginForm
//
//	@Summary	Login
//	@Produce	json
//	@Accept		json
//	@Router		/login [post]
//	@Param		username	formData	string	true	"User ID"
//	@Param		password	formData	string	true	"Password"
//	@Success	200
func PostLoginForm(hasher hashing.Driver, session *session.Store, db *database.Database) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		username := ctx.FormValue("username")
		// Find user
		user, err := FindUserByUsername(db, username)
		if err != nil {
			log.Fatalf("Error when finding user: %v", err)
		}

		// Check if password matches hash
		if hasher != nil {
			password := ctx.FormValue("password")
			match, err := hasher.MatchHash(password, user.Password)
			if err != nil {
				err := ctx.SendString("The entered details do not match our records.")
				if err != nil {
					err.Error()
				}
				return nil
			}
			if match {
				store, err := session.Get(ctx)
				if err != nil {
					panic(err)
				}
				//goland:noinspection ALL
				// Set the user ID in the session store
				store.Set("userid", user.ID)
				store.Set("email", user.Email)
				store.Set("username", user.Name)
				_ = store.Save()
				fmt.Printf("User set in session store with ID: %v\n", user.ID)
				err = ctx.Redirect("/")
				if err != nil {
					//if err := ctx.SendString("You should be logged in successfully!"); err != nil {
					panic(err.Error())
				}
			} else {
				if err := ctx.SendString("The entered details do not match our records."); err != nil {
					panic(err.Error())
				}
			}
		} else {
			panic("Hash provider was not set")
		}
		return nil
	}
}

// PostLogoutForm
//
//	@Summary	Logout
//	@Produce	json
//	@Accept		json
//	@Router		/logout [get]
//	@Success	200
func PostLogoutForm(session *session.Store) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if IsAuthenticated(session, ctx) {
			store, err := session.Get(ctx)
			if err != nil {
				panic(err)
			}
			_ = store.Regenerate()
		}
		return ctx.Redirect("/")

	}
}
