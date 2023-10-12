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

func NewAuthRoute(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group fiber.Router) {
	wr := repository.NewUserRepository(db, domain.CollectionUser)
	wc := &httpHandler.AuthHandler{
		AuthUsecase: usecase.NewAuthUsecase(wr, timeout, env),
	}

	group.Post("register", wc.Register)
	group.Post("login", wc.Login)
	group.Post("refresh", wc.RefreshToken)
}
