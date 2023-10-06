package middlewares

import (
	"errors"
	"sportsync/bootstrap"
	"sportsync/internal"
	"sportsync/internal/tokenutil"

	"github.com/gofiber/fiber/v2"
)

// Middleware JWT function
func NewAuthMiddleware(env *bootstrap.Env) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return errors.New("authorization header is required")
		}

		userId, err := tokenutil.ExtractIDFromToken(authHeader, env.AccessTokenSecret)
		if err != nil {
			return internal.SendErrorRespond(c, 401, map[string][]string{"error": {err.Error()}})
		}

		if userId == "" {
			return internal.SendErrorRespond(c, 401, map[string][]string{"error": {"unauthorized"}})
		}

		c.Locals("userId", userId)

		return c.Next()
	}
}
