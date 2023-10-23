package models

type CupBody struct {
	Name           string `json:"name"`
	MaxParticipant int    `json:"max_participant"`
	Format         string `form:"format"`
	StartDate      string `json:"start_date"`
	Season         string `json:"season"`
}
