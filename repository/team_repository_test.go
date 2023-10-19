package repository_test

import (
	"context"
	"errors"
	"sportsync/domain"
	"sportsync/entities"
	"sportsync/models"
	"sportsync/mongo/mocks"
	"sportsync/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCreateTeam(t *testing.T) {
	var databaseHelper *mocks.Database
	var collectionHelper *mocks.Collection

	databaseHelper = &mocks.Database{}
	collectionHelper = &mocks.Collection{}

	collectionName := domain.CollectionTeam

	mockUserID := primitive.NewObjectID()
	mockTeam := &entities.Team{
		ID:        primitive.NewObjectID().Hex(),
		Name:      "Test",
		SportName: "football",
		OwnedByID: mockUserID.Hex(),
	}

	mockEmptyTeam := &entities.Team{}

	t.Run("success", func(t *testing.T) {

		collectionHelper.On("InsertOne", mock.Anything, mock.AnythingOfType("*entities.Team")).Return(mockTeam.ID, nil).Once()

		databaseHelper.On("Collection", collectionName).Return(collectionHelper)

		tr := repository.NewTeamRepository(databaseHelper, collectionName)

		err := tr.Create(context.Background(), mockTeam)

		assert.NoError(t, err)

		collectionHelper.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		collectionHelper.On("InsertOne", mock.Anything, mock.AnythingOfType("*entities.Team")).Return(mockEmptyTeam, errors.New("Unexpected")).Once()

		databaseHelper.On("Collection", collectionName).Return(collectionHelper)

		tr := repository.NewTeamRepository(databaseHelper, collectionName)

		err := tr.Create(context.Background(), mockEmptyTeam)

		assert.Error(t, err)

		collectionHelper.AssertExpectations(t)
	})
}

func TestGetByNameAndSport(t *testing.T) {
	teamId := string(primitive.NewObjectID().Hex())
	userId := string(primitive.NewObjectID().Hex())

	var databaseHelper *mocks.Database
	var collectionHelper *mocks.Collection
	var singleResultHelper *mocks.SingleResult

	databaseHelper = &mocks.Database{}
	collectionHelper = &mocks.Collection{}
	singleResultHelper = &mocks.SingleResult{}

	collectionName := domain.CollectionTeam
	t.Run("Success", func(t *testing.T) {
		// err = collection.FindOne(c, bson.M{"name": name, "sport_name": sportName}).Decode(&team)
		expected := &entities.Team{Name: "Test Team", ID: teamId, OwnedByID: userId, SportName: "football"}

		collectionHelper.On("FindOne", mock.Anything, mock.Anything).Return(singleResultHelper).Once()
		singleResultHelper.On("Decode", &entities.Team{}).Run(func(args mock.Arguments) {
			expectedPnt := args.Get(0).(*entities.Team)
			*expectedPnt = *expected
		}).Return(nil).Once()

		databaseHelper.On("Collection", collectionName).Return(collectionHelper).Once()

		tu := repository.NewTeamRepository(databaseHelper, collectionName)
		team, err := tu.GetByNameAndSport(context.Background(), expected.Name, expected.SportName)
		assert.NoError(t, err)
		assert.NotNil(t, team)
		assert.Equal(t, team.Name, expected.Name)

	})
}
func TestGetByIDTeam(t *testing.T) {
	teamId := string(primitive.NewObjectID().Hex())
	userId := string(primitive.NewObjectID().Hex())

	var databaseHelper *mocks.Database
	var collectionHelper *mocks.Collection
	var singleResultHelper *mocks.SingleResult

	databaseHelper = &mocks.Database{}
	collectionHelper = &mocks.Collection{}
	singleResultHelper = &mocks.SingleResult{}

	collectionName := domain.CollectionTeam
	t.Run("Success", func(t *testing.T) {
		// err = collection.FindOne(c, bson.M{"name": name, "sport_name": sportName}).Decode(&team)
		expected := &entities.Team{Name: "Test Team", ID: teamId, OwnedByID: userId, SportName: "football"}

		collectionHelper.On("FindOne", mock.Anything, mock.Anything).Return(singleResultHelper).Once()
		singleResultHelper.On("Decode", &entities.Team{}).Run(func(args mock.Arguments) {
			expectedPnt := args.Get(0).(*entities.Team)
			*expectedPnt = *expected
		}).Return(nil).Once()

		databaseHelper.On("Collection", collectionName).Return(collectionHelper).Once()

		tu := repository.NewTeamRepository(databaseHelper, collectionName)
		team, err := tu.GetByID(context.Background(), expected.ID)
		assert.NoError(t, err)
		assert.NotNil(t, team)
		assert.Equal(t, team.Name, expected.Name)

	})
}

