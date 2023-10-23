package http

import (
	"sportsync/bootstrap"
	"sportsync/domain"
	"sportsync/internal"
	"sportsync/models"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CompetitionHandler struct {
	CompetitionUsecase domain.CompetitionUsecase
	Env                *bootstrap.Env
}

func (ch *CompetitionHandler) CreateCup(c *fiber.Ctx) error {
	var validate = validator.New()

	cup := new(models.CupBody)
	if err := c.BodyParser(cup); err != nil {
		return err
	}

	errs := validate.Struct(cup)
	if errs != nil {
		errorMap := make(map[string][]string)
		for _, err := range errs.(validator.ValidationErrors) {
			fieldName := strings.ToLower(err.Field())
			errorMap[fieldName] = append(errorMap[fieldName], internal.MsgForTag(err))
		}

		return internal.SendErrorRespond(c, 400, errorMap)
	}

	// insert to database
	if err := ch.CompetitionUsecase.CreateCup(c.UserContext(), cup); err != nil {
		return internal.SendErrorRespond(c, 400, map[string]string{"error": err.Error()})

	}

	return c.JSON(fiber.Map{"OK": 124})
}
