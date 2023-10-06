package domain

import (
	"context"
	"sportsync/entities"
	"sportsync/models"
)

type TeamRepository interface {
	Create(c context.Context, team *entities.Team) error
	GetByNameAndSport(c context.Context, name string, sportName string) (team entities.Team, err error)
	GetMyTeam(c context.Context, filter models.GetMyTeamBody, userId string) (teams []entities.Team, page models.Page, err error)
	GetByID(c context.Context, id string) (entities.Team, error)
}

type TeamUsecase interface {
	Create(ctx context.Context, team models.TeamBody, userId string) error
	GetMyTeam(ctx context.Context, filter models.GetMyTeamBody, userId string) (teams []entities.Team, page models.Page, err error)
}
