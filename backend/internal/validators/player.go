package validators

import (
	"strings"

	"github.com/go-playground/validator/v10"
	cErrors "github.com/quenyu/deadlock-stats/internal/errors"
)

func validateSteamIDTag(fl validator.FieldLevel) bool {
	steamID := fl.Field().String()
	return ValidateSteamID(steamID) == nil
}

func ValidatePlayerSearchQuery(query string) error {
	query = strings.TrimSpace(query)

	if query == "" {
		return cErrors.ErrInvalidQuery
	}

	if len(query) < 2 {
		return cErrors.ErrFieldTooShort
	}

	if len(query) > 100 {
		return cErrors.ErrFieldTooLong
	}

	return nil
}
