package usecase

import (
	"context"
	"errors"
	"sportsync/bootstrap"
	"sportsync/domain"
	"sportsync/entities"
	"sportsync/models"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type teamUsecase struct {
	teamRepository domain.TeamRepository
	contextTimeout time.Duration
	env            *bootstrap.Env
}

func NewTeamUsecase(teamRepository domain.TeamRepository, timeout time.Duration, env *bootstrap.Env) domain.TeamUsecase {
	return &teamUsecase{
		teamRepository: teamRepository,
		contextTimeout: timeout,
		env:            env,
	}
}

func (tu *teamUsecase) Create(ctx context.Context, team models.TeamBody, userId string) (err error) {
	// checl if name and sport already exists
	var teamExists entities.Team
	teamExists, err = tu.teamRepository.GetByNameAndSport(ctx, team.Name, team.SportName)

	if err != nil && err != mongo.ErrNoDocuments {
		return
	} else {
		err = nil
	}
	if teamExists.Name != "" {
		err = errors.New("name and sport has already exists")
		return
	}
	teamInsert := entities.Team{
		Name:      team.Name,
		SportName: team.SportName,
		Logo:      team.Logo,
		OwnedByID: userId,
	}
	tu.teamRepository.Create(ctx, &teamInsert)

	return
}

func (tu *teamUsecase) GetMyTeam(ctx context.Context, filter models.GetMyTeamBody, userId string) (teams []entities.Team, page models.Page, err error) {
	teams, page, err = tu.teamRepository.GetMyTeam(ctx, filter, userId)
	return
}
