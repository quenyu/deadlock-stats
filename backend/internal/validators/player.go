package validators

import (
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	cErrors "github.com/quenyu/deadlock-stats/internal/errors"
)

func ValidateSteamID(steamID string) error {
	steamID = strings.TrimSpace(steamID)

	if steamID == "" {
		return cErrors.ErrInvalidSteamID
	}

	id, err := strconv.ParseInt(steamID, 10, 64)
	if err != nil {
		return cErrors.ErrInvalidSteamID
	}

	if id <= 0 {
		return cErrors.ErrInvalidSteamID
	}

	if id < 76561197960265728 {
		return cErrors.ErrInvalidSteamID
	}

	return nil
}

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
