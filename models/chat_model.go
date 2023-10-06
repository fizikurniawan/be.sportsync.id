package models

type ChatResponse struct {
	Message   string `json:"message"`
	SenderId  string `json:"sender_id"`
	MessageId string `json:"message_id"`
}
