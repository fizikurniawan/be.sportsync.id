package domain

import (
	"context"
	"sportsync/entities"
)

const CollectionChat = "chats"

type ChatRepository interface {
	GetRecentMessages(c context.Context, teamId string, userId string, senderId string) (chats []entities.Chat, err error)
	Insert(c context.Context, chat entities.Chat) error
}

type ChatUsecase interface {
	ChatOnTeam(teamId string) error
	GetRecentMessages(c context.Context, teamId string, userId string, senderId string) (chats []entities.Chat, err error)
	Insert(c context.Context, chat entities.Chat) error
	IsTeamExists(c context.Context, teamId string) (bool, error)
}
