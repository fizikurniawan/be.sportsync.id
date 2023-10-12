package routes

import (
	"sportsync/bootstrap"
	wsHandler "sportsync/delivery/websocket"
	"sportsync/domain"
	"sportsync/internal"
	"sportsync/internal/tokenutil"
	"sportsync/mongo"
	"sportsync/repository"
	"sportsync/usecase"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func NewChatRoute(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group fiber.Router) {
	cr := repository.NewChatRepository(db, domain.CollectionChat)
	tr := repository.NewTeamRepository(db, domain.CollectionTeam)
	ch := &wsHandler.ChatWSHandler{
		ChatUsecase: usecase.NewChatUsecase(cr, tr, timeout, env),
	}

	group.Use("/", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)

			authHeader := c.Queries()["token"]
			userData, err := tokenutil.ExtractDataFromToken(authHeader, env.AccessTokenSecret)
			if err != nil {
				return internal.SendErrorRespond(c, 401, map[string][]string{"error": {err.Error()}})
			}

			c.Locals("user", userData)
			return c.Next()
		}

		return fiber.ErrUpgradeRequired
	})

	group.Get("/:teamId", websocket.New(ch.ChatOnTeam))
}
