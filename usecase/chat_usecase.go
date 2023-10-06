package usecase

import (
	"context"
	"sportsync/bootstrap"
	"sportsync/domain"
	"sportsync/entities"
	"time"
)

type chatUsecase struct {
	chatRepository domain.ChatRepository
	teamRepository domain.TeamRepository
	contextTimeout time.Duration
	env            *bootstrap.Env
}

func NewChatUsecase(chatRepository domain.ChatRepository, teamRepository domain.TeamRepository, timeout time.Duration, env *bootstrap.Env) domain.ChatUsecase {
	return &chatUsecase{
		chatRepository: chatRepository,
		teamRepository: teamRepository,
		contextTimeout: timeout,
		env:            env,
	}
}

func (cu *chatUsecase) ChatOnTeam(teamId string) error {
	return nil
}

func (cu *chatUsecase) GetRecentMessages(c context.Context, teamId string, userId string, senderId string) (chats []entities.Chat, err error) {
	chats, err = cu.chatRepository.GetRecentMessages(c, teamId, userId, senderId)

	return
}

func (cu *chatUsecase) Insert(c context.Context, chat entities.Chat) error {
	if err := cu.chatRepository.Insert(c, chat); err != nil {
		return err
	}
	return nil
}

func (cu *chatUsecase) IsTeamExists(c context.Context, teamId string) (bool, error) {
	team, err := cu.teamRepository.GetByID(c, teamId)
	if err != nil {
		return false, err
	}

	if team.Name != "" {
		return true, nil
	}
	return false, nil
}
