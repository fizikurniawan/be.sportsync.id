package repository

import (
	"context"
	"sportsync/domain"
	"sportsync/entities"
	"sportsync/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

type competitionRepository struct {
	database         mongo.Database
	cupCollection    string
	leagueCollection string
}

func NewCompetitionRepository(database mongo.Database, cupCollection string, leagueCollection string) domain.CompetiotionRepository {
	return &competitionRepository{
		database:         database,
		cupCollection:    cupCollection,
		leagueCollection: leagueCollection,
	}
}

func (cr *competitionRepository) CreateCup(c context.Context, cup *entities.Cup) error {
	collection := cr.database.Collection(cr.cupCollection)

	_, err := collection.InsertOne(c, cup)

	return err
}

func (cr *competitionRepository) CreateLeague(c context.Context, league *entities.League) error {
	collection := cr.database.Collection(cr.leagueCollection)

	_, err := collection.InsertOne(c, league)

	return err
}

func (cr *competitionRepository) GetCupByID(c context.Context, cupId string) (cup entities.Cup, err error) {
	collection := cr.database.Collection(cr.cupCollection)
	err = collection.FindOne(c, bson.M{"_id": cupId}).Decode(&cup)
	return
}

func (cr *competitionRepository) GetLeagueByID(c context.Context, cupId string) (cup entities.Cup, err error) {
	collection := cr.database.Collection(cr.cupCollection)
	err = collection.FindOne(c, bson.M{"_id": cupId}).Decode(&cup)
	return
}
