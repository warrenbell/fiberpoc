package handlers

import (
	"time"

	fiber "github.com/gofiber/fiber/v2"
	"gitlab.com/sandstone2/fiberpoc/common/services"
	"go.uber.org/zap"
)

type AuthcHandler struct {
	authcService *services.AuthcServiceInterface
	logger       *zap.Logger
}

func NewAuthcHandler(authcService services.AuthcServiceInterface, logger *zap.Logger) *AuthcHandler {
	return &AuthcHandler{authcService: &authcService, logger: logger}
}

func (authcHandler *AuthcHandler) HandleRoot(c *fiber.Ctx) error {
	return c.Render("home", fiber.Map{
		"LoggedIn": false,
		"Error":    false,
		"Name":     "",
		"Email":    "",
	})
}

func (authcHandler *AuthcHandler) HandleLogin(c *fiber.Ctx) error {
	// 1. Generate a secure, random state string
	state, err := (*authcHandler.authcService).GenerateState()
	if err != nil {
		(*authcHandler.logger).Sugar().Errorf("Error: NKUM7E - Logging in. Error: %v", err)
		return c.Render("home", fiber.Map{
			"LoggedIn": false,
			"Error":    true,
			"Name":     "",
			"Email":    "",
		})
	}

	// 2. Store the state in a cookie (HttpOnly for security)
	c.Cookie(&fiber.Cookie{
		Name:     "oidc_state",
		Value:    state,
		Expires:  time.Now().Add(5 * time.Minute),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	})

	// 3. Redirect to the OIDC provider with the state
	url := (*authcHandler.authcService).GetOauthConfig().AuthCodeURL(state)
	return c.Redirect(url, fiber.StatusFound)
}

func (authcHandler *AuthcHandler) HandleOauthCallback(c *fiber.Ctx) error {
	expectedState := c.Cookies("oidc_state", "")
	receivedState := c.Query("state", "")

	if receivedState != expectedState {
		(*authcHandler.logger).Error("Error: 92ASWW - Logging in. CSRF attempted. States do not match.")
		return c.Render("home", fiber.Map{
			"LoggedIn": false,
			"Error":    true,
			"Name":     "",
			"Email":    "",
		})
	}

	code := c.Query("code", "")
	if code == "" {
		(*authcHandler.logger).Error("Error: TDUSAL - Getting oidc code from query string.")
		return c.Render("home", fiber.Map{
			"LoggedIn": false,
			"Error":    true,
			"Name":     "",
			"Email":    "",
		})
	}

	claims, jwt, err := (*authcHandler.authcService).ProcessOauth(code)
	if err != nil {
		(*authcHandler.logger).Sugar().Errorf("Error: 0GLO1T - Processing OAuth. Error: %v", err)
		return c.Render("home", fiber.Map{
			"LoggedIn": false,
			"Error":    true,
			"Name":     "",
			"Email":    "",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt_token",
		Value:    *jwt,
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	})

	return c.Render("home", fiber.Map{
		"LoggedIn": true,
		"Error":    false,
		"Name":     claims.Name,
		"Email":    claims.Email,
	})
}
