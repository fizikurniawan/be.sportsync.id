package domain

import (
	"context"
	"sportsync/entities"
	"sportsync/models"
)

type AuthUsecase interface {
	Register(ctx context.Context, userReq models.RegisterBody) error
	Login(ctx context.Context, userReq models.LoginBody) (user models.LoginResponse, err error)
	CreateAccessToken(user *entities.User, secret string, expiry int) (accessToken string, err error)
	CreateRefreshToken(user *entities.User, secret string, expiry int) (refreshToken string, err error)
	RefreshToken(ctx context.Context, refreshToken string) (res models.LoginResponse, err error)
}
