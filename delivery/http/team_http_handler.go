package http

import (
	"sportsync/bootstrap"
	"sportsync/domain"
	"sportsync/entities"
	"sportsync/internal"
	"sportsync/models"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type TeamHandler struct {
	TeamUsecase domain.TeamUsecase
	Env         *bootstrap.Env
}

func (th *TeamHandler) Create(c *fiber.Ctx) (err error) {
	var validate = validator.New()
	_ = validate.RegisterValidation("password", internal.PasswordValidation)

	team := new(models.TeamBody)
	if err := c.BodyParser(team); err != nil {
		return err
	}

	errs := validate.Struct(team)
	if errs != nil {
		errorMap := make(map[string][]string)
		for _, err := range errs.(validator.ValidationErrors) {
			fieldName := strings.ToLower(err.Field())
			switch fieldName {
			case "sportname":
				fieldName = "sport_name"
			}
			errorMap[fieldName] = append(errorMap[fieldName], internal.MsgForTag(err))
		}

		return internal.SendErrorRespond(c, 400, errorMap)
	}

	uid := c.Locals("userId")
	err = th.TeamUsecase.Create(c.UserContext(), *team, uid.(string))
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return internal.SendErrorRespond(c, 400, map[string][]string{"name": {err.Error()}})
		}
		return internal.SendErrorRespond(c, 400, map[string][]string{"error": {err.Error()}})
	}

	return internal.SendSuccessRespond(c, 201, team)
}

func (th *TeamHandler) MyTeam(c *fiber.Ctx) (err error) {
	var teams []entities.Team
	var page models.Page

	filter := new(models.GetMyTeamBody)
	if err := c.BodyParser(filter); err != nil {
		return err
	}

	teams, page, err = th.TeamUsecase.GetMyTeam(c.UserContext(), *filter, c.Locals("userId").(string))

	if err != nil {
		return
	}
	return internal.SendSuccessPaginationRespond(c, 200, teams, page)
}
