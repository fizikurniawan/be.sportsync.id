package repository_test

import (
	"context"
	"errors"
	"testing"

	"sportsync/domain"
	"sportsync/entities"
	"sportsync/mongo/mocks"
	"sportsync/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCreate(t *testing.T) {

	var databaseHelper *mocks.Database
	var collectionHelper *mocks.Collection

	databaseHelper = &mocks.Database{}
	collectionHelper = &mocks.Collection{}

	collectionName := domain.CollectionUser

	mockUser := &entities.User{
		ID:       primitive.NewObjectID().Hex(),
		Name:     "Test",
		Email:    "test@gmail.com",
		Password: "password",
	}

	mockEmptyUser := &entities.User{}
	mockUserID := primitive.NewObjectID()

	t.Run("success", func(t *testing.T) {

		collectionHelper.On("InsertOne", mock.Anything, mock.AnythingOfType("*entities.User")).Return(mockUserID, nil).Once()

		databaseHelper.On("Collection", collectionName).Return(collectionHelper)

		ur := repository.NewUserRepository(databaseHelper, collectionName)

		err := ur.Create(context.Background(), mockUser)

		assert.NoError(t, err)

		collectionHelper.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		collectionHelper.On("InsertOne", mock.Anything, mock.AnythingOfType("*entities.User")).Return(mockEmptyUser, errors.New("Unexpected")).Once()

		databaseHelper.On("Collection", collectionName).Return(collectionHelper)

		ur := repository.NewUserRepository(databaseHelper, collectionName)

		err := ur.Create(context.Background(), mockEmptyUser)

		assert.Error(t, err)

		collectionHelper.AssertExpectations(t)
	})

}

func TestFetch(t *testing.T) {
	var databaseHelper *mocks.Database
	var collectionHelper *mocks.Collection
	var cursorHelper *mocks.Cursor

	databaseHelper = &mocks.Database{}
	collectionHelper = &mocks.Collection{}
	cursorHelper = &mocks.Cursor{}

	collectionName := domain.CollectionUser

	t.Run("Success", func(t *testing.T) {

		filter := bson.D{}
		projection := bson.D{{Key: "password", Value: 0}}
		expectedUsers := []entities.User{
			{ID: "1", Name: "User 1", Email: "user1@example.com", Password: "pass1"},
			{ID: "2", Name: "User 2", Email: "user2@example.com", Password: "pass2"},
		}

		collectionHelper.On("Find", mock.Anything, filter, options.Find().SetProjection(projection)).Return(cursorHelper, nil).Once()
		databaseHelper.On("Collection", collectionName).Return(collectionHelper)
		cursorHelper.On("All", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			usersPointer := args.Get(1).(*[]entities.User)
			*usersPointer = expectedUsers
		}).Return(nil).Once()

		ur := repository.NewUserRepository(databaseHelper, collectionName)

		users, err := ur.Fetch(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.Equal(t, users[0].Name, expectedUsers[0].Name)

		collectionHelper.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {

		filter := bson.D{}
		projection := bson.D{{Key: "password", Value: 0}}

		collectionHelper.On("Find", mock.Anything, filter, options.Find().SetProjection(projection)).Return(cursorHelper, errors.New("Unexpected")).Once()
		databaseHelper.On("Collection", collectionName).Return(collectionHelper)

		ur := repository.NewUserRepository(databaseHelper, collectionName)

		users, err := ur.Fetch(context.Background())

		assert.Error(t, err)
		assert.Nil(t, users)

		collectionHelper.AssertExpectations(t)
	})
	t.Run("Success Fetech Empty data", func(t *testing.T) {

		filter := bson.D{}
		projection := bson.D{{Key: "password", Value: 0}}

		collectionHelper.On("Find", mock.Anything, filter, options.Find().SetProjection(projection)).Return(cursorHelper, nil).Once()
		databaseHelper.On("Collection", collectionName).Return(collectionHelper)
		cursorHelper.On("All", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			args.Get(1)
		}).Return(nil).Once()

		ur := repository.NewUserRepository(databaseHelper, collectionName)

		users, err := ur.Fetch(context.Background())

		assert.NoError(t, err)
		assert.Nil(t, users)

		collectionHelper.AssertExpectations(t)
	})
}

func TestGetByEmail(t *testing.T) {
	var databaseHelper *mocks.Database
	var collectionHelper *mocks.Collection
	var singleResultHelper *mocks.SingleResult

	databaseHelper = &mocks.Database{}
	collectionHelper = &mocks.Collection{}
	singleResultHelper = &mocks.SingleResult{}

	collectionName := domain.CollectionUser
	t.Run("Success", func(t *testing.T) {
		expectedUser := &entities.User{Name: "User1", Email: "user1@gmail.com"}

		collectionHelper.On("FindOne", mock.Anything, mock.Anything).Return(singleResultHelper).Once()
		singleResultHelper.On("Decode", &entities.User{}).Run(func(args mock.Arguments) {
			usersPointer := args.Get(0).(*entities.User)
			*usersPointer = *expectedUser
		}).Return(nil).Once()

		databaseHelper.On("Collection", collectionName).Return(collectionHelper).Once()
		ur := repository.NewUserRepository(databaseHelper, collectionName)

		user, err := ur.GetByEmail(context.Background(), expectedUser.Email)

		collectionHelper.AssertExpectations(t)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.Email, user.Email, "User's Email should match expected value.")

	})
}
func TestGetByID(t *testing.T) {
	var databaseHelper *mocks.Database
	var collectionHelper *mocks.Collection
	var singleResultHelper *mocks.SingleResult

	databaseHelper = &mocks.Database{}
	collectionHelper = &mocks.Collection{}
	singleResultHelper = &mocks.SingleResult{}

	collectionName := domain.CollectionUser
	t.Run("Success", func(t *testing.T) {
		userId := string(primitive.NewObjectID().Hex())
		expectedUser := &entities.User{Name: "User1", Email: "user1@gmail.com", ID: userId}

		collectionHelper.On("FindOne", mock.Anything, mock.Anything).Return(singleResultHelper).Once()
		singleResultHelper.On("Decode", &entities.User{}).Run(func(args mock.Arguments) {
			usersPointer := args.Get(0).(*entities.User)
			*usersPointer = *expectedUser
		}).Return(nil).Once()

		databaseHelper.On("Collection", collectionName).Return(collectionHelper).Once()
		ur := repository.NewUserRepository(databaseHelper, collectionName)

		user, err := ur.GetByID(context.Background(), expectedUser.ID)

		collectionHelper.AssertExpectations(t)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.Email, user.Email, "User's Email should match expected value.")
	})
	t.Run("Error Hex ID", func(t *testing.T) {
		expectedUser := &entities.User{Name: "User1", Email: "user1@gmail.com", ID: "INVALID-ID"}

		databaseHelper.On("Collection", collectionName).Return(collectionHelper).Once()
		ur := repository.NewUserRepository(databaseHelper, collectionName)

		_, err := ur.GetByID(context.Background(), expectedUser.ID)
		assert.Error(t, err)
	})
}
