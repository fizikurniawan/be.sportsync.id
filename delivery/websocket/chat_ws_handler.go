package websocket

import (
	"context"
	"sportsync/bootstrap"
	"sportsync/domain"
	"sportsync/entities"
	"sportsync/internal"
	"sportsync/models"
	"sportsync/ws"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type ChatWSHandler struct {
	ChatUsecase domain.ChatUsecase
	Env         *bootstrap.Env
}

func NewChatWSHandler(chatUsecase domain.ChatUsecase, env *bootstrap.Env) *ChatWSHandler {
	return &ChatWSHandler{
		ChatUsecase: chatUsecase,
		Env:         env,
	}
}

func (ch *ChatWSHandler) ChatOnTeam(c *websocket.Conn) {
	userData := c.Locals("user").(models.UserClaims)
	teamId := c.Params("teamId")
	userId := userData.ID

	// validate teamId is exists
	isTeamExists, _ := ch.ChatUsecase.IsTeamExists(context.TODO(), teamId)
	if !isTeamExists {
		c.WriteJSON(fiber.Map{"team_id": "team not found"})
		c.CloseHandler()(1000, "")
	}

	// lock clients
	ws.ClientsMu.Lock()
	ws.Clients[c] = userId
	ws.ClientsMu.Unlock()

	// retrieve saved message
	chats, errRetrieve := ch.ChatUsecase.GetRecentMessages(context.TODO(), teamId, userId, "")
	if errRetrieve == nil {
		c.WriteJSON(fiber.Map{"chats": chats})
	}

	for {
		var chat entities.Chat
		if err := c.ReadJSON(&chat); err != nil {
			break
		}

		msgId := internal.GenerateMessageID()
		chat.ID = msgId
		chat.SenderID = userId
		chat.TeamID = teamId

		if err := ch.ChatUsecase.Insert(context.TODO(), chat); err != nil {
			break
		}

		ws.ClientsMu.Lock()
		for clientConn, clientID := range ws.Clients {
			if clientID != userId {
				if err := clientConn.WriteJSON(chat); err != nil {
					break
				}
			}
		}
		ws.ClientsMu.Unlock()
	}
}
