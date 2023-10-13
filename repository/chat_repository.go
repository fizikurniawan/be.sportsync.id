package repository

import (
	"context"
	"sportsync/domain"
	"sportsync/entities"
	"sportsync/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

type chatRepository struct {
	database   mongo.Database
	collection string
}

func NewChatRepository(db mongo.Database, collection string) domain.ChatRepository {
	return &chatRepository{
		database:   db,
		collection: collection,
	}
}

func (cr *chatRepository) GetRecentMessages(c context.Context, teamId string, userId string, senderId string) (chats []entities.Chat, err error) {
	collection := cr.database.Collection(cr.collection)

	filterMongo := bson.M{
		"$and": []bson.M{
			{"receiver_id": userId},
			{"sender_id": senderId},
		},
	}

	if teamId != "" {
		filterMongo = bson.M{
			"$and": []bson.M{
				{"team_id": teamId},
			},
		}
	}

	cursor, err := collection.Find(c, filterMongo)
	if err != nil {
		return
	}

	err = cursor.All(c, &chats)
	return

}

func (cr *chatRepository) Insert(c context.Context, chat entities.Chat) error {
	collection := cr.database.Collection(cr.collection)

	_, err := collection.InsertOne(c, chat)

	return err
}
