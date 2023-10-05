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

type AuthHandler struct {
	AuthUsecase domain.AuthUsecase
	Env         *bootstrap.Env
}

func (uh *AuthHandler) Register(c *fiber.Ctx) (err error) {
	var validate = validator.New()
	user := new(models.RegisterRequest)
	if err := c.BodyParser(user); err != nil {
		return err
	}

	errs := validate.Struct(user)
	if errs != nil {
		errorMap := make(map[string][]string)
		for _, err := range errs.(validator.ValidationErrors) {
			fieldName := strings.ToLower(err.Field())
			errorMap[fieldName] = append(errorMap[fieldName], internal.MsgForTag(err))
		}

		return internal.SendErrorRespond(c, 400, errorMap)

	}

	err = uh.AuthUsecase.Register(c.UserContext(), *user)
	if err != nil {
		return internal.SendErrorRespond(c, 400, map[string]string{"a": "baba"})
	}

	return internal.SendSuccessRespond(c, 201, user)
}
