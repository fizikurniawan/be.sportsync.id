package repository

import (
	"context"
	"math"

	"sportsync/domain"
	"sportsync/entities"
	"sportsync/models"
	"sportsync/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type teamRepository struct {
	database   mongo.Database
	collection string
}

func NewTeamRepository(db mongo.Database, collection string) domain.TeamRepository {
	return &teamRepository{
		database:   db,
		collection: collection,
	}
}

func (tr *teamRepository) Create(c context.Context, team *entities.Team) error {
	collection := tr.database.Collection(tr.collection)

	_, err := collection.InsertOne(c, team)

	return err
}

func (tr *teamRepository) GetByNameAndSport(c context.Context, name string, sportName string) (team entities.Team, err error) {
	collection := tr.database.Collection(tr.collection)
	err = collection.FindOne(c, bson.M{"name": name, "sport_name": sportName}).Decode(&team)
	return
}

func (tr *teamRepository) GetMyTeam(c context.Context, filter models.GetMyTeamBody, userId string) (teams []entities.Team, page models.Page, err error) {
	var totalDataCount interface{}
	collection := tr.database.Collection(tr.collection)

	filterMongo := bson.M{
		"$and": []bson.M{
			{"name": bson.M{"$regex": filter.Search, "$options": "i"}}, // Case-insensitive search on the "name" field
			{"owned_by_id": userId}, // Match documents with the specified owned_by_id
		},
	}

	// Define options for pagination
	options := options.Find().SetSkip(int64((filter.Page - 1) * filter.Size)).SetLimit(int64(filter.Size))

	cursor, err := collection.Find(c, filterMongo, options)
	if err != nil {
		return
	}
	totalDataCount, err = collection.CountDocuments(c, filterMongo)

	if err != nil {
		return
	}

	totalDataCountInt := totalDataCount.(int64)
	totalPages := int(math.Ceil(float64(totalDataCountInt) / float64(filter.Size)))

	page.Total = int(totalDataCountInt)
	page.Size = filter.Size
	page.Current = filter.Page
	page.TotalPages = totalPages

	err = cursor.All(c, &teams)
	if teams == nil {
		return
	}

	return
}
