package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"sportsync/bootstrap"
	"sportsync/routes"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := bootstrap.App()
	env := app.Env

	db := app.Mongo.Database(env.DBName)
	defer app.CloseDBConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second

	fiberApp := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		Prefork:     false,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			status := http.StatusText(code)

			// Send custom error page
			if err != nil {
				// In case the SendFile fails
				return c.Status(code).JSON(fiber.Map{
					"status":  status,
					"code":    code,
					"message": err.Error(),
				})
			}
			return nil
		},
	})

	// middleware
	fiberApp.Use(logger.New())
	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins: env.AllowOrigins,
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	routes.Setup(env, timeout, db, fiberApp)

	// running API
	apiPort := "4000"
	if env.ApiPort != "" {
		apiPort = env.ApiPort
	}
	apiServer := fmt.Sprintf(":%s", apiPort)
	log.Fatal(fiberApp.Listen(apiServer))
}
