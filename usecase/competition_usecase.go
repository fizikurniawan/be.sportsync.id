package usecase

import (
	"context"
	"sportsync/bootstrap"
	"sportsync/domain"
	"sportsync/entities"
	"sportsync/models"
	"time"
)

type competitionUsecase struct {
	competitionRepository domain.CompetitionRepository
	contextTimeout        time.Duration
	env                   *bootstrap.Env
}

func NewCompetitionUsecase(competitionRepository domain.CompetitionRepository, timeout time.Duration, env *bootstrap.Env) domain.CompetitionUsecase {
	return &competitionUsecase{
		competitionRepository: competitionRepository,
		contextTimeout:        timeout,
		env:                   env,
	}
}

func (cu *competitionUsecase) CreateCup(c context.Context, cup *models.CupBody) error {
	inputStartDate := cup.StartDate
	startDate, err := time.Parse("2006-01-02", inputStartDate)
	if err != nil {
		return err
	}

	cupInsert := entities.Cup{
		Name:           cup.Name,
		MaxParticipant: cup.MaxParticipant,
		Format:         cup.Format,
		StartDate:      startDate,
		Season:         cup.Season,
	}
	if err := cu.competitionRepository.CreateCup(c, &cupInsert); err != nil {
		return err
	}

	return nil
}

func (cu *competitionUsecase) CreateLeague(c context.Context, league *entities.League) error {
	if err := cu.competitionRepository.CreateLeague(c, league); err != nil {
		return err
	}
	return nil
}

func (cu *competitionUsecase) GetCupByID(c context.Context, cupId string) (cup entities.Cup, err error) {
	cup, err = cu.competitionRepository.GetCupByID(c, cupId)
	return
}

func (cu *competitionUsecase) GetLeagueByID(c context.Context, leagueId string) (league entities.League, err error) {
	league, err = cu.competitionRepository.GetLeagueByID(c, leagueId)
	return
}
