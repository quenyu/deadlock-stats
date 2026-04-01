package validators

import (
	"strconv"
	"strings"

	cErrors "github.com/quenyu/deadlock-stats/internal/errors"
)

// ValidateSteamID checks the validity of a Steam ID
func ValidateSteamID(steamID string) error {
	if steamID == "" {
		return cErrors.ErrInvalidSteamID
	}

	steamID = strings.TrimSpace(steamID)

	id, err := strconv.ParseInt(steamID, 10, 64)
	if err != nil {
		return cErrors.ErrInvalidSteamID
	}

	if id <= 0 {
		return cErrors.ErrInvalidSteamID
	}

	if id > 999999999999999999 {
		return cErrors.ErrInvalidSteamID
	}

	return nil
}
