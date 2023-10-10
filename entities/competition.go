package entities

import "time"

type Cup struct {
	ID   string `bson:"_id,omitempty" json:"id"`
	Name string `json:"name" bson:"name"`
}

type League struct {
	ID   string `bson:"_id,omitempty" json:"id"`
	Name string `json:"name" bson:"name"`
}

type Fixtures struct {
	CupID    string    `bson:"cup_id" json:"-"`
	LeagueID string    `bson:"league_id" json:"-"`
	Date     time.Time `bson:"date" json:"date"`
}
