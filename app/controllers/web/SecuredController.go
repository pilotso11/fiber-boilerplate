package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func Secured(session *session.Store) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userid := ctx.Locals("userid")
		if userid == nil || userid == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}

		// Render template
		err := ctx.Render("secured", GetMapOfSessionData(ctx, session))
		if err != nil {
			err2 := ctx.Status(500).SendString(err.Error())
			if err2 != nil {
				panic(err2.Error())
			}
		}
		return err

	}
}

func GetMapOfSessionData(ctx *fiber.Ctx, session *session.Store) fiber.Map {
	data := fiber.Map{}
	sess, _ := session.Get(ctx)
	if sess != nil {
		for _, k := range sess.Keys() {
			data[k] = sess.Get(k)
		}
	}
	BindAuthenticationData(ctx, session, data)

	return data
}
