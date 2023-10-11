package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sportsync/bootstrap"
	"sportsync/entities"
	"sportsync/mocks"
	"sportsync/models"
	"sportsync/usecase"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestCreate(t *testing.T) {
	mockTeamRepo := new(mocks.TeamRepository)
	env := bootstrap.Env{}

	t.Run("Success", func(t *testing.T) {
		mockTeamExists := entities.Team{}
		mockTeamBody := models.TeamBody{
			Name:      "Exists Team",
			SportName: "Exists Sport",
		}
		userId := string(primitive.NewObjectID().Hex())
		mockTeam := entities.Team{
			Name:      mockTeamBody.Name,
			SportName: mockTeamBody.SportName,
			OwnedByID: userId,
		}
		mockTeamRepo.On("GetByNameAndSport", mock.Anything, mockTeamBody.Name, mockTeamBody.SportName).Return(mockTeamExists, nil).Once()
		mockTeamRepo.On("Create", mock.Anything, &mockTeam).Return(nil).Once()

		tu := usecase.NewTeamUsecase(mockTeamRepo, time.Second*2, &env)
		err := tu.Create(context.Background(), mockTeamBody, userId)
		assert.NoError(t, err)
	})
	t.Run("Error data exists", func(t *testing.T) {
		mockTeamExists := entities.Team{
			Name:      "Exists Team",
			SportName: "Exists Sport",
		}
		mockTeamBody := models.TeamBody{
			Name:      "Exists Team",
			SportName: "Exists Sport",
		}
		userId := string(primitive.NewObjectID().Hex())
		mockTeamRepo.On("GetByNameAndSport", mock.Anything, mockTeamBody.Name, mockTeamBody.SportName).Return(mockTeamExists, nil).Once()

		tu := usecase.NewTeamUsecase(mockTeamRepo, time.Second*2, &env)
		err := tu.Create(context.Background(), mockTeamBody, userId)
		assert.Error(t, err)
	})
	t.Run("Error creating", func(t *testing.T) {
		mockTeamExists := entities.Team{}
		mockTeamBody := models.TeamBody{
			Name:      "Exists Team",
			SportName: "Exists Sport",
		}
		userId := string(primitive.NewObjectID().Hex())
		mockTeam := entities.Team{
			Name:      mockTeamBody.Name,
			SportName: mockTeamBody.SportName,
			OwnedByID: userId,
		}
		mockTeamRepo.On("GetByNameAndSport", mock.Anything, mockTeamBody.Name, mockTeamBody.SportName).Return(mockTeamExists, nil).Once()
		mockTeamRepo.On("Create", mock.Anything, &mockTeam).Return(errors.New("error")).Once()

		tu := usecase.NewTeamUsecase(mockTeamRepo, time.Second*2, &env)
		err := tu.Create(context.Background(), mockTeamBody, userId)
		assert.Error(t, err)
	})
	t.Run("Error Document Team Not Exists", func(t *testing.T) {
		mockTeamExists := entities.Team{}
		mockTeamBody := models.TeamBody{
			Name:      "Exists Team",
			SportName: "Exists Sport",
		}
		userId := string(primitive.NewObjectID().Hex())
		mockTeam := entities.Team{
			Name:      mockTeamBody.Name,
			SportName: mockTeamBody.SportName,
			OwnedByID: userId,
		}
		mockTeamRepo.On("GetByNameAndSport", mock.Anything, mockTeamBody.Name, mockTeamBody.SportName).Return(mockTeamExists, mongo.ErrNoDocuments).Once()
		mockTeamRepo.On("Create", mock.Anything, &mockTeam).Return(errors.New("error")).Once()

		tu := usecase.NewTeamUsecase(mockTeamRepo, time.Second*2, &env)
		err := tu.Create(context.Background(), mockTeamBody, userId)
		assert.Error(t, err)
	})
	t.Run("Error Non Document Team Not Exists", func(t *testing.T) {
		mockTeamExists := entities.Team{}
		mockTeamBody := models.TeamBody{
			Name:      "Exists Team",
			SportName: "Exists Sport",
		}
		userId := string(primitive.NewObjectID().Hex())
		mockTeam := entities.Team{
			Name:      mockTeamBody.Name,
			SportName: mockTeamBody.SportName,
			OwnedByID: userId,
		}
		mockTeamRepo.On("GetByNameAndSport", mock.Anything, mockTeamBody.Name, mockTeamBody.SportName).Return(mockTeamExists, errors.New("error mongo")).Once()
		mockTeamRepo.On("Create", mock.Anything, &mockTeam).Return(errors.New("error")).Once()

		tu := usecase.NewTeamUsecase(mockTeamRepo, time.Second*2, &env)
		err := tu.Create(context.Background(), mockTeamBody, userId)
		assert.Error(t, err)
	})
}

func TestGetMyTeam(t *testing.T) {
	mockTeamRepo := new(mocks.TeamRepository)
	env := bootstrap.Env{}
	mockFilter := models.GetMyTeamBody{
		Search: "",
		Size:   10,
		Page:   1,
	}
	userId := string(primitive.NewObjectID().Hex())
	teamId := string(primitive.NewObjectID().Hex())

	t.Run("Success", func(t *testing.T) {
		var mockTeams []entities.Team
		totalData := 15
		for i := 0; i < totalData; i++ {
			mockTeam := entities.Team{
				ID:        teamId,
				Name:      fmt.Sprintf("Team %d", i),
				SportName: "volley",
				OwnedByID: userId,
			}
			mockTeams = append(mockTeams, mockTeam)
		}
		mockPage := models.Page{
			Total:      totalData,
			Size:       mockFilter.Size,
			Current:    mockFilter.Page,
			TotalPages: int(math.Ceil(float64(totalData) / float64(mockFilter.Size))),
		}

		mockTeamRepo.On("GetMyTeam", mock.Anything, mockFilter, userId).Return(mockTeams, mockPage, nil).Once()
		tu := usecase.NewTeamUsecase(mockTeamRepo, time.Second*2, &env)
		teams, page, err := tu.GetMyTeam(context.Background(), mockFilter, userId)
		assert.NoError(t, err)
		assert.Equal(t, teams, mockTeams)
		assert.Equal(t, page, mockPage)
		for _, team := range teams {
			assert.Equal(t, team.OwnedByID, userId)
		}
	})
	t.Run("Error", func(t *testing.T) {
		var mockTeams []entities.Team

		mockTeamRepo.On("GetMyTeam", mock.Anything, mockFilter, userId).Return(mockTeams, models.Page{}, errors.New("error")).Once()
		tu := usecase.NewTeamUsecase(mockTeamRepo, time.Second*2, &env)
		_, _, err := tu.GetMyTeam(context.Background(), mockFilter, userId)
		assert.Error(t, err)

	})
}
