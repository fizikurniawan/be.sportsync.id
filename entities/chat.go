package entities

import "time"

type Chat struct {
	ID         string    `bson:"_id,omitempty" json:"id"`
	Message    string    `bson:"message" json:"message"`
	Photo      string    `bson:"photo" json:"photo"`
	SendAt     time.Time `bson:"send_at" json:"send_at"`
	SenderID   string    `bson:"sender_id" json:"sender_id"`
	ReceiverID string    `bson:"receiver_id" json:"receiver_id"`
	TeamID     string    `bson:"team_id" json:"team_id"`
}
