package routes

import (
	"sportsync/bootstrap"
	"sportsync/mongo"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db mongo.Database, app *fiber.App) {
	// routes
	publicRouter := app.Group("")
	publicRouter.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "OK",
		})
	})

	NewAuthRoute(env, timeout, db, app)
}
