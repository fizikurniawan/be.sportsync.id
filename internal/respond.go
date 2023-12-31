package internal

import (
	"net/http"
	"sportsync/models"

	"github.com/gofiber/fiber/v2"
)

func SendSuccessPaginationRespond(c *fiber.Ctx, code int, data interface{}, page models.Page) error {
	status := http.StatusText(code)
	res := models.SuccessPaginationResponse{
		Status: status,
		Code:   code,
		Data:   data,
		Page:   page,
	}

	return c.Status(code).JSON(res)
}

func SendSuccessRespond(c *fiber.Ctx, code int, data interface{}) error {
	status := http.StatusText(code)
	res := models.SuccessResponse{
		Status: status,
		Code:   code,
		Data:   data,
	}
	return c.Status(code).JSON(res)
}

func SendErrorRespond(c *fiber.Ctx, code int, errors interface{}) error {
	status := http.StatusText(code)
	res := models.ErrorResponse{
		Code:   code,
		Status: status,
		Errors: errors,
	}

	return c.Status(code).JSON(res)
}
