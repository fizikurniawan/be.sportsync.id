package domain

import (
	"context"
	"sportsync/models"
)

type AuthUsecase interface {
	Register(ctx context.Context, userReq models.RegisterRequest) error
}
