package models

type TeamBody struct {
	Name      string `json:"name" validate:"required"`
	SportName string `json:"sport_name" validate:"required"`
	Logo      string `json:"logo"`
}

type GetMyTeamBody struct {
	Search string `json:"search"`
	Size   int    `json:"size"`
	Page   int    `json:"page"`
}
