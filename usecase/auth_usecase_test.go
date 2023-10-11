package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"sportsync/bootstrap"
	"sportsync/entities"
	"sportsync/internal/tokenutil"
	"sportsync/mocks"
	"sportsync/models"
	"sportsync/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	mockUserRepository := new(mocks.UserRepository)
	objectID := primitive.NewObjectID()
	id := objectID.Hex()
	var hashedPassword []byte
	hashedPassword, _ = bcrypt.GenerateFromPassword([]byte("Adm!n1234"), bcrypt.DefaultCost)
	env := bootstrap.Env{}
	mockUser := entities.User{
		ID:       id,
		Name:     "Test Name",
		Email:    "test@mail.com",
		Password: string(hashedPassword),
	}
	reqBody := models.RegisterBody{
		Name:     "Test Name",
		Email:    "test@mail.com",
		Password: "Adm!n1234",
	}

	t.Run("Success", func(t *testing.T) {
		var mockUserGet entities.User

		mockUserRepository.On("GetByEmail", mock.Anything, mockUser.Email).Return(mockUserGet, nil).Once()

		// bcrypt always create diff hashes, so compare password
		userPassMatcher := mock.MatchedBy(func(user *entities.User) bool {
			return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password)) == nil
		})
		mockUserRepository.On("Create", mock.Anything, userPassMatcher).Return(nil).Once()

		u := usecase.NewAuthUsecase(mockUserRepository, time.Second*2, &env)
		err := u.Register(context.Background(), reqBody)

		assert.NoError(t, err)
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error exist email", func(t *testing.T) {
		mockUserRepository.On("GetByEmail", mock.Anything, mockUser.Email).Return(mockUser, nil).Once()

		u := usecase.NewAuthUsecase(mockUserRepository, time.Second*2, &env)
		err := u.Register(context.Background(), reqBody)
		assert.Error(t, err)

		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error on GetByEmail", func(t *testing.T) {
		// Expect GetByEmail to return an error
		mockUserRepository.On("GetByEmail", mock.Anything, mockUser.Email).Return(entities.User{}, errors.New("error get data")).Once()

		u := usecase.NewAuthUsecase(mockUserRepository, time.Second*2, &env)
		err := u.Register(context.Background(), reqBody)
		assert.Error(t, err)

		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error on Create", func(t *testing.T) {
		var mockUserGet entities.User

		// Expect GetByEmail to return empty user (user not found)
		mockUserRepository.On("GetByEmail", mock.Anything, mockUser.Email).Return(mockUserGet, nil).Once()

		// bcrypt always creates different hashes, so compare the password using MatchedBy
		userPassMatcher := mock.MatchedBy(func(user *entities.User) bool {
			return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password)) == nil
		})

		// Expect Create to be called with the matched password and return an error
		mockUserRepository.On("Create", mock.Anything, userPassMatcher).Return(errors.New("error")).Once()

		u := usecase.NewAuthUsecase(mockUserRepository, time.Second*2, &env)
		err := u.Register(context.Background(), reqBody)
		assert.Error(t, err)

		mockUserRepository.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	mockUserRepository := new(mocks.UserRepository)
	objectID := primitive.NewObjectID()
	id := objectID.Hex()
	var hashedPassword []byte
	hashedPassword, _ = bcrypt.GenerateFromPassword([]byte("Adm!n1234"), bcrypt.DefaultCost)
	env := bootstrap.Env{}
	mockUser := entities.User{
		ID:       id,
		Name:     "Test Name",
		Email:    "test@mail.com",
		Password: string(hashedPassword),
	}
	reqBody := models.LoginBody{
		Email:    "test@mail.com",
		Password: "Adm!n1234",
	}

	t.Run("Success", func(t *testing.T) {
		mockUserRepository.On("GetByEmail", mock.Anything, mockUser.Email).Return(mockUser, nil).Once()

		u := usecase.NewAuthUsecase(mockUserRepository, time.Second*2, &env)
		user, err := u.Login(context.Background(), reqBody)
		assert.NoError(t, err)
		assert.Equal(t, user.Name, mockUser.Name)
		assert.Equal(t, user.Email, mockUser.Email)
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error on GetByEmail", func(t *testing.T) {
		// Expect GetByEmail to return an error
		mockUserRepository.On("GetByEmail", mock.Anything, mockUser.Email).Return(entities.User{}, errors.New("error get data")).Once()

		u := usecase.NewAuthUsecase(mockUserRepository, time.Second*2, &env)
		_, err := u.Login(context.Background(), reqBody)
		assert.Error(t, err)

		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error on Difference Password", func(t *testing.T) {
		// Expect GetByEmail to return an error
		mockUser.Password = "invalidpassword"
		mockUserRepository.On("GetByEmail", mock.Anything, mockUser.Email).Return(mockUser, nil).Once()

		u := usecase.NewAuthUsecase(mockUserRepository, time.Second*2, &env)
		_, err := u.Login(context.Background(), reqBody)
		assert.Error(t, err)

		mockUserRepository.AssertExpectations(t)
	})
}

func TestRefreshToken(t *testing.T) {
	mockUserRepository := new(mocks.UserRepository)
	userId := primitive.NewObjectID().Hex()
	userMock := entities.User{
		ID:    userId,
		Name:  "Test Name",
		Email: "Test Email",
	}

	env := bootstrap.Env{
		RefreshTokenExpiryHour: 2,
		AccessTokenExpiryHour:  2,
		AccessTokenSecret:      "secret",
		RefreshTokenSecret:     "secret",
	}

	t.Run("Error user not found", func(t *testing.T) {

		newUserMock := new(entities.User)
		newUserMock.ID = userId
		newUserMock.Email = "Test Email"
		token, _ := tokenutil.CreateRefreshToken(newUserMock, env.RefreshTokenSecret, env.AccessTokenExpiryHour)

		mockUserRepository.On("GetByID", mock.Anything, userId).Return(entities.User{Email: ""}, nil).Once()

		au := usecase.NewAuthUsecase(mockUserRepository, time.Second*2, &env)
		_, err := au.RefreshToken(context.Background(), token)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "user not found")
	})

	t.Run("Success", func(t *testing.T) {
		token, _ := tokenutil.CreateRefreshToken(&userMock, env.RefreshTokenSecret, env.AccessTokenExpiryHour)

		mockUserRepository.On("GetByID", mock.Anything, userId).Return(userMock, nil).Once()

		au := usecase.NewAuthUsecase(mockUserRepository, time.Second*2, &env)
		res, err := au.RefreshToken(context.Background(), token)
		assert.NoError(t, err)

		assert.Equal(t, userMock.Name, res.Name)
		assert.Equal(t, userMock.Email, res.Email)
		mockUserRepository.AssertExpectations(t)
	})
	t.Run("Error invalid token", func(t *testing.T) {

		mockUserRepository.On("GetByID", mock.Anything, userId).Return(userMock, nil).Once()

		au := usecase.NewAuthUsecase(mockUserRepository, time.Second*2, &env)
		_, err := au.RefreshToken(context.Background(), "invalid")
		assert.Error(t, err)

	})
	t.Run("Error invalid userId", func(t *testing.T) {

		mockUserRepository.On("GetByID", mock.Anything, userId).Return(entities.User{}, nil).Once()

		au := usecase.NewAuthUsecase(mockUserRepository, time.Second*2, &env)
		_, err := au.RefreshToken(context.Background(), "invalid")

		assert.Error(t, err)
	})

}
