package entities

import "time"

type Cup struct {
	Name           string    `json:"name"`
	MaxParticipant int       `json:"max_participant"`
	Format         string    `form:"format"`
	StartDate      time.Time `json:"start_date"`
	Season         string    `json:"season"`
}

type League struct {
	Name           string    `json:"name"`
	MaxParticipant int       `json:"max_participant"`
	Format         string    `form:"format"`
	StartDate      time.Time `json:"start_date"`
	Season         string    `json:"season"`
}

type Fixtures struct {
	CupID    string    `bson:"cup_id" json:"-"`
	LeagueID string    `bson:"league_id" json:"-"`
	Date     time.Time `bson:"date" json:"date"`
}
