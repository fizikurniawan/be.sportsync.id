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
	"strings"
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

	user, err = au.userRepository.GetByEmail(ctx, userReq.Email)
	if err != nil && !strings.Contains(err.Error(), "no documents in result") {
		return
	}

	if user.Email != "" {
		err = errors.New("email has already exists")
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	user.Email = userReq.Email
	user.Name = userReq.Name
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

	accessToken, _ = au.CreateAccessToken(&u, au.env.AccessTokenSecret, au.env.AccessTokenExpiryHour)
	refreshToken, _ = au.CreateRefreshToken(&u, au.env.RefreshTokenSecret, au.env.RefreshTokenExpiryHour)

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

func (au *authUsecase) RefreshToken(ctx context.Context, refreshToken string) (res models.LoginResponse, err error) {
	id, _ := tokenutil.ExtractIDFromToken(refreshToken, au.env.RefreshTokenSecret)

	if id == "" {
		err = errors.New("user not found")
		return
	}

	user, err := au.userRepository.GetByID(ctx, id)
	if user.Email == "" {
		err = errors.New("user not found")
		return
	}

	accessToken, _ := tokenutil.CreateAccessToken(&user, au.env.AccessTokenSecret, au.env.AccessTokenExpiryHour)
	newRefreshToken, _ := tokenutil.CreateRefreshToken(&user, au.env.RefreshTokenSecret, au.env.RefreshTokenExpiryHour)

	res.Email = user.Email
	res.Name = user.Name
	res.RefreshToken = newRefreshToken
	res.AccessToken = accessToken

	return
}
