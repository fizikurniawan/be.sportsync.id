package domain

import (
	"context"
	"sportsync/models"
)

type AuthUsecase interface {
	Register(ctx context.Context, userReq models.RegisterBody) error
	Login(ctx context.Context, userReq models.LoginBody) (user models.LoginResponse, err error)
}
