package entities

import "time"

type Team struct {
	ID        string     `bson:"_id,omitempty" json:"id"`
	Name      string     `bson:"name" json:"name"`
	SportName string     `bson:"sport_name" json:"sport_name"`
	OwnedByID string     `bson:"owned_by_id" json:"owned_by_id"`
	Logo      string     `bson:"logo" json:"logo"`
	CreatedAt *time.Time `bson:"created_at" json:"created_at"`
}
