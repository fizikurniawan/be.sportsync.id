package usecase

import (
	"context"
	"sportsync/domain"
	"sportsync/entities"
	"sportsync/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewAuthUsecase(userRepository domain.UserRepository, timeout time.Duration) domain.AuthUsecase {
	return &authUsecase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (uu *authUsecase) Register(ctx context.Context, userReq models.RegisterRequest) (err error) {
	var user entities.User
	user.Email = userReq.Email
	user.Name = userReq.Name

	var hashedPassword []byte
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	user.Password = string(hashedPassword)
	err = uu.userRepository.Create(ctx, &user)

	return
}
