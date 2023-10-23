package domain

import (
	"context"
	"sportsync/entities"
	"sportsync/models"
)

const (
	CollectionCup    = "cups"
	CollectionLeague = "leagues"
)

type CompetitionRepository interface {
	CreateCup(c context.Context, cup *entities.Cup) error
	CreateLeague(c context.Context, league *entities.League) error
	GetCupByID(c context.Context, cupId string) (cup entities.Cup, err error)
	GetLeagueByID(c context.Context, cupId string) (league entities.League, err error)
}

type CompetitionUsecase interface {
	CreateCup(c context.Context, cup *models.CupBody) error
	CreateLeague(c context.Context, league *entities.League) error
	GetCupByID(c context.Context, cupId string) (cup entities.Cup, err error)
	GetLeagueByID(c context.Context, cupId string) (league entities.League, err error)
}
