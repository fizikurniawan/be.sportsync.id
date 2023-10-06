package routes

import (
	"sportsync/bootstrap"
	httpHandler "sportsync/delivery/http"
	"sportsync/mongo"
	"sportsync/repository"
	"sportsync/usecase"
	"time"

	"github.com/gofiber/fiber/v2"
)

func NewTeamRoute(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group fiber.Router) {
	tr := repository.NewTeamRepository(db, "teams")
	th := &httpHandler.TeamHandler{
		TeamUsecase: usecase.NewTeamUsecase(tr, timeout, env),
	}

	group.Post("create", th.Create)
	group.Post("my-team", th.MyTeam)
}
