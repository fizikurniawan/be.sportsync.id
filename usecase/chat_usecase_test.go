package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"sportsync/bootstrap"
	"sportsync/entities"
	"sportsync/mocks"
	"sportsync/usecase"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestChatOnTeam(t *testing.T) {
	env := bootstrap.Env{}

	t.Run("Success", func(t *testing.T) {
		teamId := primitive.NewObjectID().Hex()
		mockChatRepo := new(mocks.ChatRepository)
		mockTeamRepo := new(mocks.TeamRepository)
		cu := usecase.NewChatUsecase(mockChatRepo, mockTeamRepo, time.Second*2, &env)
		err := cu.ChatOnTeam(string(teamId))
		assert.NoError(t, err)
	})
}

func TestGetRecentMessages(t *testing.T) {

	mockChatRepo := new(mocks.ChatRepository)
	mockTeamRepo := new(mocks.TeamRepository)
	teamId := string(primitive.NewObjectID().Hex())
	env := bootstrap.Env{}
	t.Run("Success", func(t *testing.T) {
		cu := usecase.NewChatUsecase(mockChatRepo, mockTeamRepo, time.Second*2, &env)
		totalData := 3
		var mockChats []entities.Chat
		for i := 0; i < totalData; i++ {
			mockChat := entities.Chat{ID: string(primitive.NewObjectID().Hex()), Message: fmt.Sprintf("Test Message %d", i), TeamID: teamId}
			mockChats = append(mockChats, mockChat)
		}

		mockChatRepo.On("GetRecentMessages", mock.Anything, teamId, "", "").Return(mockChats, nil).Once()
		res, err := cu.GetRecentMessages(context.Background(), teamId, "", "")
		assert.NoError(t, err)
		assert.Equal(t, len(res), totalData)
		mockChatRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		cu := usecase.NewChatUsecase(mockChatRepo, mockTeamRepo, time.Second*2, &env)
		mockChatRepo.On("GetRecentMessages", mock.Anything, teamId, "", "").Return([]entities.Chat{}, errors.New("error")).Once()
		_, err := cu.GetRecentMessages(context.Background(), teamId, "", "")
		assert.Error(t, err)
	})

}

func TestInsert(t *testing.T) {
	env := bootstrap.Env{}
	teamId := string(primitive.NewObjectID().Hex())

	mockChatRepo := new(mocks.ChatRepository)
	mockTeamRepo := new(mocks.TeamRepository)
	t.Run("Success", func(t *testing.T) {
		mockChatRepo := new(mocks.ChatRepository)
		mockTeamRepo := new(mocks.TeamRepository)
		mockChat := entities.Chat{
			Message: "Test Message",
			TeamID:  teamId,
		}

		mockChatRepo.On("Insert", mock.Anything, mockChat).Return(nil).Once()
		cu := usecase.NewChatUsecase(mockChatRepo, mockTeamRepo, time.Second*2, &env)
		err := cu.Insert(context.Background(), mockChat)
		assert.NoError(t, err)
	})
	t.Run("Error", func(t *testing.T) {
		mockChat := entities.Chat{
			Message: "Test Message",
			TeamID:  teamId,
		}

		mockChatRepo.On("Insert", mock.Anything, mockChat).Return(errors.New("error")).Once()
		cu := usecase.NewChatUsecase(mockChatRepo, mockTeamRepo, time.Second*2, &env)
		err := cu.Insert(context.Background(), mockChat)
		assert.Error(t, err)
	})
}

func TestIsTeamExists(t *testing.T) {
	env := bootstrap.Env{}
	teamId := string(primitive.NewObjectID().Hex())

	mockChatRepo := new(mocks.ChatRepository)
	mockTeamRepo := new(mocks.TeamRepository)

	t.Run("Success", func(t *testing.T) {
		mockTeam := entities.Team{
			ID:   teamId,
			Name: "Test team",
		}
		mockTeamRepo.On("GetByID", mock.Anything, teamId).Return(mockTeam, nil).Once()
		au := usecase.NewChatUsecase(mockChatRepo, mockTeamRepo, time.Second*2, &env)
		isExists, err := au.IsTeamExists(context.Background(), teamId)
		assert.NoError(t, err)
		assert.Equal(t, isExists, true)
		mockTeamRepo.AssertExpectations(t)
	})
	t.Run("Error", func(t *testing.T) {
		mockTeam := entities.Team{}
		mockTeamRepo.On("GetByID", mock.Anything, teamId).Return(mockTeam, errors.New("error")).Once()
		au := usecase.NewChatUsecase(mockChatRepo, mockTeamRepo, time.Second*2, &env)
		isExists, err := au.IsTeamExists(context.Background(), teamId)
		assert.Error(t, err)
		assert.Equal(t, isExists, false)
		mockTeamRepo.AssertExpectations(t)
	})
	t.Run("Error Data", func(t *testing.T) {
		mockTeam := entities.Team{}
		mockTeamRepo.On("GetByID", mock.Anything, teamId).Return(mockTeam, nil).Once()
		au := usecase.NewChatUsecase(mockChatRepo, mockTeamRepo, time.Second*2, &env)
		isExists, err := au.IsTeamExists(context.Background(), teamId)
		assert.NoError(t, err)
		assert.Equal(t, isExists, false)
		mockTeamRepo.AssertExpectations(t)
	})
}
