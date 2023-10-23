package routes

import (
	"sportsync/bootstrap"
	httpHandler "sportsync/delivery/http"
	"sportsync/domain"
	"sportsync/mongo"
	"sportsync/repository"
	"sportsync/usecase"
	"time"

	"github.com/gofiber/fiber/v2"
)

func NewCompetitionRoute(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group fiber.Router) {
	cr := repository.NewCompetitionRepository(db, domain.CollectionCup, domain.CollectionLeague)
	cu := usecase.NewCompetitionUsecase(cr, timeout, env)
	ch := &httpHandler.CompetitionHandler{
		CompetitionUsecase: cu, Env: env,
	}

	cupGroup := group.Group("cup")
	cupGroup.Post("create", ch.CreateCup)
}
