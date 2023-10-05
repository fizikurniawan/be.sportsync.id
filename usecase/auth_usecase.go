package usecase

import (
	"context"
	"errors"
	"fmt"
	"sportsync/bootstrap"
	"sportsync/domain"
	"sportsync/entities"
	"sportsync/internal/tokenutil"
	"sportsync/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
	env            *bootstrap.Env
}

func NewAuthUsecase(userRepository domain.UserRepository, timeout time.Duration, env *bootstrap.Env) domain.AuthUsecase {
	return &authUsecase{
		userRepository: userRepository,
		contextTimeout: timeout,
		env:            env,
	}
}

func (au *authUsecase) Register(ctx context.Context, userReq models.RegisterBody) (err error) {
	var user entities.User
	user.Email = userReq.Email
	user.Name = userReq.Name

	user, err = au.userRepository.GetByEmail(ctx, userReq.Email)
	if err != nil {
		return
	}

	if user.Email != "" {
		err = errors.New("email has already exists")
		return
	}

	var hashedPassword []byte
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	user.Password = string(hashedPassword)
	err = au.userRepository.Create(ctx, &user)

	return
}

func (au *authUsecase) Login(ctx context.Context, userReq models.LoginBody) (user models.LoginResponse, err error) {
	defaultErr := fmt.Errorf("invalid credentials")

	var u entities.User
	u, err = au.userRepository.GetByEmail(ctx, userReq.Email)
	if err != nil {
		return
	}

	var accessToken string
	var refreshToken string

	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(userReq.Password)) != nil {
		err = defaultErr
		return
	}

	accessToken, err = au.CreateAccessToken(&u, au.env.AccessTokenSecret, au.env.AccessTokenExpiryHour)
	if err != nil {
		return
	}

	refreshToken, err = au.CreateRefreshToken(&u, au.env.RefreshTokenSecret, au.env.RefreshTokenExpiryHour)
	if err != nil {
		return
	}

	user.AccessToken = accessToken
	user.RefreshToken = refreshToken
	user.Email = u.Email
	user.Name = u.Name

	return
}

func (au *authUsecase) CreateAccessToken(user *entities.User, secret string, expiry int) (accessToken string, err error) {
	return tokenutil.CreateAccessToken(user, secret, expiry)
}

func (au *authUsecase) CreateRefreshToken(user *entities.User, secret string, expiry int) (refreshToken string, err error) {
	return tokenutil.CreateRefreshToken(user, secret, expiry)
}
