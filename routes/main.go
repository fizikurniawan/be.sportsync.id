package routes

import (
	"sportsync/bootstrap"
	"sportsync/delivery/http/middlewares"
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

	authMiddlerware := middlewares.NewAuthMiddleware(env)
	authGroup := app.Group("/auth")
	NewAuthRoute(env, timeout, db, authGroup)

	teamGroup := app.Group("/team").Use(authMiddlerware)
	NewTeamRoute(env, timeout, db, teamGroup)

	authGroup = app.Group("/ws")
	NewChatRoute(env, timeout, db, authGroup)
}
