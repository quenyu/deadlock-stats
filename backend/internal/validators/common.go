package validators

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	cErrors "github.com/quenyu/deadlock-stats/internal/errors"
)

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)

	if email == "" {
		return cErrors.ErrInvalidEmail
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)

	if !emailRegex.MatchString(email) {
		return cErrors.ErrInvalidEmail
	}

	return nil
}

func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)

	if username == "" {
		return cErrors.ErrInvalidUsername
	}

	if len(username) < 3 {
		return cErrors.ErrFieldTooShort
	}

	if len(username) > 50 {
		return cErrors.ErrFieldTooLong
	}

	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !usernameRegex.MatchString(username) {
		return cErrors.ErrInvalidUsername
	}

	return nil
}

func validateUsernameTag(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	return ValidateUsername(username) == nil
}

func ValidateDescription(description string, maxLen int) error {
	description = strings.TrimSpace(description)

	if len(description) > maxLen {
		return cErrors.ErrFieldTooLong
	}

	return nil
}

func ValidateRating(rating int) error {
	if rating < 1 || rating > 5 {
		return cErrors.ErrInvalidRating
	}
	return nil
}
