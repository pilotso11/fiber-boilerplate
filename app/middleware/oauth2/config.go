package oauth2

import (
	"github.com/gofiber/fiber/v2"
)

type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// keys to store the username in Locals
	ContextUserID          string
	ContextUsername        string
	ContextEmail           string
	Provider               string
	BaseURL                string
	CallbackURL            string
	ClientKey              string
	Secret                 string
	OrgURL                 string
	Auth0Domain            string
	AfterLoginRedirectURL  string
	AfterLogoutRedirectURL string
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	Next:                   nil,
	ContextUserID:          "userid",
	ContextUsername:        "username",
	ContextEmail:           "email",
	Provider:               "",
	BaseURL:                "/oauth",
	CallbackURL:            "http://localhost:8080",
	ClientKey:              "",
	Secret:                 "",
	OrgURL:                 "",
	Auth0Domain:            "",
	AfterLoginRedirectURL:  "/",
	AfterLogoutRedirectURL: "http://localhost:8080/",
}

// Helper function to set default values
func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	// Set default values
	if cfg.Next == nil {
		cfg.Next = ConfigDefault.Next
	}
	if cfg.ContextUserID == "" {
		cfg.ContextUserID = ConfigDefault.ContextUserID
	}
	if cfg.ContextUsername == "" {
		cfg.ContextUsername = ConfigDefault.ContextUsername
	}
	if cfg.ContextEmail == "" {
		cfg.ContextEmail = ConfigDefault.ContextEmail
	}
	if cfg.Provider == "" {
		cfg.Provider = ConfigDefault.Provider
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = ConfigDefault.BaseURL
	}
	if cfg.CallbackURL == "" {
		cfg.CallbackURL = ConfigDefault.CallbackURL
	}
	if cfg.ClientKey == "" {
		cfg.ClientKey = ConfigDefault.ClientKey
	}
	if cfg.Secret == "" {
		cfg.Secret = ConfigDefault.Secret
	}
	if cfg.OrgURL == "" {
		cfg.OrgURL = ConfigDefault.OrgURL
	}
	if cfg.AfterLoginRedirectURL == "" {
		cfg.AfterLoginRedirectURL = ConfigDefault.AfterLoginRedirectURL
	}
	if cfg.AfterLogoutRedirectURL == "" {
		cfg.AfterLogoutRedirectURL = ConfigDefault.AfterLogoutRedirectURL
	}
	return cfg
}
