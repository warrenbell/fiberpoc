package middleware

import (
	"strings"

	oidc "github.com/coreos/go-oidc"
	fiber "github.com/gofiber/fiber/v2"
	"gitlab.com/sandstone2/fiberpoc/common/models"
	"go.uber.org/zap"
)

func AuthcMiddleware(verifier *oidc.IDTokenVerifier, logger *zap.Logger) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// 1. Extract Bearer token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Error("Error: 3R7WBW - Getting authorization header.")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error 3R7WBW - Getting authorization header."})
		}

		rawToken := strings.TrimPrefix(authHeader, "Bearer ")

		// 2. Verify the token with Google
		idToken, err := verifier.Verify(c.Context(), rawToken)
		if err != nil {
			logger.Sugar().Errorf("Error: S2UU5K - Verifying token. Error: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Error S2UU5K - Verifying token."})
		}

		// 3. Extract claims
		claims := &models.Claims{}

		if err := idToken.Claims(&claims); err != nil {
			logger.Sugar().Errorf("Error: KH1NV5 - Parsing claims. Error: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error KH1NV5 - Parsing claims."})
		}

		// 4. Store user info in context
		c.Locals("user", claims)

		// 5. Proceed to next handler
		return c.Next()
	}
}
