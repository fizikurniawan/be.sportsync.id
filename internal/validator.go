package internal

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

func MsgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "password":
		return "Invalid password format"
	}

	return fe.Error() // default error
}

func PasswordValidation(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Minimal panjang password adalah 8 karakter
	if len(password) < 8 {
		return false
	}

	var (
		hasUpperCase bool
		hasLowerCase bool
		hasDigit     bool
		hasSpecial   bool
	)

	// Memeriksa setiap karakter dalam password
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpperCase = true
		case unicode.IsLower(char):
			hasLowerCase = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// Memastikan password memenuhi semua kriteria
	return hasUpperCase && hasLowerCase && hasDigit && hasSpecial
}
