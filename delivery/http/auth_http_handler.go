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
	_ = validate.RegisterValidation("password", internal.PasswordValidation)

	user := new(models.RegisterBody)
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
	if err != nil && !strings.Contains(err.Error(), "no documents in result") {
		if strings.Contains(err.Error(), "email") {
			return internal.SendErrorRespond(c, 400, map[string]string{"email": err.Error()})
		}
		return internal.SendErrorRespond(c, 400, map[string]string{"error": err.Error()})
	}

	return internal.SendSuccessRespond(c, 201, user)
}

func (uh *AuthHandler) Login(c *fiber.Ctx) (err error) {
	var validate = validator.New()
	_ = validate.RegisterValidation("password", internal.PasswordValidation)

	body := new(models.LoginBody)
	if err := c.BodyParser(body); err != nil {
		return err
	}

	errs := validate.Struct(body)
	if errs != nil {
		errorMap := make(map[string][]string)
		for _, err := range errs.(validator.ValidationErrors) {
			fieldName := strings.ToLower(err.Field())
			errorMap[fieldName] = append(errorMap[fieldName], internal.MsgForTag(err))
		}

		return internal.SendErrorRespond(c, 400, errorMap)

	}

	var user models.LoginResponse
	user, err = uh.AuthUsecase.Login(c.UserContext(), *body)
	if err != nil {
		return internal.SendErrorRespond(c, 401, map[string][]string{"password": {"invalid credentials"}})
	}

	return internal.SendSuccessRespond(c, 200, user)
}