func TestGetMyTeam(t *testing.T) {
	userId := string(primitive.NewObjectID().Hex())

	var databaseHelper *mocks.Database
	var collectionHelper *mocks.Collection
	var cursorHelper *mocks.Cursor

	databaseHelper = &mocks.Database{}
	collectionHelper = &mocks.Collection{}
	cursorHelper = &mocks.Cursor{}

	collectionName := domain.CollectionTeam
	t.Run("Success", func(t *testing.T) {
		filter := models.GetMyTeamBody{
			Search: "",
			Size:   10,
			Page:   1,
		}
		filterMongo := bson.M{
			"$and": []bson.M{
				{"name": bson.M{"$regex": filter.Search, "$options": "i"}},
				{"owned_by_id": userId},
			},
		}
		expectedResult := []entities.Team{
			{ID: "1", Name: "Team 1", OwnedByID: userId, SportName: "football"},
			{ID: "2", Name: "Team 2", OwnedByID: userId, SportName: "football"},
		}
		options := options.Find().SetSkip(int64((filter.Page - 1) * filter.Size)).SetLimit(int64(filter.Size))

		collectionHelper.On("Find", mock.Anything, filterMongo, options).Return(cursorHelper, nil).Once()
		collectionHelper.On("CountDocuments", mock.Anything, filterMongo).Return(int64(len(expectedResult)), nil).Once()

		cursorHelper.On("All", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			chatPointer := args.Get(1).(*[]entities.Team)
			*chatPointer = expectedResult
		}).Return(nil).Once()
		databaseHelper.On("Collection", collectionName).Return(collectionHelper).Once()

		tu := repository.NewTeamRepository(databaseHelper, collectionName)
		teams, _, err := tu.GetMyTeam(context.Background(), filter, userId)
		assert.NoError(t, err)
		assert.NotNil(t, teams)
		collectionHelper.AssertExpectations(t)

	})
	t.Run("Error Find", func(t *testing.T) {
		filter := models.GetMyTeamBody{
			Search: "",
			Size:   10,
			Page:   1,
		}
		filterMongo := bson.M{
			"$and": []bson.M{
				{"name": bson.M{"$regex": filter.Search, "$options": "i"}},
				{"owned_by_id": userId},
			},
		}
		options := options.Find().SetSkip(int64((filter.Page - 1) * filter.Size)).SetLimit(int64(filter.Size))

		collectionHelper.On("Find", mock.Anything, filterMongo, options).Return(cursorHelper, errors.New("unexpected")).Once()
		collectionHelper.On("CountDocuments", mock.Anything, filterMongo).Return(int64(0), nil).Once()
		databaseHelper.On("Collection", collectionName).Return(collectionHelper).Once()

		tu := repository.NewTeamRepository(databaseHelper, collectionName)
		teams, _, err := tu.GetMyTeam(context.Background(), filter, userId)
		assert.Error(t, err)
		assert.Nil(t, teams)

	})
	t.Run("Error CountDocuments", func(t *testing.T) {
		var databaseHelper *mocks.Database
		var collectionHelper *mocks.Collection
		var cursorHelper *mocks.Cursor

		databaseHelper = &mocks.Database{}
		collectionHelper = &mocks.Collection{}
		cursorHelper = &mocks.Cursor{}

		collectionName := domain.CollectionTeam

		filter := models.GetMyTeamBody{
			Search: "",
			Size:   10,
			Page:   1,
		}
		filterMongo := bson.M{
			"$and": []bson.M{
				{"name": bson.M{"$regex": filter.Search, "$options": "i"}},
				{"owned_by_id": userId},
			},
		}
		expectedResult := []entities.Team{}
		options := options.Find().SetSkip(int64((filter.Page - 1) * filter.Size)).SetLimit(int64(filter.Size))

		collectionHelper.On("Find", mock.Anything, filterMongo, options).Return(cursorHelper, nil).Once()
		cursorHelper.On("All", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			chatPointer := args.Get(1).(*[]entities.Team)
			*chatPointer = expectedResult
		}).Return(nil).Once()
		collectionHelper.On("CountDocuments", mock.Anything, filterMongo).Return(int64(0), errors.New("unexpected")).Once()
		databaseHelper.On("Collection", collectionName).Return(collectionHelper).Once()

		tu := repository.NewTeamRepository(databaseHelper, collectionName)
		teams, _, err := tu.GetMyTeam(context.Background(), filter, userId)
		assert.Error(t, err)
		assert.Nil(t, teams)
	})
}
