package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/markbates/goth"
)

func Index(session *session.Store) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		// Bind data to template
		bind := fiber.Map{
			"name": "Fiber",
		}

		BindAuthenticationData(ctx, session, bind)

		// Copy auth details from goth.User to session

		// Render template
		err := ctx.Render("index", bind)
		if err != nil {
			err2 := ctx.Status(500).SendString(err.Error())
			if err2 != nil {
				panic(err2.Error())
			}
		}
		return err
	}
}

func BindAuthenticationData(ctx *fiber.Ctx, session *session.Store, bind fiber.Map) {
	sess, err := session.Get(ctx)
	if err != nil {
		panic(err)
	}
	bind["auth"] = false
	if IsAuthenticatedOauth(session, ctx) {
		user := sess.Get("user")
		if user != nil {
			gu := user.(goth.User)
			userID := gu.UserID
			email := gu.Email
			name := gu.Name
			if name == "" {
				name = email
			}
			bind["username"] = name
			bind["userid"] = userID
			bind["email"] = email
			bind["auth"] = true
		}
	} else if IsAuthenticated(session, ctx) {
		// Get User ID from session store
		userID := sess.Get("userid")
		email := sess.Get("email")
		name := sess.Get("username")
		bind["username"] = name
		bind["userid"] = userID
		bind["email"] = email
		bind["auth"] = true

	}
}

func IsAuthenticatedOauth(store *session.Store, ctx *fiber.Ctx) bool {
	sess, err := store.Get(ctx)
	if err != nil {
		return false
	}

	user, ok := sess.Get("user").(goth.User)
	if ok && user.UserID > "" {
		return true
	}
	return false
}
