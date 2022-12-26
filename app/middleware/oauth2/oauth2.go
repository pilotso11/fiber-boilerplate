package oauth2

import (
	"fiber-boilerplate/auth/cognito"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/amazon"
	"github.com/markbates/goth/providers/auth0"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/okta"
	"github.com/shareed2k/goth_fiber"
	"go.uber.org/zap"
	"log"
	"net/url"
)

// OAUTH2 middleware
// User login via goth and goth_fiber
// login page at /oauth/login/{provider}
// logout at /oauth/logout
// The callback page is /oauth/callback/{provider}
//

const sessionUserKey = "user"

// RegisterOauth2
// Register the Oauth2 login/logout pages
func RegisterOauth2(cfg Config, app *fiber.App, session *session.Store) {
	CallbackURL := cfg.BaseURL + "/callback/"
	LoginURL := cfg.BaseURL + "/login/"
	LogoutURL := cfg.BaseURL + "/logout"

	if len(cfg.Secret) == 0 {
		log.Println("Empty OAUTH2 Secret, likely a misconfiguration")
	}

	callbackURL, _ := url.JoinPath(cfg.CallbackURL, CallbackURL, cfg.Provider)
	switch cfg.Provider {
	case "google":
		goth.UseProviders(google.New(cfg.ClientKey, cfg.Secret, callbackURL))
		break
	case "amazon":
		goth.UseProviders(amazon.New(cfg.ClientKey, cfg.Secret, callbackURL))
		break
	case "auth0":
		goth.UseProviders(auth0.New(cfg.ClientKey, cfg.Secret, cfg.Auth0Domain, callbackURL))
		break
	case "okta":
		goth.UseProviders(okta.New(cfg.ClientKey, cfg.Secret, cfg.OrgURL, callbackURL))
		break
	case "github":
		goth.UseProviders(github.New(cfg.ClientKey, cfg.Secret, callbackURL))
		break
	case "cognito":
		goth.UseProviders(cognito.New(cfg.ClientKey, cfg.Secret, cfg.OrgURL, callbackURL))
		break
	default:
		log.Fatalf("Unknown OATH2 provider %s: \n", cfg.Provider)
	}

	app.Get(LoginURL+":provider", goth_fiber.BeginAuthHandler)

	app.Get(CallbackURL+":provider", func(ctx *fiber.Ctx) error {
		user, err := goth_fiber.CompleteUserAuth(ctx, goth_fiber.CompleteUserAuthOptions{ShouldLogout: false})
		if err != nil {
			panic(err)
		}
		// Save the user in the session
		sess, _ := session.Get(ctx)
		sess.Set(sessionUserKey, user)
		_ = sess.Save()
		err = ctx.Redirect(cfg.AfterLoginRedirectURL)
		return err
	})

	app.Get(LogoutURL, func(ctx *fiber.Ctx) error {
		if err := goth_fiber.Logout(ctx); err != nil {
			panic(err)
		}
		sess, _ := session.Get(ctx)
		_ = sess.Regenerate() // destroy session on logout

		err := ctx.Redirect(cfg.AfterLogoutRedirectURL)
		return err
	})
}

// New
// Creates new oauth2 middleware
func New(app *fiber.App, sessionStore *session.Store, config ...Config) fiber.Handler {
	var cfg Config
	if len(config) < 1 {
		cfg = configDefault()
	} else {
		cfg = configDefault(config[0])
	}

	// register the login & logout handlers
	RegisterOauth2(cfg, app, sessionStore)

	return func(ctx *fiber.Ctx) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(ctx) {
			return ctx.Next()
		}

		// Check for user in session
		sess, err := sessionStore.Get(ctx)
		if err != nil {
			return err
		}

		userid := ""
		username := ""
		email := ""
		if sess != nil {
			user := sess.Get(sessionUserKey)
			if user != nil {
				gu := user.(goth.User)
				userid = gu.UserID
				username = gu.Name
				email = gu.Email
				if username == "" {
					username = email
				}
			}
		}
		zap.S().Debugf("id: %s, name: %s, email: %s", userid, username, email)

		ctx.Locals(cfg.ContextUserID, userid)
		ctx.Locals(cfg.ContextUsername, username)
		ctx.Locals(cfg.ContextEmail, email)

		return ctx.Next()
	}

}
