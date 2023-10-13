package repository_test

import (
	"context"
	"errors"
	"sportsync/domain"
	"sportsync/entities"
	"sportsync/mongo/mocks"
	"sportsync/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetRecentMessages(t *testing.T) {
	var databaseHelper *mocks.Database
	var collectionHelper *mocks.Collection
	var cursorHelper *mocks.Cursor

	databaseHelper = &mocks.Database{}
	collectionHelper = &mocks.Collection{}
	cursorHelper = &mocks.Cursor{}

	userId := string(primitive.NewObjectID().Hex())
	teamId := string(primitive.NewObjectID().Hex())
	senderId := string(primitive.NewObjectID().Hex())
	collectionName := domain.CollectionChat

	t.Run("Success Get by UserID", func(t *testing.T) {

		// test room isn't team
		filter := bson.M{
			"$and": []bson.M{
				{"receiver_id": userId},
				{"sender_id": senderId},
			},
		}

		expectedResult := []entities.Chat{
			{ID: "1", Message: "Message 1", SenderID: senderId, ReceiverID: userId},
			{ID: "2", Message: "Message 2", SenderID: senderId, ReceiverID: userId},
		}

		collectionHelper.On("Find", mock.Anything, filter).Return(cursorHelper, nil).Once()
		databaseHelper.On("Collection", collectionName).Return(collectionHelper)
		cursorHelper.On("All", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			chatPointer := args.Get(1).(*[]entities.Chat)
			*chatPointer = expectedResult
		}).Return(nil).Once()

		cr := repository.NewChatRepository(databaseHelper, collectionName)

		chats, err := cr.GetRecentMessages(context.Background(), "", userId, senderId)
		assert.NoError(t, err)
		assert.NotNil(t, chats)
		for _, v := range chats {
			assert.Equal(t, v.ReceiverID, userId)
			assert.Equal(t, v.SenderID, senderId)
		}
		collectionHelper.AssertExpectations(t)
		cursorHelper.AssertExpectations(t)
		databaseHelper.AssertExpectations(t)

	})

	t.Run("Success Get by TeamId", func(t *testing.T) {
		filter := bson.M{
			"$and": []bson.M{
				{"team_id": teamId},
			},
		}
		expectedResult := []entities.Chat{
			{ID: "1", Message: "Message 1", SenderID: "", ReceiverID: "", TeamID: teamId},
			{ID: "2", Message: "Message 2", SenderID: "", ReceiverID: "", TeamID: teamId},
		}

		collectionHelper.On("Find", mock.Anything, filter).Return(cursorHelper, nil).Once()
		databaseHelper.On("Collection", collectionName).Return(collectionHelper)
		cursorHelper.On("All", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			chatPointer := args.Get(1).(*[]entities.Chat)
			*chatPointer = expectedResult
		}).Return(nil).Once()

		cr := repository.NewChatRepository(databaseHelper, collectionName)

		chats, err := cr.GetRecentMessages(context.Background(), teamId, "", "")
		assert.NoError(t, err)
		assert.NotNil(t, chats)
		for _, v := range chats {
			assert.Equal(t, v.ReceiverID, "")
			assert.Equal(t, v.SenderID, "")
			assert.Equal(t, v.TeamID, teamId)
		}
		collectionHelper.AssertExpectations(t)
		cursorHelper.AssertExpectations(t)
		databaseHelper.AssertExpectations(t)
	})

	t.Run("Error Get by UserID", func(t *testing.T) {
		// test room isn't team
		filter := bson.M{
			"$and": []bson.M{
				{"receiver_id": userId},
				{"sender_id": senderId},
			},
		}

		collectionHelper.On("Find", mock.Anything, filter).Return(cursorHelper, errors.New("unexpected")).Once()
		databaseHelper.On("Collection", collectionName).Return(collectionHelper)

		cr := repository.NewChatRepository(databaseHelper, collectionName)

		chats, err := cr.GetRecentMessages(context.Background(), "", userId, senderId)
		assert.Error(t, err)
		assert.Nil(t, chats)
		collectionHelper.AssertExpectations(t)
		cursorHelper.AssertExpectations(t)
		databaseHelper.AssertExpectations(t)
	})

	t.Run("Error Get by TeamID", func(t *testing.T) {
		filter := bson.M{
			"$and": []bson.M{
				{"team_id": teamId},
			},
		}

		collectionHelper.On("Find", mock.Anything, filter).Return(cursorHelper, errors.New("unexpected")).Once()
		databaseHelper.On("Collection", collectionName).Return(collectionHelper)

		cr := repository.NewChatRepository(databaseHelper, collectionName)

		chats, err := cr.GetRecentMessages(context.Background(), teamId, "", "")
		assert.Error(t, err)
		assert.Nil(t, chats)
		collectionHelper.AssertExpectations(t)
		cursorHelper.AssertExpectations(t)
		databaseHelper.AssertExpectations(t)
	})
}

func TestInsert(t *testing.T) {
	var databaseHelper *mocks.Database
	var collectionHelper *mocks.Collection

	databaseHelper = &mocks.Database{}
	collectionHelper = &mocks.Collection{}

	userId := string(primitive.NewObjectID().Hex())
	// teamId := string(primitive.NewObjectID().Hex())
	senderId := string(primitive.NewObjectID().Hex())
	collectionName := domain.CollectionChat

	mockChat := &entities.Chat{
		ID:         primitive.NewObjectID().Hex(),
		Message:    "Message",
		SenderID:   senderId,
		ReceiverID: userId,
	}

	mockEmptyChat := &entities.Chat{}
	mockChatID := primitive.NewObjectID()

	t.Run("success", func(t *testing.T) {

		collectionHelper.On("InsertOne", mock.Anything, mock.AnythingOfType("entities.Chat")).Return(mockChatID, nil).Once()

		databaseHelper.On("Collection", collectionName).Return(collectionHelper)

		cr := repository.NewChatRepository(databaseHelper, collectionName)

		err := cr.Insert(context.Background(), *mockChat)

		assert.NoError(t, err)

		collectionHelper.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		collectionHelper.On("InsertOne", mock.Anything, mock.AnythingOfType("entities.Chat")).Return(mockEmptyChat, errors.New("Unexpected")).Once()

		databaseHelper.On("Collection", collectionName).Return(collectionHelper)

		ur := repository.NewChatRepository(databaseHelper, collectionName)

		err := ur.Insert(context.Background(), *mockEmptyChat)

		assert.Error(t, err)

		collectionHelper.AssertExpectations(t)
	})
}
