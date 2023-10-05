package models

type Page struct {
	Total      int `json:"total"`
	Size       int `json:"size"`
	Current    int `json:"current"`
	TotalPages int `json:"total_pages"`
}

type SuccessResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Data   any    `json:"data"`
}

type SuccessPaginationResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Data   any    `json:"data"`
	Page   Page   `json:"page"`
}

type DetailErr struct{}

type ErrorResponse struct {
	Status string      `json:"status"`
	Code   int         `json:"code"`
	Errors interface{} `json:"errors"`
}
